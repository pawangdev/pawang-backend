package routes

import (
	"pawang-backend/controllers"
	"pawang-backend/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRouter(e *echo.Echo) {
	api := e.Group("/api")

	e.GET("/auth/:provider/authorize", controllers.LoginWithGoogle)
	e.GET("/auth/:provider/callback", controllers.LoginWithGoogleCallback)

	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)
	api.GET("/profile", controllers.Profile, middleware.IsAuthenticated)
	api.POST("/profile/change-password", controllers.ChangePassword, middleware.IsAuthenticated)
	api.POST("/profile/change-profile", controllers.UpdateProfile, middleware.IsAuthenticated)
	api.Static("/profile_image/", "public/users/")

	api.GET("/categories", controllers.CategoryIndex, middleware.IsAuthenticated)
	api.GET("/categories/:categoryId", controllers.CategoryShow, middleware.IsAuthenticated)
	api.POST("/categories/create", controllers.CategoryStore, middleware.IsAuthenticated)
	api.PUT("/categories/:categoryId/update", controllers.CategoryUpdate, middleware.IsAuthenticated)
	api.DELETE("/categories/:categoryId/delete", controllers.CategoryDestroy, middleware.IsAuthenticated)
	api.Static("/categories/", "public/users/")

	api.GET("/wallets", controllers.WalletIndex, middleware.IsAuthenticated)
	api.GET("/wallets/:walletId", controllers.WalletShow, middleware.IsAuthenticated)
	api.POST("/wallets/create", controllers.WalletStore, middleware.IsAuthenticated)
	api.PUT("/wallets/:walletId/update", controllers.WalletUpdate, middleware.IsAuthenticated)
	api.DELETE("/wallets/:walletId/delete", controllers.WalletDestroy, middleware.IsAuthenticated)

	api.GET("/transactions", controllers.TransactionIndex, middleware.IsAuthenticated)
	api.GET("/transactions/:transactionId", controllers.TransactionShow, middleware.IsAuthenticated)
	api.POST("/transactions/create", controllers.TransactionStore, middleware.IsAuthenticated)
	api.PUT("/transactions/:transactionId/update", controllers.TransactionUpdate, middleware.IsAuthenticated)
	api.DELETE("/transactions/:transactionId/delete", controllers.TransactionDestroy, middleware.IsAuthenticated)
	api.Static("/transactions/", "public/users/")
}
