package server

import (
	"go-auth/internal/app"
	"go-auth/internal/cart"
	"go-auth/internal/config"
	"go-auth/internal/middleware"
	"go-auth/internal/order"
	"go-auth/internal/payment"
	"go-auth/internal/product"
	"go-auth/internal/user"
	"time"

	//"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(a *app.App, cfg config.Config) *gin.Engine {

	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// user
	userRep := user.NewRepo(a.DB)
	userServ := user.NewService(userRep, a.Config.JWTSecret)
	userHand := user.NewHandler(userServ)

	// product
	prodRepo := product.NewRepo(a.DB)
	prodService := product.NewService(prodRepo)
	prodHandler := product.NewHandler(prodService)

	// cart
	cartRepo := cart.NewRepo(a.Redis, prodRepo)
	cartHandler := cart.NewHandler(cartRepo)

	// order
	orderRepo := order.NewRepo(a.DB)
	orderService := order.NewService(orderRepo, cartRepo, prodRepo)
	orderHander := order.NewHandler(orderService)

	// payment
	paymentRepo := payment.NewRepo(a.DB)
	paymentService := payment.NewService(paymentRepo, orderRepo, orderService)
	paymentHandler := payment.NewHandler(paymentService)

	// public route
	r.GET("/health", health)
	r.POST("/register", userHand.RegisterUser)
	r.POST("/login", userHand.LoginUser)

	r.GET("/products", prodHandler.GetAllProduct)
	r.GET("/products/:id", prodHandler.GetProductByID)
	r.POST("/payment/webhook", paymentHandler.HandleWebhook)

	authenticated := r.Group("")

	authenticated.Use(middleware.AuthRequired(cfg.JWTSecret))

	authenticated.POST("/cart", cartHandler.AddProductToCart)
	authenticated.GET("/cart/all", cartHandler.GetCartProduct)
	authenticated.DELETE("/cart/:productId", cartHandler.DeleteAnItemFromCart)
	authenticated.DELETE("/cart/delete", cartHandler.DeleteAllItemFromCart)

	authenticated.POST("/order/create", orderHander.CreateNewOrder)
	authenticated.DELETE("/order/delete/:id", orderHander.DeleteOrder)
	authenticated.GET("/orders", orderHander.GetAllUserOrder)
	authenticated.GET("/order/:id", orderHander.GetUserOrder)
	//authenticated.PUT("order/update/:id", orderHander.UpdateStatusById)

	authenticated.POST("/payment/initiate/:id", paymentHandler.InitiatePayment)
	authenticated.GET("/payment/verify/:reference", paymentHandler.VerifyPayment)
	authenticated.GET("/payments", paymentHandler.GetAllPayments)
	authenticated.GET("/payment/:id", paymentHandler.GetPaymentByID)
	//authenticated.PUT("/payment/update/:id", paymentHandler.UpdatePaymentStatus)

	// authenticated.GET("/profile", func(c *gin.Context) {
	// 	userId, _ := middleware.GetUserId(c)

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"ok":     true,
	// 		"userId": userId,
	// 		"data":   []any{},
	// 	})
	// })

	// admin routes

	ad := authenticated.Group("/admin")

	ad.Use(middleware.RequireAdmin())

	// ad.GET("", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"ok":   true,
	// 		"data": []any{},
	// 	})
	// })

	ad.POST("/product", prodHandler.CreateProduct)
	ad.PUT("/update-product/:id", prodHandler.UpdateProductByID)
	ad.DELETE("/delete-product/:id", prodHandler.DeleteProductByID)
	ad.GET("/all/orders", orderHander.GellAllOrder)
	ad.GET("/all/payments", paymentHandler.GellAllPayments)
	// profile route -> must be auth with Bearer tkn
	// authenticated.Use(middleware.AuthRequired(cfg.JWTSecret))
	// {
	// 	authenticated.GET("/profile", profile)

	// 	ad := authenticated.Group("/admin")

	// 	ad.Use(middleware.RequireAdmin())
	// 	{
	// 		ad.GET("/get-admin", admin)
	// 	}

	// }

	return r
}
