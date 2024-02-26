package controllers

import "github.com/gin-gonic/gin"
import (
	"goHotel/app/models"
	"goHotel/app/config"
	"time"
	"fmt"
	"net/http"
	"strconv"
)

func GetTrashRequest(ctx *gin.Context) {
	var trash models.TrashRequest
	var data []models.TrashRequest
	db := config.Connect()
	defer db.Close()

	limit := ctx.DefaultQuery("limit", "10")
	search := ctx.Query("search")
	id := ctx.Query("id")

	// trash request
	query := "SELECT id, trash_type,proof_payment,point,trash_weight,latitude,longitude,thumb,user_id,status FROM trash_requests"

	if id != "" {
		query += " WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&trash.ID, &trash.TrashType, &trash.ProofPayment, &trash.Point, &trash.TrashWeight, &trash.Latitude, &trash.Longitude, &trash.Thumb, &trash.UserID, &trash.Status)
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
		data = append(data, trash)
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
			err := rows.Scan(&trash.ID, &trash.TrashType, &trash.ProofPayment, &trash.Point, &trash.TrashWeight, &trash.Latitude, &trash.Longitude, &trash.Thumb, &trash.UserID, &trash.Status)
			if err != nil {
			fmt.Println("Error querying the database:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status":     "error",
					"statusCode": http.StatusInternalServerError,
					"message":    "Gagal mengambil data",
					"data":       nil,
					"timestamp":  time.Now().Format(time.RFC3339),
				})
				return
			}
			data = append(data, trash)
		}
	}

	var totalCount int
	db.QueryRow("SELECT COUNT(id) FROM trash_requests").Scan(&totalCount)
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

func StoreTrashRequest(ctx *gin.Context) {
	var trash models.TrashRequest
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

	trash_type := ctx.PostForm("trash_type")
	trash_weight := ctx.PostForm("trash_weight")
	latitude := ctx.PostForm("latitude")
	longitude := ctx.PostForm("longitude")
	userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
        return
    }

	_, err = db.Exec("INSERT INTO trash_requests (trash_type, trash_weight, latitude, longitude,user_id) VALUES (?, ?, ?, ?, ?)",
		trash_type, trash_weight, latitude, longitude,userID)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to store trash data",
			"data":       nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Success create data",
		"data":       trash,
	})
}

func UpdateTrashRequest(ctx *gin.Context) {
    var trash models.TrashRequest
	db := config.Connect()
	defer db.Close()

    // err := ctx.Request.ParseMultipartForm(10 << 20) // Maksimum ukuran berkas: 10 MB
    // if err != nil {
    //     ctx.JSON(http.StatusBadRequest, gin.H{
    //         "statusCode": http.StatusBadRequest,
    //         "message":    "Failed to process input data",
    //         "data":       nil,
    //     })
    //     return
    // }
    
    id := ctx.Param("id")
	trash_type := ctx.PostForm("trash_type")
	trashWeightStr := ctx.PostForm("trash_weight")
	
	trash_weight, err := strconv.ParseFloat(trashWeightStr, 64)
    if err != nil {
        ctx.JSON(400, gin.H{
            "error": "Harga tidak valid",
        })
        return
    }


	latitude := ctx.PostForm("latitude")
	longitude := ctx.PostForm("longitude")


    trash.TrashType = trash_type
    trash.TrashWeight = trash_weight
    trash.Latitude = latitude
    trash.Longitude = longitude

    query := "UPDATE trash_requests SET trash_type=?, trash_weight=?, latitude=?, longitude=?"
    params := []interface{}{trash.TrashType, trash.TrashWeight, trash.Latitude, trash.Longitude}
    
    _, err = db.Exec(query, append(params, id)...)
    if err != nil {
        fmt.Println("Error querying the database:", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "statusCode": http.StatusInternalServerError,
            "message":    "Failed to update trash data",
            "data":       nil,
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "statusCode": http.StatusOK,
        "message":    "Success update data",
        "data":       trash,
    })
}

func DeleteTrashRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	db := config.Connect()
	defer db.Close()

	_, err := db.Exec("DELETE FROM trashs WHERE id=?", id)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Failed to delete trash data",
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

func ChangeStatus(ctx *gin.Context) {

	status := ctx.PostForm("status")
	proof_payment := ctx.PostForm("proof_payment")

	fmt.Println(status)
	fmt.Println(proof_payment)
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
	    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
	    return
	}
	
	userID, ok := userIDInterface.(uint)
	if !ok {
	    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert user ID to uint"})
	    return
	}



	var trash models.TrashRequest
	var data []models.TrashRequest
	db := config.Connect()
	defer db.Close()

	id := ctx.Param("id")

	query := "SELECT id, trash_type,proof_payment,point,trash_weight,latitude,longitude,thumb,user_id,status FROM trash_requests"

		query += " WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&trash.ID, &trash.TrashType, &trash.ProofPayment, &trash.Point, &trash.TrashWeight, &trash.Latitude, &trash.Longitude, &trash.Thumb, &trash.UserID, &trash.Status)
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
		data = append(data, trash)
	

		fmt.Println(data)
		fmt.Println(trash)



	usPointer := &userID
	trash.Status = status
	trash.DriverID = usPointer

	switch status {
	case "Delivered":
		_, err = db.Exec("UPDATE trash_requests SET status = ?, driver_id = ? WHERE id = ?", trash.Status,trash.DriverID, trash.ID)
		if err != nil {
		    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trash request"})
		    return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"statusCode": 200,
			"message":    "Data berhasil diupdate",
			"data":       nil,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		return

	case "Approved":
		// Rp 10.000 / Kg
		coin := 10000 * trash.TrashWeight
		point := &coin
		trash.Point = point
		_, err = db.Exec("UPDATE trash_requests SET status = ?, driver_id = ? WHERE id = ?", trash.Status,trash.DriverID, trash.ID)
		if err != nil {
		    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trash request"})
		    return
		}

		_, err = db.Exec("UPDATE users SET balance_coin = balance_coin + ? WHERE id = ?", trash.Point, trash.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user balance"})						    
			return	
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"statusCode": 200,
			"message":    "Data berhasil diupdate",
			"data":       nil,
			"timestamp":  time.Now().Format(time.RFC3339),
		})
		return


	case "Finished":
		pp := &proof_payment
		trash.ProofPayment = pp
	case "Received":

	default:
		// Lakukan tindakan sesuai kebutuhan
	}

	// Simpan perubahan status dan driver ID
	_, err = db.Exec("UPDATE trash_requests SET status = ?, proof_payment = ? WHERE id = ?", trash.Status, trash.ProofPayment, trash.ID)
	if err != nil {
	    ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trash request"})
	    return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"statusCode": 200,
		"message":    "Data berhasil diupdate",
		"data":       nil,
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}