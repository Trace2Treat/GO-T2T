package models

type TrashRequest struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TrashType     string  `json:"trash_type"`
	ProofPayment  *string  `json:"proof_payment"`
	Point         *float64 `json:"point"`
	TrashWeight   float64 `json:"trash_weight"`
	Latitude      string  `json:"latitude"`
	Longitude     string  `json:"longitude"`
	Thumb         *string  `json:"thumb"`
	UserID        uint    `json:"user_id"`
	DriverID        *uint    `json:"user_id"`
	Status        string  `json:"status"`
}

func (TrashRequest) TableName() string {
	return "trash_requests"
}

// func GetTrashRequestByID(id uint, db *sql.DB) (*TrashRequest, error) {
//     // Query SQL untuk mengambil data permintaan sampah berdasarkan ID
//     query := "SELECT id, status, trash_weight, point, proof_payment, user_id, driver_id FROM trash_requests WHERE id = ?"
//     row := db.QueryRow(query, id)

//     // Variabel untuk menyimpan data hasil query
//     var trashRequest TrashRequest

//     // Mendeklarasikan variabel untuk menyimpan nilai-nilai yang akan di-scan dari row
//     var (
//         status       string
//         trashWeight  float64
//         point        float64
//         proofPayment string
//         userID       uint
//         driverID     uint
//     )

//     // Men-scan nilai-nilai dari row ke variabel yang sesuai
//     err := row.Scan(&trashRequest.ID, &status, &trashWeight, &point, &proofPayment, &userID, &driverID)
//     if err != nil {
//         return nil, err
//     }

//     // Mengisi data ke dalam struct TrashRequest
//     trashRequest.Status = status
//     trashRequest.TrashWeight = trashWeight
//     trashRequest.Point = point
//     trashRequest.ProofPayment = proofPayment
//     trashRequest.UserID = userID
//     trashRequest.DriverID = driverID

//     return &trashRequest, nil
// }

// // UpdateUserBalance memperbarui saldo pengguna berdasarkan ID pengguna dan jumlah koin yang diberikan
// func UpdateUserBalance(userID uint, coin float64, db *sql.DB) error {
//     // Query SQL untuk memperbarui saldo pengguna
//     query := "UPDATE users SET balance_coin = balance_coin + ? WHERE id = ?"
//     result, err := db.Exec(query, coin, userID)
//     if err != nil {
//         return err
//     }

//     // Mengecek apakah query berhasil dijalankan
//     rowsAffected, err := result.RowsAffected()
//     if err != nil {
//         return err
//     }

//     if rowsAffected == 0 {
//         return fmt.Errorf("User with ID %d not found", userID)
//     }

//     return nil
// }
