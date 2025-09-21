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

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (uint, error) {
	args := m.Called(token)
	return args.Get(0).(uint), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Positive: Valid registration",
			requestBody: models.RegisterRequest{
				FullName:    "Test User",
				Username:    "testuser",
				Email:       "test@example.com",
				Password:    "password123",
				Address:     "Test Address",
				PhoneNumber: "1234567890",
				City:        "Test City",
				PostalCode:  "12345",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", mock.AnythingOfType("*models.RegisterRequest")).Return(&models.User{ID: 1}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "Negative: Invalid JSON",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			handler := NewAuthHandler(mockService)
			router := gin.New()
			router.POST("/register", handler.Register)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "Positive: Valid login",
			requestBody: models.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", mock.AnythingOfType("*models.LoginRequest")).Return(&models.LoginResponse{Token: "token123"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Negative: Invalid JSON",
			requestBody:    "invalid",
			mockSetup:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Negative: Invalid credentials",
			requestBody: models.LoginRequest{
				Username: "testuser",
				Password: "wrong",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", mock.AnythingOfType("*models.LoginRequest")).Return(nil, errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAuthService)
			tt.mockSetup(mockService)

			handler := NewAuthHandler(mockService)
			router := gin.New()
			router.POST("/login", handler.Login)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewAuthHandler(nil)
	router := gin.New()
	router.POST("/logout", handler.Logout)

	req := httptest.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Successfully logged out", response["message"])
}

func TestNewAuthHandler(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.authService)
}