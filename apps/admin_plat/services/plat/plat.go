package plat

import (
	"errors"
)

type Plat struct {
}

func (c *Plat) GetUserUnit(userId string) (interface{}, error) {
	if userId == "" {
		return nil, errors.New("获取组织信息失败，请先登录！")
	}

	return nil, nil
}
