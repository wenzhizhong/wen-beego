package dto

// 基础参数
type BaseParamDto struct {
	ModuleName string
	UnitId     string
	UserId     string
}

//  响应参数
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 请求页面数据列表参数
type ReqDataListDto struct {
	PageSize int `json:"pageSize"`
	Page     int `json:"page"`
	Offset   int
}

// 响应页面数据列表参数
type RespDataListDto struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	PageSize int         `json:"pageSize"`
	Page     int         `json:"page"`
}
