package user

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log/slog"
)

func RegisterUserRoutes(e *echo.Group, logger *slog.Logger, db *gorm.DB) {
	userRepository := NewUserRepository(db, logger)
	userService := NewUserService(logger, userRepository)
	userHandler := NewUserHandler(logger, userService)

	auth := e.Group("/auth")
	{
		auth.POST("/login", userHandler.Login)
		auth.POST("/register", userHandler.Register)
		//auth.POST("/verify", userHandler.VerifyAccount)
		//auth.POST("/initiate-reset", userHandler.SendResetToken)
		//auth.POST("/reset-password", userHandler.ResetPassword)
	}

	profile := e.Group("/user")
	{
		// profile.GET("/profile", userHandler.GetProfile)
		// profile.PUT("/profile", userHandler.UpdateProfile)
		profile.POST("/refresh-token", userHandler.RefreshToken)

	}
}
