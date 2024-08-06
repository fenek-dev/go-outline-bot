package payment_service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/fenek-dev/go-outline-bot/pkg/payment_service/mocks"
	"github.com/stretchr/testify/mock"
)

func generateRequest(hooks ...func(transaction *CreateTransactionRequest)) CreateTransactionRequest {
	transaction := CreateTransactionRequest{
		TxUUID:       "195ede64-d103-4b0b-8fc6-5a1aa1b568b2",
		Amount:       100,
		UserID:       "1",
		MethodID:     1,
		MethodType:   "payment",
		CurrencyCode: "RUB",
		Description:  "fuck",
		Service:      "test",
		PostbackURL:  "http://test.com/postback",
		SuccessURL:   "http://test.com/success",
		FailURL:      "http://test.com/fail",
		Customer: TransactionCustomer{
			ID:    "1",
			Email: "test@gmail.com",
			Phone: "+79779486965",
		},
		Items: []TransactionItem{
			{
				ID:       "1",
				Name:     "test",
				Type:     "test",
				Price:    100,
				Quantity: 1,
				SKU:      nil,
			},
		},
	}
	for _, hook := range hooks {
		hook(&transaction)
	}

	return transaction
}

func Test_CreateTransaction(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		request  CreateTransactionRequest
		response *http.Response
		check    func(*testing.T, *CreateTransactionResponse, error)
		error    error
	}{
		{
			name:    "should_create_transaction",
			request: generateRequest(),
			response: &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"id": 104,
						"uuid": "195ede64-d103-4b0b-8fc6-5a1aa1b568b4",
						"amount": 100,
						"method_type": null,
						"status": "pending",
						"result": {
							"success": true,
							"error": null,
							"data": {
								"type": "redirect",
								"message": "Please go to the payment form.",
								"redirectUrl": "https://securepayments.tinkoff.ru/dTw1GuWr"
							},
							"externalId": "4804909179",
							"postbacks": null,
							"convertedAmount": null,
							"tryRepeatLater": true
						},
						"error": null,
						"success": true,
						"action": {
							"type": "redirect",
							"message": "Please go to the payment form.",
							"redirectUrl": "https://securepayments.tinkoff.ru/dTw1GuWr"
						},
						"created_at": "2024-08-02T19:57:59.000000Z"
					}
				`))),
			},
			check: func(t *testing.T, response *CreateTransactionResponse, err error) {
				assert.Equal(t, err, nil)
				assert.Equal(t, 104, response.ID)
			},
			error: nil,
		},
		{
			name: "validation_error",
			request: generateRequest(func(transaction *CreateTransactionRequest) {
				transaction.Items = []TransactionItem{}
			}),
			response: &http.Response{
				StatusCode: 422,
				Body: io.NopCloser(bytes.NewReader([]byte(`
					{
						"message": "The items field is required.",
						"errors": {
							"items": [
								"The items field is required."
							]
						}
					}
				`))),
			},
			check: func(t *testing.T, response *CreateTransactionResponse, err error) {
				assert.Equal(t, err.Error(), "create transaction: The items field is required.")
			},
			error: &ErrorResponse{
				Message: "The items field is required.",
				Errors:  []string{"The items field is required."},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			httpClient := mocks.NewHTTPClient(t)
			httpClient.On("Do", mock.Anything).Return(tc.response, tc.error)

			client := NewClient("", WithHTTPClient(httpClient))

			ctx := context.Background()

			// Act
			response, err := client.CreateTransaction(ctx, 1, tc.request)

			// Assert
			tc.check(t, response, err)
		})
	}
}
