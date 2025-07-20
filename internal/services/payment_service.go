package services

import (
	"fmt"
	"log"
	"time"

	"github.com/bagussubagja/backend-payment-gateway-go/internal/models"
	repository "github.com/bagussubagja/backend-payment-gateway-go/internal/repositories"
	"github.com/midtrans/midtrans-go"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(req *models.CreatePaymentRequest, user *models.User) (*models.CreatePaymentResponse, error)
	GetPaymentStatus(orderID string) (*models.Transaction, error)
	HandleNotification(notificationPayload map[string]interface{}) error
	GetPaymentHistory(userID uint) ([]models.Transaction, error)
	CreateQrisPayment(req *models.CreateQrisPaymentRequest, user *models.User) (*models.CreateQrisPaymentResponse, error)
}

type paymentService struct {
	txRepo      repository.TransactionRepository
	midtransSvc MidtransService
}

func NewPaymentService(txRepo repository.TransactionRepository, midtransSvc MidtransService) PaymentService {
	return &paymentService{txRepo, midtransSvc}
}

func (s *paymentService) CreateQrisPayment(req *models.CreateQrisPaymentRequest, user *models.User) (*models.CreateQrisPaymentResponse, error) {
	orderID := fmt.Sprintf("QRIS-%d", time.Now().UnixNano())

	var totalAmount int64
	var midtransItems []midtrans.ItemDetails

	for _, item := range req.Items {
		totalAmount += item.Price * int64(item.Quantity)
		midtransItems = append(midtransItems, midtrans.ItemDetails{
			ID:    item.ID,
			Price: item.Price,
			Qty:   item.Quantity,
			Name:  item.Name,
		})
	}

	midtransResp, midtransErr := s.midtransSvc.CreateQrisTransaction(orderID, totalAmount, midtransItems, user)
	if midtransErr != nil {
		return nil, midtransErr
	}

	dbTransactionErr := s.txRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		newTx := &models.Transaction{
			ID:                    orderID,
			UserID:                user.ID,
			Amount:                totalAmount,
			Status:                "pending",
			MidtransTransactionID: midtransResp.TransactionID,
			PaymentURL:            midtransResp.Actions[0].URL,
		}
		if err := tx.Create(newTx).Error; err != nil {
			return err
		}

		for _, item := range req.Items {
			txItem := models.TransactionItem{
				TransactionID: orderID,
				ItemID:        item.ID,
				Name:          item.Name,
				Price:         item.Price,
				Quantity:      item.Quantity,
			}
			if err := tx.Create(&txItem).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if dbTransactionErr != nil {
		log.Printf("ERROR: Gagal menyimpan transaksi QRIS ke DB: %v", dbTransactionErr)
		return nil, dbTransactionErr
	}

	log.Printf("SUKSES: Transaksi QRIS dengan Order ID: %s berhasil disimpan ke DB.", orderID)

	return &models.CreateQrisPaymentResponse{
		OrderID:    orderID,
		QrCodeUrl:  midtransResp.Actions[0].URL,
		ExpiryTime: midtransResp.ExpiryTime,
	}, nil
}

func (s *paymentService) CreatePayment(req *models.CreatePaymentRequest, user *models.User) (*models.CreatePaymentResponse, error) {

	orderID := fmt.Sprintf("ORDER-%d", time.Now().UnixNano())

	var totalAmount int64
	var midtransItems []midtrans.ItemDetails

	for _, item := range req.Items {
		totalAmount += item.Price * int64(item.Quantity)
		midtransItems = append(midtransItems, midtrans.ItemDetails{
			ID:    item.ID,
			Price: item.Price,
			Qty:   item.Quantity,
			Name:  item.Name,
		})
	}

	customer := midtrans.CustomerDetails{
		FName: req.CustomerDetails.FirstName,
		LName: req.CustomerDetails.LastName,
		Email: req.CustomerDetails.Email,
		Phone: req.CustomerDetails.Phone,
		BillAddr: &midtrans.CustomerAddress{
			FName:       req.CustomerDetails.FirstName,
			LName:       req.CustomerDetails.LastName,
			Phone:       req.CustomerDetails.Phone,
			Address:     req.CustomerDetails.Address,
			City:        req.CustomerDetails.City,
			Postcode:    req.CustomerDetails.PostalCode,
			CountryCode: "IDN",
		},
	}

	midtransResp, midtransErr := s.midtransSvc.CreateTransaction(orderID, totalAmount, midtransItems, customer)
	if midtransErr != nil {
		return nil, midtransErr
	}

	newTx := &models.Transaction{
		ID:                    orderID,
		UserID:                user.ID,
		Amount:                totalAmount,
		Status:                "pending",
		MidtransTransactionID: midtransResp.Token,
		PaymentURL:            midtransResp.RedirectURL,
	}

	var txItems []models.TransactionItem
	for _, item := range req.Items {
		txItems = append(txItems, models.TransactionItem{
			TransactionID: orderID,
			ItemID:        item.ID,
			Name:          item.Name,
			Price:         item.Price,
			Quantity:      item.Quantity,
		})
	}

	dbTransactionErr := s.txRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newTx).Error; err != nil {
			return err
		}
		if err := tx.Create(&txItems).Error; err != nil {
			return err
		}
		return nil
	})

	if dbTransactionErr != nil {
		log.Printf("Failed save to DB: %v", dbTransactionErr)
		return nil, dbTransactionErr
	}

	return &models.CreatePaymentResponse{
		OrderID:       orderID,
		RedirectURL:   midtransResp.RedirectURL,
		TransactionID: midtransResp.Token,
	}, nil
}

func (s *paymentService) GetPaymentStatus(orderID string) (*models.Transaction, error) {
	return s.txRepo.FindByID(orderID)
}

func (s *paymentService) GetPaymentHistory(userID uint) ([]models.Transaction, error) {
	return s.txRepo.FindByUserID(userID)
}

func (s *paymentService) HandleNotification(payload map[string]interface{}) error {
	orderID, _ := payload["order_id"].(string)
	transactionStatus, _ := payload["transaction_status"].(string)
	fraudStatus, _ := payload["fraud_status"].(string)

	tx, err := s.txRepo.FindByID(orderID)
	if err != nil {
		return fmt.Errorf("transaction not found: %s", orderID)
	}

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			tx.Status = "challenge"
		} else if fraudStatus == "accept" {
			tx.Status = "success"
		}
	} else if transactionStatus == "settlement" {
		tx.Status = "success"
	} else if transactionStatus == "deny" || transactionStatus == "expire" || transactionStatus == "cancel" {
		tx.Status = "failed"
	} else {
		tx.Status = transactionStatus
	}

	return s.txRepo.Update(tx)
}
