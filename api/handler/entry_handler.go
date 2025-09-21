package handler

import (
	"github.com/gin-gonic/gin"
)

type EntryHandler struct{}

func NewEntryHandler() *EntryHandler {
	return &EntryHandler{}
}

func (h *EntryHandler) GetEntry(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Payment API Gateway with Midtrans Golang by Bagus Subagja",
	})
}