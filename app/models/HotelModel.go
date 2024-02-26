package models

type Hotel struct {
	Id        string `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
	Address string `form:"address" json:"address"`
	Thumb  string `form:"thumb" json:"thumb"`
	CreatedAt  *string `form:"created_at" json:"created_at"`
	UpdatedAt  *string `form:"updated_at" json:"updated_at"`
}
