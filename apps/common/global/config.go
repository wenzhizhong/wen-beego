package global

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/server/web"
)

/**
* 获取配置
* @param string section 配置节
* @return map[string]string, error
 */
func GetConfig(section string) (map[string]interface{}, error) {
	interfaceConfig, err := web.AppConfig.DIY(section)
	if err != nil {
		return nil, err
	}
	tmpMapConfig, ok := interfaceConfig.(map[string]interface{})
	if !ok {
		fmt.Println("类型转换错误, \ninterfaceConfig=", interfaceConfig)
		return nil, errors.New("类型转换错误")
	}
	return tmpMapConfig, nil
}

func GetConfigDiy(section string) (interface{}, error) {
	interfaceConfig, err := web.AppConfig.DIY(section)
	if err != nil {
		return nil, err
	}
	return interfaceConfig, nil
}
