package models

type ResposeV1 struct {
	V1 string `json:"w1"`
}

type ResponseUI struct {
	Ui string `json:"ui"`
}

type Item struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
