package payment

import (
	"fmt"
	"log/slog"
	"mamlaka/internal/app/common"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// PaymentService defines the methods available in the payment service.
type PaymentService interface {
	MakePayment(c echo.Context) error
	GetPaymentDetail(c echo.Context) error
	GetAllTransactions(c echo.Context) error
}

// paymentService is the implementation of PaymentService.
type paymentService struct {
	logger     *slog.Logger
	repository PaymentRepository
}

// isValidPaymentMethod checks if the given payment method is valid.
func isValidPaymentMethod(method PaymentMethod) bool {
	validMethods := map[PaymentMethod]struct{}{
		CreditCard: {},
		EWallet:    {},
		Mpesa:      {},
	}

	_, exists := validMethods[method]
	return exists
}

// MakePayment handles payment requests
func (p paymentService) MakePayment(c echo.Context) error {
	var makePaymentRequest PaymentRequestDto
	if err := c.Bind(&makePaymentRequest); err != nil {
		p.logger.Error("Error parsing payment request body", err)
		return p.handleError(c, err, http.StatusBadRequest)
	}

	// Validate the incoming payment request
	if err := common.ValidateModel(makePaymentRequest); err != nil {
		p.logger.Error("Invalid payment request body", err)
		return p.handleError(c, err, http.StatusBadRequest)
	}

	// Convert payment method and validate
	paymentMethod := PaymentMethod(makePaymentRequest.PaymentMethod)
	if !isValidPaymentMethod(paymentMethod) { // Implement isValidPaymentMethod function
		err := fmt.Errorf("invalid payment method: %s", makePaymentRequest.PaymentMethod)
		p.logger.Error("Invalid payment method", err)
		return p.handleError(c, err, http.StatusBadRequest)
	}

	// Create and populate the Payment struct
	payment := &Payment{ // Use a pointer here
		Amount:        makePaymentRequest.Amount,
		Currency:      makePaymentRequest.Currency,
		PaymentMethod: paymentMethod,
		PaymentDetails: PaymentDetails{
			CardNumber:  makePaymentRequest.PaymentDetails.CardNumber,
			ExpiryDate:  makePaymentRequest.PaymentDetails.ExpiryDate,
			CVV:         makePaymentRequest.PaymentDetails.CVV,
			PhoneNumber: makePaymentRequest.PaymentDetails.PhoneNumber,
			Email:       makePaymentRequest.PaymentDetails.Email,
		},
	}

	// Simulate payment processing
	response, err := p.ProcessPayment(payment) // Pass pointer to ProcessPayment
	if err != nil {
		p.logger.Error("Error processing payment", err)
		return p.handleError(c, err, http.StatusInternalServerError)
	}

	// Save the payment to the database
	if _, err := p.repository.CreatePayment(payment); err != nil { // Pass pointer to CreatePayment
		p.logger.Error("Error creating payment", err)
		return p.handleError(c, err, http.StatusInternalServerError)
	}

	// Log success and return the response
	p.logger.Info("Payment successful")
	return c.JSON(http.StatusOK, response) // Return the PaymentResponseDto
}

// ProcessPayment simulates payment processing
func (p paymentService) ProcessPayment(payment *Payment) (*PaymentResponseDto, error) {
	// Simulate processing delay
	time.Sleep(2 * time.Second)

	// Here you would integrate with an actual payment gateway
	// For example purposes, we're just returning a simulated response

	// Check if payment amount is valid
	if payment.Amount == "" || payment.Currency == "" {
		return nil, fmt.Errorf("invalid payment details")
	}

	// Generate a simulated transaction ID
	transactionID := fmt.Sprintf("TXN-%d", time.Now().UnixNano())

	// Create a simulated response
	response := &PaymentResponseDto{
		TransactionID: transactionID,
		Status:        "Success",
		Message:       "Payment processed successfully",
	}

	return response, nil
}

// GetPaymentDetail is a helper function for creating error responses.
func (p paymentService) GetPaymentDetail(c echo.Context) error {

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Assuming you want a 32-bit unsigned integer
	if err != nil {
		p.logger.Error("Error parsing ID", err)
		return p.handleError(c, err, http.StatusBadRequest)
	}

	// Convert id to uint
	uintID := uint(id)

	// Retrieve all payments from the repository
	payment, err := p.repository.GetPaymentByID(uintID)
	if err != nil {
		p.logger.Error("Error fetching payments", err)
		return p.handleError(c, err, http.StatusInternalServerError)
	}

	p.logger.Info("Payment detail fetched successful")
	return c.JSON(http.StatusOK, common.BaseResponse{
		Status:  http.StatusOK,
		Message: "Payment detail fetched successful",
		Data:    payment,
	})
}

// GetAllTransactions handles payment payment requests
func (p paymentService) GetAllTransactions(c echo.Context) error {
	// Retrieve all payments from the repository
	payments, err := p.repository.GetPayments()
	if err != nil {
		p.logger.Error("Error fetching payments", err)
		return p.handleError(c, err, http.StatusInternalServerError)
	}

	// Log the successful retrieval of payments
	p.logger.Info("Payments fetched successfully", "count", len(payments))

	// Return a successful response with the fetched payments
	return c.JSON(http.StatusOK, common.BaseResponse{
		Status:  http.StatusOK,
		Message: "Payments fetched successfully",
		Data:    payments, // Include the payments in the response data
	})
}

// handleError is a helper function for creating error responses.
func (p paymentService) handleError(c echo.Context, err error, status int) error {
	return c.JSON(status, common.ErrorResponse{
		Status: status,
		Error:  err.Error(),
	})
}

// NewPaymentService creates a new instance of paymentService.
func NewPaymentService(logger *slog.Logger, repository PaymentRepository) PaymentService {
	return paymentService{logger: logger, repository: repository}
}
