package server

import (
	"github.com/fenek-dev/go-outline-bot/pkg/payment_service"
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) PaymentWebhookHandler(c *gin.Context) {
	var request payment_service.WebhookRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, payment_service.WebhookResponse{OK: false})
		return
	}

	if request.Status == payment_service.TransactionStatusCancelled || request.Status == payment_service.TransactionStatusDeclined {
		err := s.service.CancelDeposit(c.Request.Context(), request.UUID)
		if err != nil {
			s.log.Error("can not cancel deposit %s: %e", request.UUID, err)
			c.JSON(500, payment_service.WebhookResponse{OK: false})
			return
		}
	}

	if request.Status == payment_service.TransactionStatusCompleted {
		err := s.service.ConfirmDeposit(c.Request.Context(), request.UUID)
		if err != nil {
			s.log.Error("can not confirm deposit %s: %e", request.UUID, err)
			c.JSON(500, payment_service.WebhookResponse{OK: false})
			return
		}
	}

	c.JSON(204, payment_service.WebhookResponse{OK: true})
	return
}
