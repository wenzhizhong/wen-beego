package global

import (
	"github.com/beego/beego/v2/core/logs"
)

// 注册中间件，设置日志全局变量
var Log *logs.BeeLogger
