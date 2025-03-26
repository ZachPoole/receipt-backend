package main

import (
	"github.com/ZachPoole/receipt-backend/api"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	api.Setup(router)

	router.Run("localhost:8080")
}
