package payment

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log/slog"
	"mamlaka/internal/app/middlewares"
)

func RegisterPaymentRoutes(e *echo.Group, logger *slog.Logger, db *gorm.DB) {
	paymentRepository := NewPaymentRepository(db, logger)
	paymentService := NewPaymentService(logger, paymentRepository)
	paymentHandler := NewPaymentHandler(logger, paymentService)

	payment := e.Group("/payments")
	{
		payment.Use(middlewares.JWTMiddleware)

		payment.POST("/payment", paymentHandler.MakePayment)
		payment.GET("/:id", paymentHandler.GetPaymentDetail)
		payment.GET("/transactions", paymentHandler.Transactions)
	}
}
