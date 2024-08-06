package payment_service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/fenek-dev/go-outline-bot/pkg/payment_service/mocks"
	"github.com/stretchr/testify/mock"
)

func NewTestLogger(t *testing.T) *slog.Logger {
	t.Helper()
	return slog.New(
		slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{}),
	)
}

type MockBody struct {
	Param string `json:"value"`
}

func Test_NewRequest(t *testing.T) {
	t.Run("should create new request", func(t *testing.T) {
		// Arrange
		client := NewClient("", WithHTTPClient(mocks.NewHTTPClient(t)))
		ctx := context.Background()

		// Act
		request, err := client.NewRequest(ctx, "GET", "test", nil)

		// Assert
		assert.Equal(t, err, nil)
		assert.Equal(t, request.URL.String(), "/test")
		assert.Equal(t, request.Method, "GET")
	})
}

func Test_Send(t *testing.T) {
	t.Run("should send request", func(t *testing.T) {
		// Arrange
		mockBody, err := json.Marshal(&MockBody{"test"})
		if err != nil {
			t.Fatalf("failed to marshal test body: %s", err)
		}

		httpClient := mocks.NewHTTPClient(t)
		httpClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(mockBody)),
		}, nil)

		client := NewClient("", WithHTTPClient(httpClient))

		ctx := context.Background()
		request, err := client.NewRequest(ctx, "GET", "/test", nil)
		if err != nil {
			t.Fatalf("failed to create request: %s", err)
		}

		response := &MockBody{}

		// Act
		err = client.Send(request, response)
		assert.Equal(t, err, nil)

		// Assert
		assert.Equal(t, err, nil)
		assert.Equal(t, response.Param, "test")
	})
}
