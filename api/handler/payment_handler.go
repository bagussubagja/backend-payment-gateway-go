package handler

import (
	"fmt"
	"net/http"

	"github.com/bagussubagja/backend-payment-gateway-go/internal/models"
	"github.com/bagussubagja/backend-payment-gateway-go/internal/services"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService services.PaymentService
	userService    services.UserService
}

func NewPaymentHandler(paymentService services.PaymentService, userService services.UserService) *PaymentHandler {
	return &PaymentHandler{paymentService, userService}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Authenticated user not found"})
		return
	}

	resp, err := h.paymentService.CreatePayment(&req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process after payment creation",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) GetStatus(c *gin.Context) {
	orderID := c.Param("orderID")

	transaction, err := h.paymentService.GetPaymentStatus(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	userID, _ := c.Get("userID")
	if transaction.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to view this transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *PaymentHandler) GetHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	history, err := h.paymentService.GetPaymentHistory(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payment history", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *PaymentHandler) HandleNotification(c *gin.Context) {
	var notificationPayload map[string]interface{}
	if err := c.ShouldBindJSON(&notificationPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification payload"})
		return
	}

	err := h.paymentService.HandleNotification(notificationPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle notification", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification handled successfully"})
}

func (h *PaymentHandler) CreateQrisPayment(c *gin.Context) {
	var req models.CreateQrisPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Authenticated user not found"})
		return
	}
	resp, err := h.paymentService.CreateQrisPayment(&req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create QRIS payment", "details": fmt.Sprintf("%v", err)})
		return
	}

	c.JSON(http.StatusOK, resp)
}
