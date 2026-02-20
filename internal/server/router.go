package server

import (
	"go-auth/internal/app"
	"go-auth/internal/config"
	"go-auth/internal/middleware"
	"go-auth/internal/product"
	"go-auth/internal/user"
	//"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(a *app.App, cfg config.Config) *gin.Engine {

	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	userRep := user.NewRepo(a.DB)
	userServ := user.NewService(userRep, a.Config.JWTSecret)
	userHand := user.NewHandler(userServ)

	prodRepo := product.NewRepo(a.DB)
	prodService := product.NewService(prodRepo)
	prodHandler := product.NewHandler(prodService)

	// public route
	r.GET("/health", health)
	r.POST("/register", userHand.RegisterUser)
	r.POST("/login", userHand.LoginUser)


	
	r.GET("/products", prodHandler.GetAllProduct)
	r.GET("/:id", prodHandler.GetProductByID)





	authenticated := r.Group("")

	authenticated.Use(middleware.AuthRequired(cfg.JWTSecret))

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

	ad.POST("/create-product", prodHandler.CreateProduct)

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
