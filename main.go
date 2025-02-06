package main

import (
	"net/http"
	"receipt-processor/models"
	"receipt-processor/rules"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// In-memory storage
var receipts = make(map[string]models.Receipt)
var rulesEngine = rules.NewRulesEngine()

func main() {
	router := gin.Default()

	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getPoints)

	router.Run(":8080")
}

// POST /receipts/process
func processReceipt(c *gin.Context) {
	var receipt models.Receipt

	// Bind JSON and validate schema
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid."})
		return
	}

	// Validate custom rules
	if err := models.Validate.Struct(receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid. " + err.Error()})
		return
	}

	// Validate each Item in the Items array
	for _, item := range receipt.Items {
		if err := models.Validate.Struct(item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The receipt is invalid. " + err.Error()})
			return
		}
	}

	// Generate unique ID and store receipt
	id := uuid.New().String()
	receipts[id] = receipt

	c.JSON(http.StatusOK, gin.H{"id": id})
}


// GET /receipts/{id}/points
func getPoints(c *gin.Context) {
	id := c.Param("id")

	// Ensure the ID exists in our in-memory store
	receipt, exists := receipts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "No receipt found for that ID."})
		return
	}

	// Calculate points for the receipt
	points := rulesEngine.CalculatePoints(receipt)
	c.JSON(http.StatusOK, gin.H{"points": points})
}
