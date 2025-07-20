package models

type RegisterRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Address     string `json:"address" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	City        string `json:"city" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type ItemDetailRequest struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Price    int64  `json:"price" binding:"required"`
	Quantity int32  `json:"quantity" binding:"required"`
}

type AddressDetail struct {
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone" binding:"required"`
	Address    string `json:"address" binding:"required"`
	City       string `json:"city" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`
}

type CreatePaymentRequest struct {
	Items           []ItemDetailRequest `json:"items" binding:"required"`
	CustomerDetails AddressDetail       `json:"customer_details" binding:"required"`
}

type CreatePaymentResponse struct {
	OrderID       string `json:"order_id"`
	RedirectURL   string `json:"redirect_url"`
	TransactionID string `json:"transaction_id"`
}

type PaymentNotification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
}

type CreateQrisPaymentRequest struct {
	Items []ItemDetailRequest `json:"items" binding:"required"`
}

type CreateQrisPaymentResponse struct {
	OrderID    string `json:"order_id"`
	QrCodeUrl  string `json:"qr_code_url"`
	ExpiryTime string `json:"expiry_time"`
}
