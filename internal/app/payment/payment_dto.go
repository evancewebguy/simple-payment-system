package payment

type PaymentRequestDto struct {
	Amount         string            `json:"amount"`
	Currency       string            `json:"currency"`
	PaymentMethod  string            `json:"payment_method"` // Ensure this is a string for conversion
	PaymentDetails PaymentDetailsDto `json:"payment_details"`
}

type PaymentDetailsDto struct {
	CardNumber  string `json:"card_number"`
	ExpiryDate  string `json:"expiry_date"`
	CVV         string `json:"cvv"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type PaymentResponseDto struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}
