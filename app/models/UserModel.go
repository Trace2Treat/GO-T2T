package models


type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	Avatar      string    `json:"avatar"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	BalanceCoin float64   `json:"balance_coin"`
}

func (User) TableName() string {
	return "users"
}