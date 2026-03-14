---
name: wen-beego
description: 指导在 wen-beego 框架（基于 beego 的二次开发框架）中新增接口、使用全局变量和助手函数、遵循目录结构规范。当用户提及“在 wen-beego 中开发”、“新增接口”或询问框架用法时，应加载此 Skill。
license: MIT
compatibility: OpenCode >= 1.0
metadata:
  framework: wen-beego
  version: 1.0
---

# wen-beego 二次开发指南

## 项目简介
wen-beego 框架是基于 beego 的 web 开发框架，包含平台端（`apps/admin_plat`）、商户端（`apps/admin_mchnt`）、对外端（`apps/api`），适用于二次开发。

## 快速开始
### 运行服务
- **HTTP 服务**：`go run cmd/http/main.go`
- **MQ 服务**：`go run cmd/mq/main.go`

### 测试登录接口
开发环境验证码固定为 `1234`，`authCodeId` 可任意填写。
```http
POST http://localhost:8013/api/v1/auth/login HTTP/1.1
Content-Type: application/json

{
    "phone": "15912345678",
    "password": "G72+shD3^6",
    "authCode": "1234",
    "authCodeId": "任意字符串",
    "AuthCodeType": "captcha"
}
```
**成功响应示例**：
```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "userInfo": {
            "phone": "159****5678",
            "name": "wzz",
            "username": "wzz",
            "expires": 1773303810000,
            "accessToken": "xxx",
            "refreshToken": "xxx"
        }
    }
}
```

## 项目结构
### 核心目录
- `apps/admin_plat`：平台端（MVC）
- `apps/admin_mchnt`：商户端（MVC）
- `apps/api`：API 模块
- `apps/cron_task`：定时任务
- `apps/common`：公共模块（模型、助手、全局变量等）
- `apps/mq_task`：MQ 任务消费
- `apps/swagger`：Swagger 文档
- `routers/`：路由定义
  - `admin_plat_router.go`：平台端路由
  - `admin_mchnt_router.go`：商户端路由
  - `mq_task_router.go`：MQ 任务路由
  - `crontab_task/router.go`：定时任务路由

> 📁 **完整目录说明**：请读取项目根目录下的 `doc/目录结构.md` 获取详细目录用途。

### 数据库表结构
数据库表定义保存在 `doc/数据库表结构.sql` 中。当需要查询表字段、索引或关联关系时，请读取该文件。

## 新增接口完整流程（以平台端“获取组织单位”为例）

### 步骤 1：定义路由
在对应路由文件（如 `routers/admin_plat_router.go`）中添加路由切片函数。
- **函数命名格式**：`平台类型+模块类型+slices`，例如 `platSystemSlices()`
- **路由元素格式**：`beego.NSCtrl<HTTP方法>("/path", (*控制器包.控制器).方法)`
- **示例**：
```go
func platSystemSlices() []beego.LinkNamespace {
    return []beego.LinkNamespace{
        beego.NSCtrlGet("/system-unit/get", (*adminSystem.UnitController).Get),
    }
}
```
最后将返回的切片合并到总路由中。

### 步骤 2：创建 Controller
文件路径：`apps/admin_plat/controllers/system/unit.go`
- **结构体定义**：嵌入 `commonControllers.AdminBaseController`
- **Swagger 注解**：包含 `@Summary`, `@Description`, `@Tags`, `@Param`, `@Success`, `@Router`
- **方法实现**：解析参数，调用 Service，返回 JSON
```go
package system

import (
    "WenBeego/apps/admin_plat/services/system"
    commonControllers "WenBeego/apps/common/controller"
    "WenBeego/apps/common/dto/page_dto"
    "WenBeego/apps/common/helper"
)

type UnitController struct {
    commonControllers.AdminBaseController
    UnitService system.UnitService
}

// @Summary 获取内部组织管理
// @Description 获取内部组织管理
// @Tags 系统管理-内部组织管理
// @Accept application/json
// @Produce application/json
// @Param parentUnitId query string true "父级ID"
// @Success 200 {object} dto.RespDataListDto
// @Router /admin_plat/system-unit/get [get]
func (c *UnitController) Get() {
    baseParamDto, _ := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
    reqDataListDto, _ := helper.GetReqDataListDto(&c.Controller)
    unitDto := page_dto.SystemUnitListReqDto{
        BaseParamDto:   baseParamDto,
        ReqDataListDto: reqDataListDto,
        Name:           c.GetString("name"),
        Code:           c.GetString("code"),
    }
    unitDto.Status, _ = c.GetInt("status")

    data, err := c.UnitService.GetUnitList(unitDto)
    if err != nil {
        c.Data["json"] = helper.Response(500, err.Error(), nil)
    } else {
        c.Data["json"] = helper.Response(200, "success", data)
    }
    c.ServeJSON()
}
```

