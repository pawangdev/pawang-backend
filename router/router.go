package router

import (
	"pawang-backend/config"
	"pawang-backend/handler"
	"pawang-backend/middleware"
	"pawang-backend/repository"
	"pawang-backend/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(router fiber.Router) {
	db, _ := config.Database()
	newUserRouter(router, db)
	newWalletRouter(router, db)
	newCategoryRouter(router, db)
	newSubCategoryRouter(router, db)
	newTransactionRouter(router, db)
	newStaticRouter(router)
}

func newUserRouter(router fiber.Router, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	userNotificationRepository := repository.NewUserNotificationRepository(db)
	walletRepository := repository.NewWalletRepository(db)
	userService := service.NewUserService(userRepository, walletRepository)
	authService := service.NewAuthService()
	emailService := service.NewEmailService()
	notificationService := service.NewNotificationService(userNotificationRepository)
	userHandler := handler.NewUserHandler(userService, authService, emailService, notificationService)

	router.Post("/users/register", userHandler.RegisterUser)
	router.Post("/users/login", userHandler.LoginUser)
	router.Put("/users/change-password", middleware.Authenticated(), userHandler.ChangePassword)
	router.Put("/users/change-profile", middleware.Authenticated(), userHandler.ChangeProfile)
	router.Get("/users/profile", middleware.Authenticated(), userHandler.UserProfile)
	router.Post("/users/reset-password", userHandler.RequestResetPasswordToken)
	router.Post("/users/reset-password/token", userHandler.VerifyResetPasswordToken)
	router.Post("/users/reset-password/password-confirmation", userHandler.ResetPasswordConfirmation)
}

func newWalletRouter(router fiber.Router, db *gorm.DB) {
	walletRepository := repository.NewWalletRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	walletService := service.NewWalletService(walletRepository, transactionRepository)
	authService := service.NewAuthService()
	walletHandler := handler.NewWalletHandler(walletService, authService)

	router.Get("/wallets", middleware.Authenticated(), walletHandler.GetWallets)
	router.Post("/wallets/create", middleware.Authenticated(), walletHandler.CreateWallet)
	router.Put("/wallets/update/:walletId", middleware.Authenticated(), walletHandler.UpdateWallet)
	router.Delete("/wallets/delete/:walletId", middleware.Authenticated(), walletHandler.DeleteWallet)
}

func newCategoryRouter(router fiber.Router, db *gorm.DB) {
	categoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepository)
	authService := service.NewAuthService()
	categoryHandler := handler.NewCategoryHandler(categoryService, authService)

	router.Get("/categories", middleware.Authenticated(), categoryHandler.GetCategories)
	// router.Post("/categories/create", middleware.Authenticated(), categoryHandler.CreateCategory)
	// router.Put("/categories/update/:categoryId", middleware.Authenticated(), categoryHandler.UpdateCategory)
	// router.Delete("/categories/delete/:categoryId", middleware.Authenticated(), categoryHandler.DeleteCategory)
}

func newSubCategoryRouter(router fiber.Router, db *gorm.DB) {
	subCategoryRepository := repository.NewSubCategoryRepository(db)
	subCategoryService := service.NewSubCategoryService(subCategoryRepository)
	authService := service.NewAuthService()
	subCategoryHandler := handler.NewSubCategoryHandler(subCategoryService, authService)

	router.Post("/categories/sub-categories/create", middleware.Authenticated(), subCategoryHandler.CreateSubCategory)
	router.Put("/categories/sub-categories/update/:subcategoryId", middleware.Authenticated(), subCategoryHandler.UpdateSubCategory)
	router.Delete("/categories/sub-categories/delete/:subcategoryId", middleware.Authenticated(), subCategoryHandler.DeleteSubCategory)
}

func newTransactionRouter(router fiber.Router, db *gorm.DB) {
	transactionRepository := repository.NewTransactionRepository(db)
	walletRepository := repository.NewWalletRepository(db)
	subCategoryRepository := repository.NewSubCategoryRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	transactionService := service.NewTransactionService(transactionRepository, walletRepository, subCategoryRepository, categoryRepository)
	authService := service.NewAuthService()
	transactionHandler := handler.NewTransactionHandler(transactionService, authService)

	router.Get("/transactions", middleware.Authenticated(), transactionHandler.GetTransactions)
	router.Post("/transactions/create", middleware.Authenticated(), transactionHandler.CreateTransaction)
	router.Put("/transactions/update/:transactionId", middleware.Authenticated(), transactionHandler.UpdateTransaction)
	router.Delete("/transactions/delete/:transactionId", middleware.Authenticated(), transactionHandler.DeleteTransaction)
}

func newStaticRouter(router fiber.Router) {
	router.Static("/storage", "./public")
}
