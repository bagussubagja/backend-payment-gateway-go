package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestEntryHandler_GetEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Positive: GET request returns success",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]string{
				"status":  "success",
				"message": "Payment API Gateway with Midtrans Golang by Bagus Subagja",
			},
		},
		{
			name:           "Negative: POST method not allowed",
			method:         "POST",
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name:           "Negative: PUT method not allowed",
			method:         "PUT",
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name:           "Negative: DELETE method not allowed",
			method:         "DELETE",
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewEntryHandler()
			router := gin.New()
			router.GET("/", handler.GetEntry)

			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestEntryHandler_GetEntry_ResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	handler := NewEntryHandler()
	router := gin.New()
	router.GET("/", handler.GetEntry)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Test response headers
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	
	// Test JSON structure
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "status")
	assert.Contains(t, response, "message")
	assert.IsType(t, "", response["status"])
	assert.IsType(t, "", response["message"])
}

func TestEntryHandler_GetEntry_ConcurrentRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	handler := NewEntryHandler()
	router := gin.New()
	router.GET("/", handler.GetEntry)

	// Edge case: Test concurrent requests
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, http.StatusOK, w.Code)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestNewEntryHandler(t *testing.T) {
	handler := NewEntryHandler()
	assert.NotNil(t, handler)
	assert.IsType(t, &EntryHandler{}, handler)
}