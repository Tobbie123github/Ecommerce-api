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
	return &Handler{svc: svc}
}

func (h *Handler) InitiatePayment(c *gin.Context) {
	orderId := c.Param("id")

	userId, ok := middleware.GetUserId(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		Email string `json:"email"`
	}
	c.ShouldBindJSON(&body)

	authURL, reference, err := h.svc.CreatePayment(c.Request.Context(), userId, orderId, body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authorization_url": authURL,
		"reference":         reference,
	})
}

func (h *Handler) HandleWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// Paystack sends signature in x-paystack-signature header
	signature := c.GetHeader("x-paystack-signature")

	if err := h.svc.HandleWebhook(c.Request.Context(), payload, signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *Handler) VerifyPayment(c *gin.Context) {
	reference := c.Param("reference")

	result, err := h.svc.VerifyPayment(c.Request.Context(), reference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *Handler) GetAllPayments(c *gin.Context) {
	userId, ok := middleware.GetUserId(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payments, err := h.svc.repo.GetAll(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payments": payments})
}

func (h *Handler) GetPaymentByID(c *gin.Context) {
	paymentId := c.Param("id")

	userId, ok := middleware.GetUserId(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payment, err := h.svc.repo.GetById(c.Request.Context(), userId, paymentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment})
}


func (h *Handler) GellAllPayments(c *gin.Context){

	payments, err := h.svc.repo.Orders(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": payments,
	})

}
