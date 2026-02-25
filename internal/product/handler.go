package product

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) CreateProduct(c *gin.Context) {

	var input ProductInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	product, err := h.svc.UploadItems(c.Request.Context(), input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product Added Successfully",
		"data":    product,
	})

}

func (h *Handler) GetAllProduct(c *gin.Context) {
	products, err := h.svc.repo.GetAll(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
	})
}

func (h *Handler) GetProductByID(c *gin.Context) {

	idStr := c.Param("id")

	objId, err := bson.ObjectIDFromHex(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	product, err := h.svc.repo.GetById(c.Request.Context(), objId)

	if err != nil {

		// lets use mongodb errs
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Product not Found",
			})
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func (h *Handler) UpdateProductByID(c *gin.Context) {


	idStr := c.Param("id")
	var update UpdateProduct
	if err := c.ShouldBind(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	prod, err := h.svc.UpdateItems(c.Request.Context(), update, idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data": prod,
	})

}



func (h *Handler) DeleteProductByID(c *gin.Context) {

	id := c.Param("id")


	err := h.svc.DeleteProduct(c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})

}
