package payment

import (
	"log/slog"
	_ "mamlaka/internal/app/common"

	"github.com/labstack/echo/v4"
)

type PaymentHandler interface {
	MakePayment(c echo.Context) error
	GetPaymentDetail(c echo.Context) error
	Transactions(c echo.Context) error
}

type paymentHandler struct {
	logger         *slog.Logger
	paymentService PaymentService
}

// MakePayment godoc
// @Summary make a payment
// @Description  makes a single payment transaction
// @Tags Payment
// @Accept  json
// @Produce  json
// @Param   PaymentRequestDto body PaymentRequestDto true "Make Payment Request"
// @Success 200 {object} common.BaseResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router  /payments/payment [post]
func (p paymentHandler) MakePayment(c echo.Context) error {
	return p.paymentService.MakePayment(c)
}

// GetPaymentDetail godoc
// @Summary Get Details of a payment made
// @Description Authenticates a payment and returns a JWT token
// @Tags Payments
// @Accept  json
// @Produce  json
// @Success 200 {object} common.BaseResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router  /payments/{id} [get]
func (p paymentHandler) GetPaymentDetail(c echo.Context) error {
	return p.paymentService.GetPaymentDetail(c)
}

// Transactions godoc
// @Summary Transactions for all users
// @Description Gets all Transactions
// @Tags Payments
// @Accept  json
// @Produce  json
// @Success 201 {object} common.BaseResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router  /payments/transactions [get]
func (p paymentHandler) Transactions(c echo.Context) error {
	return p.paymentService.GetAllTransactions(c)
}

func NewPaymentHandler(logger *slog.Logger, service PaymentService) PaymentHandler {
	return paymentHandler{logger: logger, paymentService: service}
}
