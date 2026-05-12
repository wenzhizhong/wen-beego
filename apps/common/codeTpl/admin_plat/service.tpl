{{if .IsMultiApp}}
package {{.MenuModule}}

import (
	commonService "WenBeego/apps/common/services/{{.MenuModule}}"
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/{{.MenuModule}}_dto"
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

func (s *{{.ModelName}}Service) GetList(baseParamDto dto.BaseParamDto, pageSize, offset int, keyword string) (*dto.RespDataListDto, error) {
	return s.common{{.ModelName}}.GetList(baseParamDto, pageSize, offset, keyword)
}

{{else}}
package {{.MenuModule}}

import (
	"WenBeego/apps/{{.AppModule}}/models_ar"
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/{{.MenuModule}}_dto"
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
		return s.{{.ModelName}}Ar.Insert(tx, &data.{{.ModelName}})
	})
}

func (s *{{.ModelName}}Service) Edit(baseParamDto dto.BaseParamDto, data {{.MenuModule}}_dto.{{.ModelName}}Dto) error {
	_ = baseParamDto
	return global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
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

func (s *{{.ModelName}}Service) GetList(baseParamDto dto.BaseParamDto, pageSize, offset int, keyword string) (*dto.RespDataListDto, error) {
	_ = baseParamDto
	data, count, err := s.{{.ModelName}}Ar.GetList(pageSize, offset, keyword)
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
