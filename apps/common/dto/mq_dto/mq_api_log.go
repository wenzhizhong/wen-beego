package mq_dto

type ApiLogDto struct {
	Host   string `json:"host"`
	Uri    string `json:"uri"`
	Method string `json:"method"`
	Ip     string `json:"ip"`
	UnitId string `json:"unit_id"`
	UserId string `json:"user_id"`
}
