package cart

import (
	"go-auth/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{
		repo: repo,
	}
}

type CartItemRequest struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
	Size      string `json:"size"`
}

func (h *Handler) AddProductToCart(c *gin.Context) {

	var req CartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	userId, _ := middleware.GetUserId(c)

	if req.UserID != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthroized",
		})
		return
	}

	cart, err := h.repo.AddCart(c.Request.Context(), req.ProductID, req.UserID, int(req.Quantity), req.Size)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"cart":   cart,
	})

}

func (h *Handler) GetCartProduct(c *gin.Context) {

	userId, ok := middleware.GetUserId(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": ok,
		})
	}

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	cart, err := h.repo.GetCart(c.Request.Context(), userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cart": cart,
	})

}

func (h *Handler) DeleteAnItemFromCart(c *gin.Context) {

	// var req CartItemRequest

	userId, ok := middleware.GetUserId(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": ok,
		})
	}

	productId := c.Param("productId")

	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// }

	if err := h.repo.DeleteCartItem(c.Request.Context(), productId, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Item Deleted from cart",
	})
}

func (h *Handler) DeleteAllItemFromCart(c *gin.Context) {

	// var req CartItemRequest

	//userId := c.Param("id")

	userId, ok := middleware.GetUserId(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": ok,
		})
	}

	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// }

	if err := h.repo.DeleteAll(c.Request.Context(), userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All items Deleted from cart",
	})
}
