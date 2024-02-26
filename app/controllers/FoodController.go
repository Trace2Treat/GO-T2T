package controllers

import "github.com/gin-gonic/gin"
import (
	"goHotel/app/models"
	"goHotel/app/config"
	"goHotel/app/utils"
	"time"
	"fmt"
	"net/http"
	"mime/multipart"
	"strconv"
	"path/filepath"
	"os"
	"io"
)

func GetFood(ctx *gin.Context) {
	var food models.Food
	var data []models.Food
	db := config.Connect()
	defer db.Close()

	limit := ctx.DefaultQuery("limit", "10")
	search := ctx.Query("search")
	id := ctx.Query("id")

	query := "SELECT id, name, description, price, stock, thumb, category_id, restaurant_id FROM foods"

	if id != "" {
		query += " WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&food.ID, &food.Name, &food.Description, &food.Price, &food.Stock, &food.Thumb, &food.CategoryID, &food.RestaurantID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":     "error",
				"statusCode": http.StatusNotFound,
				"message":    "Data tidak ditemukan",
				"data":       nil,
				"timestamp":  time.Now().Format(time.RFC3339),
			})
			return
		}
		data = append(data, food)
	} else {
		if search != "" {
			query += " WHERE name LIKE ?"
			search = "%" + search + "%"
		}

		query += " LIMIT 10"
		rows, err := db.Query(query)
		if err != nil {
			fmt.Println("Error querying the database:", err)
			fmt.Println("Error querying the database:", limit)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":     "error",
				"statusCode": http.StatusInternalServerError,
				"message":    "Gagal mengambil data",
				"data":       nil,
				"timestamp":  time.Now().Format(time.RFC3339),
			})
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&food.ID, &food.Name, &food.Description, &food.Price, &food.Stock, &food.Thumb, &food.CategoryID, &food.RestaurantID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status":     "error",
					"statusCode": http.StatusInternalServerError,
					"message":    "Gagal mengambil data",
					"data":       nil,
					"timestamp":  time.Now().Format(time.RFC3339),
				})
				return
			}
			data = append(data, food)
		}
	}

	var totalCount int
	db.QueryRow("SELECT COUNT(id) FROM foods").Scan(&totalCount)
	totalPages := (totalCount + (totalCount % 10)) / 10

	ctx.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"statusCode":     http.StatusOK,
		"message":        "Data berhasil diambil",
		"data":           data,
		"timestamp":      time.Now().Format(time.RFC3339),
		"current_page":   1,
		"first_page_url": ctx.Request.URL.String(),
		"from":           1,
		"last_page":      totalPages,
		"last_page_url":  ctx.Request.URL.String(),
		"links":          []interface{}{gin.H{"url": nil, "label": "&laquo; Previous", "active": false}, gin.H{"url": ctx.Request.URL.String(), "label": "1", "active": true}, gin.H{"url": nil, "label": "Next &raquo;", "active": false}},
		"next_page_url":  nil,
		"path":           ctx.Request.URL.Path,
		"per_page":       10,
		"prev_page_url":  nil,
		"to":             1,
		"total":          totalCount,
	})
}

func StoreFood(ctx *gin.Context) {
	var food models.Food
	db := config.Connect()
	defer db.Close()

	err := ctx.Request.ParseMultipartForm(10 << 20) // Maksimum ukuran berkas: 10 MB
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Failed to process input data",
			"data":       nil,
		})
		return
	}

	name := ctx.PostForm("name")
    description := ctx.PostForm("description")
    priceStr := ctx.PostForm("price")
    stockStr := ctx.PostForm("stock")
    categoryIDStr := ctx.PostForm("category_id")
    restaurantIDStr := 3

    price, err := strconv.ParseFloat(priceStr, 64)
    if err != nil {
        ctx.JSON(400, gin.H{
            "error": "Harga tidak valid",
        })
        return
    }

    stock, err := strconv.Atoi(stockStr)
    if err != nil {
        ctx.JSON(400, gin.H{
            "error": "Stok tidak valid",
        })
        return
    }

    categoryID, err := strconv.Atoi(categoryIDStr)
    if err != nil {
        ctx.JSON(400, gin.H{
            "error": "Category ID tidak valid",
        })
        return
    }


	food.Name = name
	food.Description = description
	food.Price = price
	food.Stock = stock
	food.CategoryID = categoryID
	food.RestaurantID = restaurantIDStr

	file, fileHeader, err := ctx.Request.FormFile("thumb")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":     "error",
			"statusCode": http.StatusBadRequest,
			"message":    "File tidak ditemukan",
			"data":       nil,
		})
		return
	}
	defer file.Close()

	filePath, err := utils.SaveFile(fileHeader, "FoodThumb")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":     "error",
			"statusCode": http.StatusInternalServerError,
			"message":    "Gagal menyimpan file",
			"data":       nil,
		})
		return
	}

	food.Thumb = filePath

	_, err = db.Exec("INSERT INTO foods (name, description, price, stock, thumb, category_id, restaurant_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		food.Name, food.Description, food.Price, food.Stock, food.Thumb, food.CategoryID, food.RestaurantID)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to store food data",
			"data":       nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Success create data",
		"data":       food,
	})
}

