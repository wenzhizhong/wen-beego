{{if .IsMultiTenant}}
package {{.MenuModule}}

import (
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/{{.MenuModule}}_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type {{.ModelName}}Service struct{}

func (s *{{.ModelName}}Service) Add(baseParamDto dto.BaseParamDto, data {{.MenuModule}}_dto.{{.ModelName}}Dto) error {
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

		switch baseParamDto.ModuleName {
		case "admin_plat":
			return base_ar.Insert{{.ModelName}}(tx, &data.{{.ModelName}}, &models.{{.PlatModelName}}{})
		case "admin_mchnt":
			return base_ar.Insert{{.ModelName}}(tx, &data.{{.ModelName}}, &models.{{.MchntModelName}}{})
		default:
			return errors.New("模块名称错误")
		}
	})
}

func (s *{{.ModelName}}Service) Edit(baseParamDto dto.BaseParamDto, data {{.MenuModule}}_dto.{{.ModelName}}Dto) error {
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

		switch baseParamDto.ModuleName {
		case "admin_plat":
			return base_ar.Update{{.ModelName}}(tx, &data.{{.ModelName}}, &models.{{.PlatModelName}}{})
		case "admin_mchnt":
			return base_ar.Update{{.ModelName}}(tx, &data.{{.ModelName}}, &models.{{.MchntModelName}}{})
		default:
			return errors.New("模块名称错误")
		}
	})
}

func (s *{{.ModelName}}Service) Del(baseParamDto dto.BaseParamDto, id string) error {
	return global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		switch baseParamDto.ModuleName {
		case "admin_plat":
			return base_ar.Delete{{.ModelName}}(tx, id, &models.{{.PlatModelName}}{})
		case "admin_mchnt":
			return base_ar.Delete{{.ModelName}}(tx, id, &models.{{.MchntModelName}}{})
		default:
			return errors.New("模块名称错误")
		}
	})
}

func (s *{{.ModelName}}Service) GetDetail(baseParamDto dto.BaseParamDto, id string) (base_model.{{.ModelName}}, error) {
	var data base_model.{{.ModelName}}
	switch baseParamDto.ModuleName {
	case "admin_plat":
		data, err := base_ar.Get{{.ModelName}}ById(id, &models.{{.PlatModelName}}{})
		if err != nil {
			return data, err
		}
		return data, nil
	case "admin_mchnt":
		data, err := base_ar.Get{{.ModelName}}ById(id, &models.{{.MchntModelName}}{})
		if err != nil {
			return data, err
		}
		return data, nil
	default:
		return data, errors.New("模块名称错误")
	}
}

func (s *{{.ModelName}}Service) GetList(baseParamDto dto.BaseParamDto, pageSize, offset int, searchDto {{.MenuModule}}_dto.{{.ModelName}}Dto) (*dto.RespDataListDto, error) {
	var data []base_model.{{.ModelName}}
	var count int64
	var err error

	switch baseParamDto.ModuleName {
	case "admin_plat":
		data, count, err = base_ar.Get{{.ModelName}}List(pageSize, offset, searchDto, &models.{{.PlatModelName}}{}, &models.PlatUser{}{{if .HasUnitId}}, &models.Plat{}{{end}})
	case "admin_mchnt":
		data, count, err = base_ar.Get{{.ModelName}}List(pageSize, offset, searchDto, &models.{{.MchntModelName}}{}, &models.MchntUser{}{{if .HasUnitId}}, &models.Mchnt{}{{end}})
	default:
		return nil, errors.New("模块名称错误")
	}
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

{{else}}
package {{.MenuModule}}

import (
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/{{.MenuModule}}_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar"
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

		return s.{{.ModelName}}Ar.Insert(tx, &data.{{.ModelName}})
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
		return s.{{.ModelName}}Ar.Update(tx, &data.{{.ModelName}})
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
