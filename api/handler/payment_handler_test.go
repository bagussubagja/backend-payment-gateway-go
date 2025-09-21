package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bagussubagja/backend-payment-gateway-go/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) CreatePayment(req *models.CreatePaymentRequest, user *models.User) (*models.CreatePaymentResponse, error) {
	args := m.Called(req, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CreatePaymentResponse), args.Error(1)
}

func (m *MockPaymentService) GetPaymentStatus(orderID string) (*models.Transaction, error) {
	args := m.Called(orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockPaymentService) GetPaymentHistory(userID uint) ([]models.Transaction, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockPaymentService) HandleNotification(payload map[string]interface{}) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *MockPaymentService) CreateQrisPayment(req *models.CreateQrisPaymentRequest, user *models.User) (*models.CreateQrisPaymentResponse, error) {
	args := m.Called(req, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CreateQrisPaymentResponse), args.Error(1)
}

func TestPaymentHandler_CreatePayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		userID         interface{}
		userIDExists   bool
		mockSetup      func(*MockPaymentService, *MockUserService)
		expectedStatus int
	}{
		{
			name: "Positive: Valid payment creation",
			requestBody: models.CreatePaymentRequest{
				Items: []models.ItemDetailRequest{{ID: "1", Name: "Test", Price: 1000, Quantity: 1}},
				CustomerDetails: models.AddressDetail{
					FirstName: "Test", Email: "test@example.com", Phone: "123", Address: "Test", City: "Test", PostalCode: "12345",
				},
			},
			userID:       uint(1),
			userIDExists: true,
			mockSetup: func(mp *MockPaymentService, mu *MockUserService) {
				user := &models.User{ID: 1, FullName: "Test User"}
				mu.On("GetUserByID", uint(1)).Return(user, nil)
				mp.On("CreatePayment", mock.AnythingOfType("*models.CreatePaymentRequest"), user).Return(&models.CreatePaymentResponse{OrderID: "order123"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Negative: Invalid JSON",
			requestBody:    "invalid",
			userID:         uint(1),
			userIDExists:   true,
			mockSetup:      func(mp *MockPaymentService, mu *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Negative: User not authenticated",
			requestBody: models.CreatePaymentRequest{
				Items: []models.ItemDetailRequest{{ID: "1", Name: "Test", Price: 1000, Quantity: 1}},
				CustomerDetails: models.AddressDetail{
					FirstName: "Test", Email: "test@example.com", Phone: "123", Address: "Test", City: "Test", PostalCode: "12345",
				},
			},
			userIDExists:   false,
			mockSetup:      func(mp *MockPaymentService, mu *MockUserService) {},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPaymentService := new(MockPaymentService)
			mockUserService := new(MockUserService)
			tt.mockSetup(mockPaymentService, mockUserService)

			handler := NewPaymentHandler(mockPaymentService, mockUserService)
			router := gin.New()
			router.POST("/payment", func(c *gin.Context) {
				if tt.userIDExists {
					c.Set("userID", tt.userID)
				}
				handler.CreatePayment(c)
			})

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/payment", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockPaymentService.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestPaymentHandler_GetStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		orderID        string
		userID         uint
		mockSetup      func(*MockPaymentService)
		expectedStatus int
	}{
		{
			name:    "Positive: Valid status request",
			orderID: "order123",
			userID:  1,
			mockSetup: func(m *MockPaymentService) {
				transaction := &models.Transaction{ID: "order123", UserID: 1, Status: "pending"}
				m.On("GetPaymentStatus", "order123").Return(transaction, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Negative: Transaction not found",
			orderID: "notfound",
			userID:  1,
			mockSetup: func(m *MockPaymentService) {
				m.On("GetPaymentStatus", "notfound").Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:    "Negative: Unauthorized access",
			orderID: "order123",
			userID:  2,
			mockSetup: func(m *MockPaymentService) {
				transaction := &models.Transaction{ID: "order123", UserID: 1, Status: "pending"}
				m.On("GetPaymentStatus", "order123").Return(transaction, nil)
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockPaymentService)
			tt.mockSetup(mockService)

			handler := NewPaymentHandler(mockService, nil)
			router := gin.New()
			router.GET("/status/:orderID", func(c *gin.Context) {
				c.Set("userID", tt.userID)
				handler.GetStatus(c)
			})

			req := httptest.NewRequest("GET", "/status/"+tt.orderID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestPaymentHandler_GetHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userIDExists   bool
		userID         uint
		mockSetup      func(*MockPaymentService)
		expectedStatus int
	}{
		{
			name:         "Positive: Valid history request",
			userIDExists: true,
			userID:       1,
			mockSetup: func(m *MockPaymentService) {
				history := []models.Transaction{{ID: "order1", UserID: 1}}
				m.On("GetPaymentHistory", uint(1)).Return(history, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Negative: User not authenticated",
			userIDExists:   false,
			mockSetup:      func(m *MockPaymentService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:         "Negative: Service error",
			userIDExists: true,
			userID:       1,
			mockSetup: func(m *MockPaymentService) {
				m.On("GetPaymentHistory", uint(1)).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockPaymentService)
			tt.mockSetup(mockService)

			handler := NewPaymentHandler(mockService, nil)
			router := gin.New()
			router.GET("/history", func(c *gin.Context) {
				if tt.userIDExists {
					c.Set("userID", tt.userID)
				}
				handler.GetHistory(c)
			})

			req := httptest.NewRequest("GET", "/history", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestPaymentHandler_HandleNotification(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockPaymentService)
		expectedStatus int
	}{
		{
			name:        "Positive: Valid notification",
			requestBody: map[string]interface{}{"order_id": "order123", "transaction_status": "settlement"},
			mockSetup: func(m *MockPaymentService) {
				m.On("HandleNotification", mock.AnythingOfType("map[string]interface {}")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Negative: Invalid JSON",
			requestBody:    "invalid",
			mockSetup:      func(m *MockPaymentService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Negative: Service error",
			requestBody: map[string]interface{}{"order_id": "order123"},
			mockSetup: func(m *MockPaymentService) {
				m.On("HandleNotification", mock.AnythingOfType("map[string]interface {}")).Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockPaymentService)
			tt.mockSetup(mockService)

			handler := NewPaymentHandler(mockService, nil)
			router := gin.New()
			router.POST("/notification", handler.HandleNotification)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/notification", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestPaymentHandler_CreateQrisPayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		userID         uint
		mockSetup      func(*MockPaymentService, *MockUserService)
		expectedStatus int
	}{
		{
			name: "Positive: Valid QRIS payment",
			requestBody: models.CreateQrisPaymentRequest{
				Items: []models.ItemDetailRequest{{ID: "1", Name: "Test", Price: 1000, Quantity: 1}},
			},
			userID: 1,
			mockSetup: func(mp *MockPaymentService, mu *MockUserService) {
				user := &models.User{ID: 1, FullName: "Test User"}
				mu.On("GetUserByID", uint(1)).Return(user, nil)
				mp.On("CreateQrisPayment", mock.AnythingOfType("*models.CreateQrisPaymentRequest"), user).Return(&models.CreateQrisPaymentResponse{OrderID: "qris123"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Negative: Invalid JSON",
			requestBody:    "invalid",
			userID:         1,
			mockSetup:      func(mp *MockPaymentService, mu *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPaymentService := new(MockPaymentService)
			mockUserService := new(MockUserService)
			tt.mockSetup(mockPaymentService, mockUserService)

			handler := NewPaymentHandler(mockPaymentService, mockUserService)
			router := gin.New()
			router.POST("/qris", func(c *gin.Context) {
				c.Set("userID", tt.userID)
				handler.CreateQrisPayment(c)
			})

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/qris", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockPaymentService.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestNewPaymentHandler(t *testing.T) {
	mockPaymentService := new(MockPaymentService)
	mockUserService := new(MockUserService)
	handler := NewPaymentHandler(mockPaymentService, mockUserService)
	
	assert.NotNil(t, handler)
	assert.Equal(t, mockPaymentService, handler.paymentService)
	assert.Equal(t, mockUserService, handler.userService)
}