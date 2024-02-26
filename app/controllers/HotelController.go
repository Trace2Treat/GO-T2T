package controllers

import "github.com/gin-gonic/gin"
import (
	"goHotel/app/models"
	"goHotel/app/config"
)




func GetHotel(ctx *gin.Context) {	
	var hotel models.Hotel
	var data []models.Hotel
	db := config.Connect()
	defer db.Close()

	result, err := db.Query("SELECT * from hotels")
	
	if err != nil {
		panic(err.Error())
		}  
		defer result.Close()  
		for result.Next() {
			err := result.Scan(&hotel.Id,&hotel.Name,&hotel.Address,&hotel.Thumb,&hotel.CreatedAt,&hotel.UpdatedAt)
			if err != nil {
				panic(err.Error())
			} else {
				hotel.Thumb = hotel.Thumb
				data = append(data, hotel)
			}
	}  
	ctx.JSON(200, gin.H{
		"statusCode": 200,
		"message": "Success get Data",
		"data": data,
	})
}


// func storeHotel(ctx *gin.Context) {	
// 	db := config.Connect()
// 	defer db.Close()
// 	h := md5.New()
// 	nip := ctx.PostForm("nip")
// 	nama := ctx.PostForm("nama")
// 	role := ctx.PostForm("role")
// 	password := h.Sum([]byte(ctx.PostForm("password")))
// 	insForm, err := db.Prepare("INSERT INTO pegawai(nip, nama, password, role) VALUES(?,?,?,?)")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	insForm.Exec(nip, nama, password, role)
	

// 	ctx.JSON(200, gin.H{
// 		"statusCode": 200,
// 		"message": "Success Post Data",
// 	})
// }
