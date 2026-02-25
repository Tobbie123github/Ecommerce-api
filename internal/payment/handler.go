package payment

import (
	"go-auth/internal/middleware"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) InitiatePayment(c *gin.Context) {

	orderId := c.Param("id")

	userId, ok := middleware.GetUserId(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error grtting user id",
		})
		return
	}

	clientSecret, err := h.svc.CreatePayment(c.Request.Context(), userId, orderId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"clientSecret": clientSecret,
	})
}

func (h *Handler) HandleWebhook(c *gin.Context) {

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	signature := c.GetHeader("Stripe-Signature")

	if err := h.svc.HandleWebhook(c.Request.Context(), payload, signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *Handler) GetAllPayments(c *gin.Context) {

	userId, ok := middleware.GetUserId(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	payment, err := h.svc.repo.GetAll(c.Request.Context(), userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payments": payment})

}

func (h *Handler) GetPaymentByID(c *gin.Context) {

	paymentId := c.Param("id")

	userId, ok := middleware.GetUserId(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	payment, err := h.svc.repo.GetById(c.Request.Context(), userId, paymentId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment})

}



