package services

import (
	"github.com/bagussubagja/backend-payment-gateway-go/config"
	"github.com/bagussubagja/backend-payment-gateway-go/internal/models"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService interface {
	CreateTransaction(orderID string, grossAmount int64, items []midtrans.ItemDetails, customer midtrans.CustomerDetails) (*snap.Response, *midtrans.Error)
	GetTransactionStatus(orderID string) (*coreapi.TransactionStatusResponse, error)
	CreateQrisTransaction(orderID string, amount int64, items []midtrans.ItemDetails, user *models.User) (*coreapi.ChargeResponse, *midtrans.Error)
}

type midtransService struct {
	snapApi snap.Client
	coreApi coreapi.Client
}

func NewMidtransService(cfg *config.Config) MidtransService {
	var snapClient snap.Client
	snapClient.New(cfg.MidtransServerKey, cfg.MidtransEnvironment)

	var coreClient coreapi.Client
	(&coreClient).New(cfg.MidtransServerKey, cfg.MidtransEnvironment)

	return &midtransService{
		snapApi: snapClient,
		coreApi: coreClient,
	}
}

func (s *midtransService) CreateTransaction(orderID string, grossAmount int64, items []midtrans.ItemDetails, customer midtrans.CustomerDetails) (*snap.Response, *midtrans.Error) {
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: grossAmount,
		},
		CustomerDetail: &customer,
		Items:          &items,
	}

	return s.snapApi.CreateTransaction(snapReq)
}

func (s *midtransService) GetTransactionStatus(orderID string) (*coreapi.TransactionStatusResponse, error) {
	return s.coreApi.CheckTransaction(orderID)
}

func (s *midtransService) CreateQrisTransaction(orderID string, amount int64, items []midtrans.ItemDetails, user *models.User) (*coreapi.ChargeResponse, *midtrans.Error) {
	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeQris,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		Items: &items,
		CustomerDetails: &midtrans.CustomerDetails{
			FName: user.FullName,
			Email: user.Email,
			Phone: user.PhoneNumber,
		},
	}
	return s.coreApi.ChargeTransaction(chargeReq)
}
