package user

import (
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

func (h Handler) RegisterUser(c *gin.Context) {

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "invalid json",
		})

		return
	}

	out, err := h.svc.Register(c.Request.Context(), input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sign up Successful",
		"user":    out,
	})

}

func (h *Handler) LoginUser(c *gin.Context) {

	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "invalid json",
		})

		return
	}

	result, err := h.svc.Login(c.Request.Context(), input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Successful",
		"user": result,
	})
}