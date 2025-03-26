package api

import (
	"github.com/ZachPoole/receipt-backend/internal"
	"github.com/gin-gonic/gin"
)

func Setup(gin *gin.Engine) {
	router := gin.Group("")

	router.POST("/receipts/process", internal.HandleProcessReceipt)
	router.GET("/receipts/:id/points", internal.HandleGetReceiptPointsById)
}
