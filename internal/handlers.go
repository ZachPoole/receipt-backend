package internal

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ZachPoole/receipt-backend/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// in-memory storage of receipts and points
var receipts = make(map[uuid.UUID](domain.ReceiptDBEntry))
var receiptPoints = make(map[uuid.UUID](int))

func HandleProcessReceipt(context *gin.Context) {
	var processReceiptBody domain.ProcessReceiptRequest

	if err := context.ShouldBindJSON(&processReceiptBody); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Error binding JSON payload to ProcessReceiptRequest model"})

		return
	}

	// converting types from string request values to DB types
	newId := uuid.New()
	purchaseDate, err := time.Parse("2006-01-02", processReceiptBody.PurchaseDate)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Error parsing PurchaseDate from Request"})

		return
	}

	purchaseTime, err := time.Parse("15:04", processReceiptBody.PurchaseTime)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Error parsing PurchaseTime from Request"})

		return
	}

	// parsing string inputs for price into floats
	var dbReceiptItems []domain.ReceiptItem

	for _, item := range processReceiptBody.Items {
		floatVal, err := strconv.ParseFloat(item.Price, 64)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   err,
				"message": "Error converting Item Price from string to float"})

			return
		}

		newDbReceiptItem := domain.ReceiptItem{
			ShortDescription: item.ShortDescription,
			Price:            floatVal,
		}

		dbReceiptItems = append(dbReceiptItems, newDbReceiptItem)
	}

	// parsing receipt total from string input to float
	receiptTotalFloat, err := strconv.ParseFloat(processReceiptBody.Total, 64)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Error converting Receipt Total from string to float"})

		return
	}

	// creating DB record
	newReceipt := domain.ReceiptDBEntry{
		Id: newId,
		Receipt: domain.Receipt{
			Retailer:     processReceiptBody.Retailer,
			Items:        dbReceiptItems,
			Total:        receiptTotalFloat,
			PurchaseDate: purchaseDate,
			PurchaseTime: purchaseTime,
		},
	}

	// "insert" new record
	receipts[newId] = newReceipt

	// process receipt to get total points
	processedReceiptPoints, err := ProcessReceipt(&newReceipt)

	if err != nil {
		receiptPoints[newId] = 0
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Error processing points for Receipt"})

		return
	}

	// update point total for new receipt
	receiptPoints[newId] = processedReceiptPoints

	newReceiptResponse := domain.ProcessReceiptResponse{
		Id: newId,
	}

	context.IndentedJSON(http.StatusCreated, newReceiptResponse)

}

func HandleGetReceiptPointsById(context *gin.Context) {
	id := context.Param("id")
	parsedId, _ := uuid.Parse(id)

	receiptPointsRecord, exists := receiptPoints[parsedId]

	if !exists {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Receipt Points Record with the provided ID was not found"})

		return
	}

	receiptPointsResponse := domain.ReceiptGetPointsResponse{
		Points: receiptPointsRecord,
	}

	context.IndentedJSON(http.StatusFound, receiptPointsResponse)
}
