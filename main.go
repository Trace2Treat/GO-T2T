package main

import "github.com/gin-gonic/gin"
import (
	// _"encoding/json"
	// "log"
	"net/http"
	"github.com/dgrijalva/jwt-go"

	// "crypto/md5"
	"fmt"
	"goHotel/app/controllers"
)


func main() {
	fmt.Println("Start Server API ....")
	router := gin.Default()

	api := router.Group("/api")

	api.POST("/login", controllers.LoginHandler)
	api.GET("/user", controllers.MeHandler)

	api.GET("/getFile/:folder/:filename", func(ctx *gin.Context) {
		folder := ctx.Param("folder")
		filename := ctx.Param("filename")
		filePath := "storage/" + folder + "/" + filename 
		ctx.File(filePath)
	})

	auth := api.Group("/")
	auth.Use(AuthMiddleware())
	{

		// auth.GET("/users", controllers.IndexUser)
		// auth.POST("/users/store", controllers.StoreUser)
		// auth.PUT("/users/update/:id", controllers.UpdateUser)
		// auth.DELETE("/users/delete/:id", controllers.DeleteUser)

		auth.GET("/foods", controllers.GetFood)
		auth.POST("/foods/store", controllers.StoreFood)
		auth.POST("/foods/update/:id", controllers.UpdateFood)
		auth.DELETE("/foods/delete/:id", controllers.DeleteFood)

		// auth.GET("/products", controllers.IndexProduct)
		// auth.POST("/products/store", controllers.StoreProduct)
		// auth.PUT("/products/update/:id", controllers.UpdateProduct)
		// auth.DELETE("/products/delete/:id", controllers.DeleteProduct)

		auth.GET("/trash-requests", controllers.GetTrashRequest)
		auth.POST("/trash-requests/store", controllers.StoreTrashRequest)
		auth.POST("/trash-requests/change-status/:id", controllers.ChangeStatus)
		// auth.PUT("/trash-requests/update/:id", controllers.UpdateTrashRequest)
		// auth.POST("/trash-requests/change-status/:id", controllers.ChangeTrashRequestStatus)
		// auth.DELETE("/trash-requests/delete/:id", controllers.DeleteTrashRequest)

		// auth.GET("/restaurant", controllers.IndexRestaurant)
		// auth.POST("/restaurant/store", controllers.StoreRestaurant)
		// auth.PUT("/restaurant/update/:id", controllers.UpdateRestaurant)
		// auth.DELETE("/restaurant/delete/:id", controllers.DeleteRestaurant)

		// auth.GET("/transaction", controllers.IndexTransaction)
		// auth.POST("/transaction/purchase-food", controllers.PurchaseFood)
		// auth.POST("/transaction/change-status/:id", controllers.ChangeTransactionStatus)
		// auth.POST("/transaction/pay/:transaction_code", controllers.PayTransaction)
		// auth.GET("/transaction/generate-qr-code/:transaction_code", controllers.GenerateQRCode)
		// auth.POST("/transaction/store", controllers.StoreTransaction)
		// auth.PUT("/transaction/update/:id", controllers.UpdateTransaction)
		// auth.DELETE("/transaction/delete/:id", controllers.DeleteTransaction)
	}

	router.Run("127.0.0.1:8000")
}

func AuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        bearerToken := ctx.GetHeader("Authorization")
        if bearerToken == "" {
            ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bearer token not provided"})
            return
        }
        tokenString := bearerToken[len("Bearer "):]
        claims := jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte("your-secret-key"), nil // Ganti dengan secret key Anda
        })
        if err != nil {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }
        if !token.Valid {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
            return
        }
        // Extract user ID from claims
        userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
		    ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from token"})
		    return
		}
		userID := int(userIDFloat)
		// email := claims["email"].(string)


        if !ok {
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from token"})
            return
        }
		ctx.Set("user_id", uint(userID))
        ctx.Next()
    }
}