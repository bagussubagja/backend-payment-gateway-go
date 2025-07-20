package models

import (
	"time"
)

type User struct {
	ID          uint   `gorm:"primaryKey"`
	FullName    string `gorm:"not null"`
	Username    string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	Address     string
	PhoneNumber string
	City        string
	PostalCode  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TransactionItem struct {
	ID            uint   `gorm:"primaryKey"`
	TransactionID string `gorm:"not null"`
	ItemID        string `gorm:"not null"`
	Name          string `gorm:"not null"`
	Price         int64  `gorm:"not null"`
	Quantity      int32  `gorm:"not null"`
}

type Transaction struct {
	ID                    string `gorm:"primaryKey"`
	UserID                uint   `gorm:"not null"`
	Amount                int64  `gorm:"not null"`
	Status                string `gorm:"not null"`
	MidtransTransactionID string
	PaymentURL            string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	User                  User              `gorm:"foreignKey:UserID"`
	Items                 []TransactionItem `gorm:"foreignKey:TransactionID"`
}
