package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReceiptDBEntry struct {
	Id uuid.UUID
	Receipt
}

type Receipt struct {
	Retailer     string
	PurchaseDate time.Time
	PurchaseTime time.Time
	Items        []ReceiptItem
	Total        float64
}

type ReceiptItem struct {
	ShortDescription string
	Price            float64
}
