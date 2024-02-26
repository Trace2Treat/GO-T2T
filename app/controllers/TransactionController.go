package controllers

import "github.com/gin-gonic/gin"
import (
	"goHotel/app/models"
	"goHotel/app/config"
	"time"
	"fmt"
	"net/http"
)

func PurchaseFood(ctx *gin.Context) {
	// Generate transaction code
	db := config.Connect()
	defer db.Close()
	orderID := generateTransactionCode()

	// Parse request body
	var requestBody struct {
		Items       []struct {
			FoodID int `json:"food_id"`
			Qty    int `json:"qty"`
		} `json:"items"`
		Address      string  `json:"address"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		RestaurantID int     `json:"restaurant_id"`
	}
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
		return
	}

	// Calculate total price
	total := 0
	for _, item := range requestBody.Items {

		food, err := db.Exec("SELECT price FROM foods WHERE id = ?", item.FoodID)
		if err != nil {
			fmt.Println("Error querying the database:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    "Failed to retrieve food information",
				"data":       nil,
			})
			return
		}
		total += item.Qty * food.Price
	}

	// Check user balance
	userID := ctx.GetInt("user_id")
	userBalance, err := GetUserBalance(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user balance"})
		return
	}
	if total > userBalance {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Insert orders
	for _, item := range requestBody.Items {
		food, err := db.Exec("SELECT price FROM foods WHERE id = ?", item.FoodID)
		if err != nil {
			fmt.Println("Error querying the database:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    "Failed to retrieve food information",
				"data":       nil,
			})
			return
		}

		order := Order{
			TransactionID: orderID,
			FoodID:        item.FoodID,
			Qty:           item.Qty,
			Total:         item.Qty * food.Price,
		}
		if err := InsertOrder(order); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert order"})
			return
		}
	}

	// Create transaction
	transaction := Transaction{
		TransactionCode: orderID,
		Address:         requestBody.Address,
		Latitude:        requestBody.Latitude,
		Longitude:       requestBody.Longitude,
		RestaurantID:    requestBody.RestaurantID,
		Status:          "pending",
		Total:           total,
		UserID:          userID,
		DriverID:        nil,
	}

	// CreateTransaction
	_, err := db.Exec("INSERT INTO transactions (transaction_code, address, latitude, longitude, restaurant_id, status, total, user_id, driver_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", 
	transaction.TransactionCode, transaction.Address, transaction.Latitude, transaction.Longitude, transaction.RestaurantID, transaction.Status, transaction.Total, transaction.UserID, transaction.DriverID)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to store trash data",
			"data":       nil,
		})
		return
	}


	// update user balance
	_, err = db.Exec("UPDATE users SET balance_coin = balance_coin - ? WHERE id = ?", total, userID)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to update user balance",
			"data":       nil,
		})
		return
	}
	

	// Prepare response
	responseData := struct {
		Status       string      `json:"status"`
		StatusCode   int         `json:"statusCode"`
		Message      string      `json:"message"`
		Transaction  Transaction `json:"transaction"`
		Timestamp    time.Time   `json:"timestamp"`
	}{
		Status:     "success",
		StatusCode: http.StatusOK,
		Message:    "Data transaksi berhasil disimpan",
		Transaction: Transaction{
			TransactionCode: orderID,
			Address:         requestBody.Address,
			Latitude:        requestBody.Latitude,
			Longitude:       requestBody.Longitude,
			RestaurantID:    requestBody.RestaurantID,
			Status:          "pending",
			Total:           total,
			UserID:          userID,
			DriverID:        nil,
		},
		Timestamp: time.Now(),
	}

	ctx.JSON(http.StatusOK, responseData)
}

func generateTransactionCode() string {
	return "TRX-" + time.Now().Format("20060102150405") + "-" + strconv.Itoa(AuthUserID) + "-" + strconv.Itoa(rand.Intn(9000)+1000)
}

func ChangeStatusTransaction(ctx *gin.Context) {
	// Get transaction ID from URL parameter
	id := ctx.Param("id")

	// Parse request body
	var requestBody struct {
		Status string `json:"status"`
	}
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
		return
	}

	// Update transaction status in database
	if err := UpdateTransactionStatus(id, requestBody.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction status"})
		return
	}

	// Prepare response
	responseData := struct {
		Status     string    `json:"status"`
		StatusCode int       `json:"statusCode"`
		Message    string    `json:"message"`
		Timestamp  time.Time `json:"timestamp"`
	}{
		Status:     "success",
		StatusCode: http.StatusOK,
		Message:    "Data berhasil diupdate",
		Timestamp:  time.Now(),
	}

	ctx.JSON(http.StatusOK, responseData)
}

func GenerateQRCode(ctx *gin.Context) {
	// Get transaction code from URL parameter
	transactionCode := ctx.Param("transaction_code")

	// Generate QR code text
	qrCodeText := "http://example.com/transaction/pay/" + transactionCode

	// Generate QR code image
	qrCode, err := qrcode.Encode(qrCodeText, qrcode.Medium, 200)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Set content type header
	ctx.Header("Content-Type", "image/png")

	// Send QR code image as response
	ctx.Writer.Write(qrCode)
}

func Pay(ctx *gin.Context) {
	// Get transaction code from URL parameter
	transactionCode := ctx.Param("transaction_code")

	// Query the database to find the transaction
	var transaction Transaction
	if err := db.Where("transaction_code = ?", transactionCode).First(&transaction).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":     "error",
			"statusCode": http.StatusNotFound,
			"message":    "Data tidak ditemukan",
			"data":       nil,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		return
	}

	// Update user balance
	userID := transaction.UserID
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal mengambil data pengguna",
			"data":       nil,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		return
	}
	user.BalanceCoin -= transaction.Total
	if err := db.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal menyimpan data pengguna",
			"data":       nil,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		return
	}

	// Update restaurant balance
	restoID := transaction.RestaurantID
	var resto User
	if err := db.First(&resto, restoID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal mengambil data restoran",
			"data":       nil,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		return
	}
	resto.BalanceCoin += transaction.Total
	if err := db.Save(&resto).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal menyimpan data restoran",
			"data":       nil,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		return
	}

	// Update transaction status
	transaction.Status = "success"
	if err := db.Save(&transaction).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal menyimpan data transaksi",
			"data":       nil,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": http.StatusOK,
		"message":    "Data berhasil diupdate",
		"data":       nil,
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}
