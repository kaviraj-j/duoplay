package model

type User struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
