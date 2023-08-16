package request

import "ginDatabaseMhs/model/entity"

type MahasiswaResponse struct {
	Status  interface{}      `json:"status"`
	Message interface{}      `json:"message"`
	Data    entity.User_data `json:"data"`
}

type ResponseToGetAll struct {
	Message string             `json:"message"`
	UserId  int64              `json:"user_id"`
	Data    int                `json:"data"`
	MHS     []entity.User_data `json:"todos"`
}

type IDResponse struct {
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type DeleteResponse struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
}

type UpdateResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"data"`
	MHS     interface{} `json:"todos"`
}
