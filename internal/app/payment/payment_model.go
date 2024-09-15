package payment

import "gorm.io/gorm"

type PaymentMethod string

const (
	CreditCard PaymentMethod = "credit_card"
	EWallet    PaymentMethod = "e_wallet"
	Mpesa      PaymentMethod = "mpesa"
)

// Payment represents a request to process a payment.
type Payment struct {
	gorm.Model
	ID             uint           `gorm:"primaryKey" json:"id"`
	Amount         string         `json:"amount"`
	Currency       string         `json:"currency"`
	PaymentMethod  PaymentMethod  `json:"payment_method"`
	PaymentDetails PaymentDetails `gorm:"foreignKey:PaymentID" json:"payment_details"` // One-to-One relationship
}

// PaymentDetails represents detailed payment information.
type PaymentDetails struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey" json:"id"`
	CardNumber  string `json:"card_number"`
	ExpiryDate  string `json:"expiry_date"`
	CVV         string `json:"cvv"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	PaymentID   uint   `json:"payment_id"` // Foreign key
}
