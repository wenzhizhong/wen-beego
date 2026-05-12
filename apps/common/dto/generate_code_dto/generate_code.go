package generate_code_dto

import "WenBeego/apps/common/models/base_model"

type GenerateCodeDto struct {
	base_model.GenerateCode
}

type SaveFormDetailDto struct {
	Id        string `json:"id"`
	TableName string `json:"table_name"`
	Data      string `json:"data"`
}

type GetTableDetailDto struct {
	TableName string `json:"tableName"`
}

type GenCodeRunDto struct {
	TableGenerateCodeId string   `json:"tableGenerateCodeId"`
	MenuName            string   `json:"menuName"`
	AppModule           string   `json:"appModule"`
	MenuModule          string   `json:"menuModule"`
	BizModule           string   `json:"bizModule"`
	CodeType            []string `json:"codeType"`
	ViewType            string   `json:"viewType"`
}

type DownloadCodeDto struct {
	ZipPath string `json:"zipPath"`
}

type DelFormDetailDto struct {
	Ids []string `json:"ids"`
}

type GenCodeParamDto struct {
	ViewTypes         map[string]string `json:"viewTypes"`
	GenerateCodeTypes map[string]string `json:"generateCodeTypes"`
	MenuName          string            `json:"menuName"`
}
