package payment_service

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type TransactionStatus string

const (
	TransactionStatusAuthorized TransactionStatus = "authorized"
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusDeclined   TransactionStatus = "declined"
	TransactionStatusCancelled  TransactionStatus = "cancelled"
)

type CreateTransactionResponse struct {
	ID         int                   `json:"id"`
	UUID       string                `json:"uuid"`
	Amount     int                   `json:"amount"`
	MethodType *string               `json:"method_type"`
	Status     TransactionStatus     `json:"status"`
	Error      *string               `json:"error"`
	Success    bool                  `json:"success"`
	Action     TransactionResultData `json:"action"`
	CreatedAt  time.Time             `json:"created_at"`

	Result TransactionResult `json:"result"`
}

//type Cre struct {
//	Id         int         `json:"id"`
//	Uuid       string      `json:"uuid"`
//	Amount     int         `json:"amount"`
//	MethodType interface{} `json:"method_type"`
//	Status     string      `json:"status"`
//	Result     struct {
//		Success bool        `json:"success"`
//		Error   interface{} `json:"error"`
//		Data    struct {
//			Type        string `json:"type"`
//			Message     string `json:"message"`
//			RedirectUrl string `json:"redirectUrl"`
//		} `json:"data"`
//		ExternalId      string      `json:"externalId"`
//		Postbacks       interface{} `json:"postbacks"`
//		ConvertedAmount interface{} `json:"convertedAmount"`
//		TryRepeatLater  bool        `json:"tryRepeatLater"`
//	} `json:"result"`
//	Error   interface{} `json:"error"`
//	Success bool        `json:"success"`
//	Action  struct {
//		Type        string `json:"type"`
//		Message     string `json:"message"`
//		RedirectUrl string `json:"redirectUrl"`
//	} `json:"action"`
//	CreatedAt time.Time `json:"created_at"`
//}

type TransactionCustomer struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type TransactionItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Price    int     `json:"price"`
	Quantity int     `json:"quantity"`
	SKU      *string `json:"sku"`
}

type CreateTransactionRequest struct {
	TxUUID       string `json:"tx_uuid"`
	Amount       int    `json:"amount"`
	UserID       string `json:"user_id"`
	MethodID     int    `json:"method_id"`
	MethodType   string `json:"method_type"`
	CurrencyCode string `json:"currency_code"`
	Description  string `json:"description"` // Описание транзакции
	Service      string `json:"service"`     // Название сервис

	PostbackURL string `json:"postback_url"`
	SuccessURL  string `json:"success_url"`
	FailURL     string `json:"fail_url"`

	Customer TransactionCustomer `json:"customer"`

	Items []TransactionItem `json:"items"`
}

func (c *Client) CreateTransaction(ctx context.Context, userId int64, payload CreateTransactionRequest) (*CreateTransactionResponse, error) {
	response := &CreateTransactionResponse{}

	req, err := c.NewRequest(ctx, http.MethodPost, "v1/transactions", payload)
	if err != nil {
		return response, err
	}

	if err = c.Send(req, response); err != nil {
		if e, ok := err.(*ErrorResponse); ok {
			return response, fmt.Errorf("create transaction: %s", e)
		}

		return response, err
	}

	return response, nil
}
