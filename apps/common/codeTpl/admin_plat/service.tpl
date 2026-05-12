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

func (s *{{.ModelName}}Service) Add(data {{.MenuModule}}_dto.{{.ModelName}}Dto) error {
	return global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		id, err := helper.GetUuid()
		if err != nil {
			return err
		}
		data.Id = id
		return s.{{.ModelName}}Ar.Insert(tx, data)
	})
}

func (s *{{.ModelName}}Service) Edit(data {{.MenuModule}}_dto.{{.ModelName}}Dto) error {
	return global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		return s.{{.ModelName}}Ar.Update(tx, data)
	})
}

func (s *{{.ModelName}}Service) Del(id string) error {
	return global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		return s.{{.ModelName}}Ar.Delete(tx, id)
	})
}

func (s *{{.ModelName}}Service) GetDetail(id string) (models.{{.ModelName}}, error) {
	data, err := s.{{.ModelName}}Ar.GetById(id)
	if err != nil {
		return models.{{.ModelName}}{}, err
	}
	return models.{{.ModelName}}{ {{.ModelName}}: data}, nil
}

func (s *{{.ModelName}}Service) GetList(pageSize, offset int, keyword string) (*dto.RespDataListDto, error) {
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
