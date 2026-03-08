package order

import (
	"go-auth/internal/middleware"
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

func (h *Handler) CreateNewOrder(c *gin.Context) {
    userId, ok := middleware.GetUserId(c)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var delivery DeliveryInfo
    if err := c.ShouldBindJSON(&delivery); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "delivery info required"})
        return
    }

    order, err := h.svc.CreateOrder(c.Request.Context(), userId, delivery)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Order created successfully",
        "order":   order,
    })
}

func (h *Handler) GetUserOrder(c *gin.Context) {

	orderId := c.Param("id")

	order, err := h.svc.repo.GetById(c.Request.Context(), orderId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})

}

func (h *Handler) GetAllUserOrder(c *gin.Context) {

	userId, ok := middleware.GetUserId(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error grtting user id",
		})

		return
	}

	order, err := h.svc.repo.GetByUserId(c.Request.Context(), userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})
}

func (h *Handler) DeleteOrder(c *gin.Context) {

	orderId := c.Param("id")

	userId, ok := middleware.GetUserId(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Error grtting user id",
		})

		return
	}

	if err := h.svc.DeleteOrder(c.Request.Context(), orderId, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order Deleted successfully",
	})
}

func (h *Handler) GellAllOrder(c *gin.Context){

	orders, err := h.svc.repo.Orders(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
	})


}

// func (h *Handler) UpdateStatusById(c *gin.Context){

// 	orderId := c.Param("id")
// 	userId, ok := middleware.GetUserId(c)

// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error": "Error grtting user id",
// 		})

// 		return
// 	}

// 	if err := h.svc.UpdateOrder(c.Request.Context(), orderId, userId, "status"); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error":"error with updating order",
// 		})
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success":"order updated succesfully",
// 	})
// }
