package controllers

import "github.com/gin-gonic/gin"
import (
	"goHotel/app/models"
	"goHotel/app/config"
	// "crypto/aes"
	// "crypto/cipher"
	// "crypto/rand"
	// "encoding/base64"
	// "encoding/json"
	// "errors"
	"time"
	"fmt"
	// "net/http"
	// "crypto/rand"
	// "encoding/hex"
	// "encoding/json"
	// "log"
	"net/http"
	// "strings"

	cache "github.com/patrickmn/go-cache"
	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

type AppContext struct {
	DB *cache.Cache
}

//Auth - simple auth token
type Auth struct {
	Token string `json:"token"`
}


func LoginHandler(ctx *gin.Context) {
	var loginData struct {
        ID                int `json:"id"`
        Email                string `json:"email"`
        Role                string `json:"role"`
        Password             string `json:"password"`
    }

    // Mengikat data JSON yang diterima ke struct
    if err := ctx.BindJSON(&loginData); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
        return
    }

    // Menggunakan data yang sudah diikat
    email := loginData.Email
    password := loginData.Password

	// Koneksi ke database
	db := config.Connect()
	defer db.Close()

	// Query untuk mencari pengguna berdasarkan email
	row := db.QueryRow("SELECT id, name, password, email, phone, address, avatar, role, status, balance_coin FROM users WHERE email = ?", email)

	var user models.User
	// Mengambil hasil query dan memasukkannya ke dalam struct User
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Phone, &user.Address, &user.Avatar, &user.Role, &user.Status, &user.BalanceCoin)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
	    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt stored password"})
	    return
	}


	fmt.Println(user.ID)
	token := jwt.New(jwt.SigningMethodHS256)
    // Menambahkan klaim (claim) ke token
    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = user.ID
    claims["email"] = user.Email
    claims["role"] = user.Role
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token akan kedaluwarsa dalam 24 jam

    // Menandatangani token dengan secret key
    secretKey := []byte("your-secret-key") // Ganti ini dengan secret key Anda sendiri
    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Mengirimkan token sebagai respons
    ctx.JSON(http.StatusOK, gin.H{"access_token": tokenString})

}

func MeHandler(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bearer token not provided"})
		return
	}
	tokenString := bearerToken[len("Bearer "):]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil // Ganti dengan secret key Anda
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		return
	}
	email := claims["email"].(string)

	var user models.User
	db := config.Connect()
	defer db.Close()
	row := db.QueryRow("SELECT id, name, password, email, phone, address, avatar, role, status, balance_coin FROM users WHERE email = ?", email)
	err = row.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Phone, &user.Address, &user.Avatar, &user.Role, &user.Status, &user.BalanceCoin)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
	}
	ctx.JSON(http.StatusOK, user)
}