func saveFile(filename string, file multipart.File, folder string) (string, error) {
	path := filepath.Join("storage", folder, filename)
	newFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, file)
	if err != nil {
		return "", err
	}
	return path, nil
}

func UpdateFood(ctx *gin.Context) {
    var food models.Food
	db := config.Connect()
	defer db.Close()

    err := ctx.Request.ParseMultipartForm(10 << 20) // Maksimum ukuran berkas: 10 MB
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "statusCode": http.StatusBadRequest,
            "message":    "Failed to process input data",
            "data":       nil,
        })
        return
    }
    
    id := ctx.Param("id")
    name := ctx.PostForm("name")
    description := ctx.PostForm("description")
    priceStr := ctx.PostForm("price")
    stockStr := ctx.PostForm("stock")
    categoryIDStr := ctx.PostForm("category_id")
    restaurantIDStr := 3

    price, err := strconv.ParseFloat(priceStr, 64)
    if err != nil {
        // Penanganan kesalahan jika harga tidak valid
        ctx.JSON(400, gin.H{
            "error": "Harga tidak valid",
        })
        return
    }

    stock, err := strconv.Atoi(stockStr)
    if err != nil {
        // Penanganan kesalahan jika stok tidak valid
        ctx.JSON(400, gin.H{
            "error": "Stok tidak valid",
        })
        return
    }

    categoryID, err := strconv.Atoi(categoryIDStr)
    if err != nil {
        // Penanganan kesalahan jika category_id tidak valid
        ctx.JSON(400, gin.H{
            "error": "Category ID tidak valid",
        })
        return
    }


    food.Name = name
    food.Description = description
    food.Price = price
    food.Stock = stock
    food.CategoryID = categoryID
    food.RestaurantID = restaurantIDStr

    file, fileHeader, err := ctx.Request.FormFile("thumb")
	if err == http.ErrMissingFile {
	    food.Thumb = ""
	} else if err != nil {
	    ctx.JSON(http.StatusBadRequest, gin.H{
	        "status":     "error",
	        "statusCode": http.StatusBadRequest,
	        "message":    "Gagal memproses berkas",
	        "data":       nil,
	    })
	    return
	} else {
	    defer file.Close()
	    filePath, err := utils.SaveFile(fileHeader, "FoodThumb")
	    if err != nil {
	        ctx.JSON(http.StatusInternalServerError, gin.H{
	            "status":     "error",
	            "statusCode": http.StatusInternalServerError,
	            "message":    "Gagal menyimpan berkas",
	            "data":       nil,
	        })
	        return
	    }

		food.Thumb = filePath
	}

    query := "UPDATE foods SET name=?, description=?, price=?, stock=?, category_id=?, restaurant_id=?"
    params := []interface{}{food.Name, food.Description, food.Price, food.Stock, food.CategoryID, food.RestaurantID}
    
    if food.Thumb != "" {
        query += ", thumb=?"
        params = append(params, food.Thumb)
    }
    query += " WHERE id=?"

    _, err = db.Exec(query, append(params, id)...)
    if err != nil {
        fmt.Println("Error querying the database:", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "statusCode": http.StatusInternalServerError,
            "message":    "Failed to update food data",
            "data":       nil,
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "statusCode": http.StatusOK,
        "message":    "Success update data",
        "data":       food,
    })
}

func DeleteFood(ctx *gin.Context) {
	id := ctx.Param("id")
	db := config.Connect()
	defer db.Close()

	_, err := db.Exec("DELETE FROM foods WHERE id=?", id)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to delete food data",
			"data":       nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Success delete data",
		"data":       nil,
	})
}