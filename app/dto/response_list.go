package dto

type ResponseList struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseToken struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
}