package dto

type AdminRequest struct {
	User string `json:"username"`
	Pass string `json:"password"`
}