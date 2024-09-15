package payment

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type PaymentRepository interface {
	GetPaymentInfoByEmail(email string) (*Payment, error)
	GetPaymentByID(id uint) (*Payment, error)
	CreatePayment(payment *Payment) (*Payment, error)
	GetPayments() ([]Payment, error)
}

type paymentRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

//func (p paymentRepository) CreatePayment(payment *Payment) (*Payment, error) {
//	if err := p.DB.Create(payment).Error; err != nil {
//		p.logger.Error("Error creating payment", err)
//		return nil, err
//	}
//	p.logger.Info("Payment created successfully", "paymentID", payment.ID)
//	return payment, nil
//}

func (p paymentRepository) CreatePayment(payment *Payment) (*Payment, error) {
	// Define the time range to check for recent payments
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)

	// Log the time range being used for checking
	p.logger.Info("Checking for duplicate payments within the last minute", "oneMinuteAgo", oneMinuteAgo)

	// Query to check for existing payments with the same amount, phone number, and card number within the last minute
	var existingPayment Payment
	err := p.DB.Joins("JOIN payment_details ON payment_details.payment_id = payments.id").
		Where("amount = ? AND payment_details.phone_number = ? AND payment_details.card_number = ? AND payments.created_at >= ?",
			payment.Amount, payment.PaymentDetails.PhoneNumber, payment.PaymentDetails.CardNumber, oneMinuteAgo).
		First(&existingPayment).Error

	// Log the result of the query
	if err == nil {
		// A payment with the same criteria was found within the last minute
		p.logger.Error("Duplicate payment detected", "amount", payment.Amount, "phone_number", payment.PaymentDetails.PhoneNumber, "card_number", payment.PaymentDetails.CardNumber)
		return nil, fmt.Errorf("a payment with the amount %s, phone number %s, and card number %s has already been processed recently", payment.Amount, payment.PaymentDetails.PhoneNumber, payment.PaymentDetails.CardNumber)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// If there is an error other than "record not found", log and return the error
		p.logger.Error("Error querying for existing payments", "error", err)
		return nil, err
	}

	// Proceed with creating the new payment only if no similar payments were found
	if err := p.DB.Create(payment).Error; err != nil {
		p.logger.Error("Error creating payment", err)
		return nil, err
	}
	p.logger.Info("Payment created successfully", "paymentID", payment.ID)
	return payment, nil
}

func (p paymentRepository) GetPaymentInfoByEmail(email string) (*Payment, error) {
	var payment Payment
	if err := p.DB.Where("email = ?", email).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.logger.Info("Payment not found", "email", email)
			return nil, nil
		}
		p.logger.Error("Error fetching payment by email", err)
		return nil, err
	}
	p.logger.Info("Payment fetched successfully", "paymentID", payment.ID)
	return &payment, nil
}

func (p paymentRepository) GetPaymentByID(paymentID uint) (*Payment, error) {
	var payment Payment
	// Use Preload to eagerly load related entities
	if err := p.DB.Preload("PaymentDetails").
		Where("id = ?", paymentID).
		First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.logger.Info("Payment not found", "paymentID", paymentID)
			return nil, nil
		}
		p.logger.Error("Error fetching payment by ID", err)
		return nil, err
	}
	p.logger.Info("Payment fetched successfully", "paymentID", paymentID)
	return &payment, nil
}

func (p paymentRepository) GetPayments() ([]Payment, error) {
	var payments []Payment // Change to a slice to hold multiple payments
	if err := p.DB.Find(&payments).Error; err != nil {
		p.logger.Error("Error fetching payments", err)
		return nil, err
	}
	p.logger.Info("Payments fetched successfully", "count", len(payments))
	return payments, nil
}

func NewPaymentRepository(db *gorm.DB, logger *slog.Logger) PaymentRepository {
	return paymentRepository{
		DB:     db,
		logger: logger,
	}
}
