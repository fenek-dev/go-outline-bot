package payment_service

import "time"

type WebhookRequest struct {
	ID        int               `json:"id"`
	UUID      string            `json:"uuid"`
	Amount    int               `json:"amount"`
	Status    TransactionStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
}

type WebhookResponse struct {
	OK bool `json:"ok"`
}
