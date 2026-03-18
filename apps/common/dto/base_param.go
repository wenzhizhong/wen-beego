package dto

// 基础参数
type BaseParamDto struct {
	Host       string
	ModuleName string
	UnitId     string
	UserId     string
	UnitUserId string
	IsOfficial bool
}

//  响应参数
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 请求页面数据列表参数
type ReqDataListDto struct {
	PageSize    int `json:"pageSize"`
	CurrentPage int `json:"currentPage"`
	Offset      int
}

// 响应页面数据列表参数
type RespDataListDto struct {
	List        interface{} `json:"list"`
	Total       int64       `json:"total"`
	PageSize    int         `json:"pageSize"`
	CurrentPage int         `json:"currentPage"`
}

type BodyEncryptDto struct {
	EncryptedData string `json:"encryptedData"`
}
