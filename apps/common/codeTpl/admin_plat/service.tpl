{{if .IsMultiApp}}
package {{.MenuModule}}

import (
	commonService "WenBeego/apps/common/services/{{.MenuModule}}"
	"WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/{{.MenuModule}}"
)

type {{.ModelName}}Service struct {
	common{{.ModelName}} commonService.{{.ModelName}}Service
}

func (s *{{.ModelName}}Service) Add(baseParamDto dto.BaseParamDto, data {{.MenuModule}}_dto.{{.ModelName}}Dto) error {
	return s.common{{.ModelName}}.Add(baseParamDto, data)
}

func (s *{{.ModelName}}Service) Edit(baseParamDto dto.BaseParamDto, data {{.MenuModule}}_dto.{{.ModelName}}Dto) error {
	return s.common{{.ModelName}}.Edit(baseParamDto, data)
}

func (s *{{.ModelName}}Service) Del(baseParamDto dto.BaseParamDto, id string) error {
	return s.common{{.ModelName}}.Del(baseParamDto, id)
}

func (s *{{.ModelName}}Service) GetDetail(baseParamDto dto.BaseParamDto, id string) (interface{}, error) {
	return s.common{{.ModelName}}.GetDetail(baseParamDto, id)
}

func (s *{{.ModelName}}Service) GetList(baseParamDto dto.BaseParamDto, pageSize, offset int, searchDto {{.MenuModule}}_dto.{{.ModelName}}Dto) (*dto.RespDataListDto, error) {
	return s.common{{.ModelName}}.GetList(baseParamDto, pageSize, offset, searchDto)
}

{{else}}
package {{.MenuModule}}

import (
	"WenBeego/apps/{{.AppModule}}/models_ar"
	"WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/{{.MenuModule}}"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"fmt"

	"gorm.io/gorm"
)

type {{.ModelName}}Service struct {
	{{.ModelName}}Model models.{{.ModelName}}
	{{.ModelName}}Ar    models_ar.{{.ModelName}}Ar
}

func (s *{{.ModelName}}Service) Add(baseParamDto dto.BaseParamDto, data {{.MenuModule}}_dto.{{.ModelName}}Dto) error {
	_ = baseParamDto
	return global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		id, err := helper.GetUuid()
		if err != nil {
			return err
		}
		data.Id = id
		{{range .Columns}}
			{{- if eq .Name "unit_id" }}		
		data.{{.GoFieldName}} = baseParamDto.UnitId
			{{else if isCreateUserIdFields .Name -}}	
		data.{{.GoFieldName}} = baseParamDto.UnitUserId
			{{else if isCreateTimeFields .Name }}
				{{- if eq .GoType "string" }}
		data.{{.GoFieldName}} = helper.GetTimeString("2006-01-02 15:04:05")
				{{ else -}}
		curTime := helper.GetTime()
		data.{{.GoFieldName}} = {{- if eq .GoType "*time.Time"}}&{{end}}curTime
				{{- end -}}
			{{end -}}
		{{end}}

		return s.{{.ModelName}}Ar.Insert(tx, &data.{{.ModelName}}.{{.ModelName}})
	})
}

func (s *{{.ModelName}}Service) Edit(baseParamDto dto.BaseParamDto, data {{.MenuModule}}_dto.{{.ModelName}}Dto) error {
	_ = baseParamDto
	return global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		{{ range .Columns}} 
			{{- if isUpdateTimeFields .Name }}
				{{- if eq .GoType "string" }}	
		data.{{.GoFieldName}} = helper.GetTimeString("2006-01-02 15:04:05")
				{{ else }}	
		curTime := helper.GetTime()
		data.{{.GoFieldName}} = {{- if eq .GoType "*time.Time"}}&{{end}}curTime
				{{- end -}}
			{{ else if isUpdateUserIdFields .Name -}}	
		data.{{.GoFieldName}} = baseParamDto.UnitUserId 
			{{end -}}
		{{end}}
		return s.{{.ModelName}}Ar.Update(tx, &data.{{.ModelName}}.{{.ModelName}})
	})
}

func (s *{{.ModelName}}Service) Del(baseParamDto dto.BaseParamDto, id string) error {
	_ = baseParamDto
	return global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		return s.{{.ModelName}}Ar.Delete(tx, id)
	})
}

func (s *{{.ModelName}}Service) GetDetail(baseParamDto dto.BaseParamDto, id string) (models.{{.ModelName}}, error) {
	_ = baseParamDto
	data, err := s.{{.ModelName}}Ar.GetById(id)
	if err != nil {
		return models.{{.ModelName}}{}, err
	}
	return models.{{.ModelName}}{ {{.ModelName}}: data}, nil
}

func (s *{{.ModelName}}Service) GetList(baseParamDto dto.BaseParamDto, pageSize, offset int, searchDto {{.MenuModule}}_dto.{{.ModelName}}Dto) (*dto.RespDataListDto, error) {
	_ = baseParamDto
	data, count, err := s.{{.ModelName}}Ar.GetList(pageSize, offset, searchDto)
	if err != nil {
		return nil, fmt.Errorf("获取列表失败: %v", err)
	}
	res := &dto.RespDataListDto{}
	res.List = data
	res.Total = count
	res.PageSize = pageSize
	res.CurrentPage = offset/pageSize + 1
	return res, nil
}
{{end}}
