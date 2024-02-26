
package models

// import "time"

type Food struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Thumb       string    `json:"thumb"`
	CategoryID  int       `json:"category_id"`
	RestaurantID int      `json:"restaurant_id"`
	// CreatedAt   string `json:"created_at"`
	// UpdatedAt   string `json:"updated_at"`
}
