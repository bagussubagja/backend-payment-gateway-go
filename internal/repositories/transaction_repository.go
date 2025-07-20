package repository

import (
	"github.com/bagussubagja/backend-payment-gateway-go/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindByID(id string) (*models.Transaction, error)
	Update(transaction *models.Transaction) error
	FindByUserID(userID uint) ([]models.Transaction, error)
	GetDB() *gorm.DB
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByID(id string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Items").Where("id = ?", id).First(&transaction).Error
	return &transaction, err
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) FindByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Items").Where("user_id = ?", userID).Order("created_at desc").Find(&transactions).Error
	return transactions, err
}
