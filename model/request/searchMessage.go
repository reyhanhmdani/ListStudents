package request

import "ginDatabaseMhs/model/entity"

type SearchResponse struct {
	Status int                `json:"status"`
	Data   []entity.User_data `json:"data"`
	Total  int64              `json:"total"`
}
