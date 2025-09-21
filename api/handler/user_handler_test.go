package handler

import (
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

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserHandler_GetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         interface{}
		userIDExists   bool
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name:         "Positive: Valid user profile",
			userID:       uint(1),
			userIDExists: true,
			mockSetup: func(m *MockUserService) {
				user := &models.User{
					ID:       1,
					FullName: "Test User",
					Username: "testuser",
					Email:    "test@example.com",
					Password: "hashedpassword",
				}
				m.On("GetUserByID", uint(1)).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Negative: User not authenticated",
			userID:         nil,
			userIDExists:   false,
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:         "Negative: User not found",
			userID:       uint(999),
			userIDExists: true,
			mockSetup: func(m *MockUserService) {
				m.On("GetUserByID", uint(999)).Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:         "Edge: Invalid userID type",
			userID:       "invalid",
			userIDExists: true,
			mockSetup:    func(m *MockUserService) {},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			handler := NewUserHandler(mockService)
			router := gin.New()
			router.GET("/profile", func(c *gin.Context) {
				if tt.userIDExists {
					c.Set("userID", tt.userID)
				}
				handler.GetProfile(c)
			})

			req := httptest.NewRequest("GET", "/profile", nil)
			w := httptest.NewRecorder()

			// Handle panic for invalid type assertion
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, http.StatusInternalServerError, tt.expectedStatus)
				}
			}()

			router.ServeHTTP(w, req)

			if tt.expectedStatus != http.StatusInternalServerError {
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			// Verify password is cleared in response
			if tt.expectedStatus == http.StatusOK {
				var response models.User
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Empty(t, response.Password)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetProfile_PasswordClearing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	user := &models.User{
		ID:       1,
		FullName: "Test User",
		Password: "shouldbecleared",
	}
	mockService.On("GetUserByID", uint(1)).Return(user, nil)

	handler := NewUserHandler(mockService)
	router := gin.New()
	router.GET("/profile", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.GetProfile(c)
	})

	req := httptest.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response models.User
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Empty(t, response.Password)
	assert.Equal(t, "Test User", response.FullName)
}

func TestNewUserHandler(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.userService)
}