### 步骤 3：创建 Service
文件路径：`apps/admin_plat/services/system/unit.go`
- 调用 Common 层的 Model 方法，组合业务逻辑
```go
package system

import (
    "WenBeego/apps/common/dto"
    "WenBeego/apps/common/dto/page_dto"
    CommonSystem "WenBeego/apps/common/services/system"
)

type UnitService struct {
    commonSystemUnit CommonSystem.Unit
}

func (s *UnitService) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (dto.RespDataListDto, error) {
    return s.commonSystemUnit.GetUnitList(unitDto)
}
```

### 步骤 4：实现 Model 层（数据访问）
如果逻辑可通用，放在 `apps/common/services/system/unit.go`；否则放在对应客户端目录。
```go
package system

import (
    "WenBeego/apps/common/dto"
    "WenBeego/apps/common/dto/page_dto"
    "WenBeego/apps/common/helper"
    "WenBeego/apps/common/models/base_model"
)

type Unit struct{}

func (s *Unit) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (dto.RespDataListDto, error) {
    var data []base_model.Unit
    var count int64
    // ... 数据库查询逻辑
    return helper.GetRespDataListDto(unitDto.PageSize, unitDto.CurrentPage, count, data)
}
```

### 步骤 5：编写单元测试
测试文件位于 `tests/xxx_test.go`，调用 Service 层方法验证。

> ✅ **检查点**：每个步骤完成后，确认文件路径、包名、导入路径正确。

## 框架全局变量与工具
全局变量集中在 `apps/common/global`，可按需导入。

### 常量与路径
```go
import "WenBeego/apps/common/global/constant"
import "WenBeego/apps/common/global"

fmt.Println(constant.ADMIN)                // 常量
fmt.Println(global.RootPath, global.AppDir) // 项目路径
```

### 配置与缓存
```go
// 读取自定义配置
key, err := global.GetConfigDiy("someKey")

// Redis 操作
data, err := global.RedisCache.Get(ctx, key).Result()
```

### 日志
```go
log := global.Log
log.Info("info message")
log.Error("error message")
```

### 数据库
```go
// 读库
global.GetReadDb().Model(&model).Where(...).Find(&results)
// 写库
global.GetWriteDb().Create(&record)
```

### MQ 任务
```go
args := []tasks.Arg{{Name: "action", Type: "string", Value: dataStr}}
result, err := (&MqClient{}).SendTask("TaskName", args)
```

## 常用助手函数（分类摘要）
助手函数位于 `apps/common/helper`，以下为常用分类及示例。

### 数组操作 (`helper/array.go`)
- `InArray(val string, array interface{})`：判断字符串是否在数组中。
- `ArrayMerge(arr1, arr2 []string)`：合并字符串数组。

### 类型转换 (`helper/conver.go`)
- `Interface2MapInterface(i interface{})`：任意类型转 `map[string]interface{}`。
- `String2Int64(s string)`：字符串转 int64。

### 时间处理 (`helper/datetime.go`)
- `GetTimestamp()`：获取当前时间戳。
- `GetTimeString(format ...string)`：获取格式化的时间字符串。

### 文件路径 (`helper/filepath.go`)
- `PathIsExist(path string)`：判断路径是否存在。
- `MkdirAll(path string)`：递归创建目录。

### 密码安全 (`helper/password.go`)
- `GenerateCryptPassword(password string)`：生成加密密码。
- `CompareCryptPassword(password, cryptPassword string)`：验证密码。

### Redis 辅助 (`helper/redis.go`)
- `RedisPut(key string, value interface{}, timeout int)`：写入 Redis。
- `RedisGet(key string)`：读取 Redis。

### 请求处理 (`helper/request.go`)
- `GetReqToken(ctx)`：从请求中提取 Token。
- `GetBaseParamDto(ctx, moduleName)`：获取基础请求参数。

### 响应处理 (`helper/response.go`)
- `Response(code int, message string, data interface{})`：标准化 JSON 响应。
- `GetRespDataListDto(pageSize, currentPage int, total int64, list interface{})`：构造分页响应。

### 通用工具 (`helper/utils.go`)
- `Ternary(condition bool, trueVal, falseVal T)`：泛型三元运算符。
- `IsCellPhone(phone string)`：验证手机号。
- `Md5(str string)`：计算 MD5。

### 树形结构转换 (`helper/tree/`)
- `BuildTree(nodes []TreeNode)`：构建树形结构。
- `ConvertSliceToTree(slice interface{}, idField, parentIdField, childrenField string)`：反射版通用树构建。

> 📘 **完整函数列表**：请直接阅读源码目录 `apps/common/helper/` 中的对应文件，需要具体实现时可打开 `.go` 文件查看。

## ⚠️ 安全与注意事项
- **数据库写操作**（如 `global.GetWriteDb().Create`）应确保在事务中或数据已验证。
- **文件操作**（删除、覆盖）前必须先请求用户确认。
- **生产环境**避免打印敏感信息（密码、Token）。
- **权限校验**：涉及组织单位或角色的接口，调用 `helper.CheckUserHasUnit` 等函数进行鉴权。

## 外部参考文件指引
- 完整目录结构：`doc/目录结构.md`（按需读取）
- 数据库表定义：`doc/数据库表结构.sql`（按需读取特定表）
- 助手函数源码：`apps/common/helper/`（按需查看）
