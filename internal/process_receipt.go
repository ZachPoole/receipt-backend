package internal

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/ZachPoole/receipt-backend/domain"
)

func ProcessReceipt(receipt *domain.ReceiptDBEntry) (int, error) {
	points := 0

	// One point for every alphanumeric character in the retailer name.
	points += countAlphanumeric(receipt.Retailer)
	fmt.Println(points, "points: 1 pt for each alphanumeric character")

	// 50 points if the total is a round dollar amount with no cents.
	if math.Mod(receipt.Total, 1) == 0 {
		points += 50
		fmt.Println(50, "points: round dollar amounts for total")
	}

	// 25 points if the total is a multiple of 0.25.
	if math.Mod(receipt.Total, 0.25) == 0 {
		points += 25
		fmt.Println(25, "points: total is multiple of 0.25")
	}

	// 5 points for every two items on the receipt.
	itemPairs := len(receipt.Items) / 2
	points += itemPairs * 5

	fmt.Println(itemPairs*5, "points: 5 points for every item pair")

	// If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer. The
	// result is the number of points earned.
	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))

		if trimmedLength%3 == 0 {
			itemPoints := math.Ceil(item.Price * 0.2)
			fmt.Println(itemPoints, "points: one of the item descriptions is multiple of 3")
			points += int(itemPoints)
		}
	}

	// 6 points if the day in the purchase date is odd.
	if receipt.PurchaseDate.Day()%2 != 0 {
		points += 6
		fmt.Println(6, "points: purchase day is odd")
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	startTime, startErr := time.Parse("15:04", "14:00")
	beforeTime, beforeErr := time.Parse("15:04", "16:00")

	if startErr != nil {
		return 0, startErr
	} else if beforeErr != nil {
		return 0, beforeErr
	}

	if receipt.PurchaseTime.After(startTime) &&
		receipt.PurchaseTime.Before(beforeTime) {
		points += 10
		fmt.Println(10, "points: purchase time between range")
	}
	fmt.Println(points, "TOTAL")
	return points, nil
}

func countAlphanumeric(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count++
		}
	}
	return count
}
