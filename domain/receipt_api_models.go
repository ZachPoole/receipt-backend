package domain

import (
	"github.com/google/uuid"
)

type ProcessReceiptRequest struct {
	Retailer     string               `json:"retailer"`
	PurchaseDate string               `json:"purchaseDate"`
	PurchaseTime string               `json:"purchaseTime"`
	Items        []InboundReceiptItem `json:"items"`
	Total        string               `json:"total"`
}

type InboundReceiptItem struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ProcessReceiptResponse struct {
	Id uuid.UUID `json:"id"`
}

type ReceiptGetPointsResponse struct {
	Points int `json:"points"`
}
