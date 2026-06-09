---
name: wen-beego-guide
description: 指导在 wen-beego 框架（基于 beego 的二次开发框架）中新增接口、使用全局变量和助手函数、遵循目录结构规范。当用户提及“在 wen-beego 中开发”、“新增接口”或询问框架用法时，应加载此 Skill。
license: MIT
compatibility: OpenCode >= 1.0
metadata:
  framework: wen-beego
  version: 1.0
---

# wen-beego 二次开发指南

## 项目简介
- wen-beego 框架是基于 beego 的 web 开发框架，包含平台端（`apps/admin_plat`）、商户端（`apps/admin_mchnt`）、对外端（`apps/api`）的多用户体系管理系统，适用于二次开发。
- 用户体系，
1.用户都来源于`user`； 
2.基于`user`，新建数据表在区分出：
平台用户（`plat_user`）、
商户用户（`mchnt_user`）、
商户会员用户（`mchnt_customer`）。

## 快速开始
### 运行服务
- **HTTP 服务**：`go run cmd/http/main.go`
- **MQ 服务**：`go run cmd/mq/main.go`
- **MQ 死信队列服务**：`go run cmd/mq/main.go dlx`

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

### 步骤 1：创建数据表
根据需求，创建表及其字段，并添加索引，输出表结构并同步到`doc\数据库表结构.sql`。如果存在则跳过。

### 步骤 2：实现 Model 层
- 根据表结构创建模型层，文件已经存在则跳过。
- 如果不需要区分平台和商户类别的模型, 比如附件就直接创建`apps/common/models/file.go`即可。
但这里‘组织单位’区分商户组织和平台组织，但是组织字段大同小异,相同字段放base_model：`apps/common/models/base_model/unit.go`；
再新增平台组织模型`apps/common/models/plat.go`、商户组织模型`apps/common/models/mchnt.go`, 然后这两个文件里面结构体继承于unit.go，再添加放各自差异字段。

```go
package base_model

type Unit struct {
	Id                string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	// 其他字段....
	Status            int    `json:"status" gorm:"type:int4;not null;default:0;comment:0未审核，1审核通过，2审核不通过，3禁用"`
}

var UNIT_STATUS_UNREVIEWED = 0
var UNIT_STATUS_PASSED = 1
var UNIT_STATUS_UNPASSED = 2
var UNIT_STATUS_DISABLED = 3
var UNIT_STATUS_MAP = map[int]string{
	UNIT_STATUS_UNREVIEWED: "未审核",
	UNIT_STATUS_PASSED:     "审核通过",
	UNIT_STATUS_UNPASSED:   "审核不通过",
	UNIT_STATUS_DISABLED:   "已禁用",
}
```

### 步骤 3：实现 ActiveRecord层（数据访问层）
代码文件已经存在则跳过。
注：如果考虑多用户体系需要共用逻辑的，可以在`apps\common\models_ar\base_ar\xxx.go`创建curd的泛型函数,否则`apps\common\models_ar\xxx.go` 写curd的方法，供`apps/common/services/xxx.go`调用
1.区分多用户体系, `apps\common\models_ar\base_ar\xxx.go`泛型函数:
```go 
package base_ar
// import ...

func GetUnitListByUserId[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf](unitDto page_dto.SystemUnitListReqDto, unitModel UnitModel, unitUserModel UnitUserModel) (listData []base_model.Unit, count int64, err error) {
	// other code ...
	query := global.GetReadDb().
		Model(unitModel).
		Joins(joinStr).
		Where(tableUnitUserName+".user_id = ?", unitDto.UserId).
		Where(tableUnitUserName + ".deleted = 0").
		Where(tableUnitName + ".deleted = 0")
    // other code ...

	err = query.Select(tableUnitName + ".id").Count(&count).Error
	// other code ...

	err = query.Select(selectStr).Order(tableUnitName + ".created_by," + tableUnitName + ".sort").Find(&listData).Error
	return
}
```
2. 如果不区分多用户体系，但是可以共用逻辑，则直接在`apps\common\models_ar`创建（如附件：`apps\common\models_ar\file.go`，编写curd逻辑）。




### 步骤 4：定义路由
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

### 步骤 4：创建 Controller
文件路径：`apps/admin_plat/controllers/system/unit.go`
- **结构体定义**：嵌入 `commonControllers.AdminBaseController`
- **Swagger 注解**：包含 `@Summary`, `@Description`, `@Tags`, `@Param`, `@Success`, `@Router`
- **方法实现**：解析参数，调用 Service，返回 JSON
接受参数如果需要自定义dto传参，则在`apps/common/dto_vo/<dtoType>/<Name>.go`

```go
package system

import (
    "WenBeego/apps/admin_plat/services/system"
    commonControllers "WenBeego/apps/common/controller"
    "WenBeego/apps/common/dto_vo/page"
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

### 步骤 6：创建 Service
1.若curd业务需要支持多用户体系逻辑，`apps/<商户或者平台>/services/system/unit.go`调用共用`apps\common\services\system\system-unit.go`的GetUnitList(unitDto page_dto.SystemUnitListReqDto)方法，再在GetUnitList()里面调用‘步骤 3’ base_ar下`apps\common\models_ar\base_ar\common_unit.go`泛型函数；
反正涉及到共用的
- 文件路径：`apps/admin_plat/services/system/unit.go`
```go
package system

import (
    "WenBeego/apps/common/dto_vo"
    "WenBeego/apps/common/dto_vo/page"
    CommonSystem "WenBeego/apps/common/services/system"
)

type UnitService struct {
    commonSystemUnit CommonSystem.Unit
}

func (s *UnitService) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (dto.RespDataListDto, error) {
    return s.commonSystemUnit.GetUnitList(unitDto)
}
```
- 文件路径 `apps\common\services\system\system-unit.go`
```go
package system

import (
    "WenBeego/apps/common/dto_vo"
    "WenBeego/apps/common/dto_vo/page"
    "WenBeego/apps/common/helper"
    "WenBeego/apps/common/models/base_model"
)

type UnitAr struct{}

func (s *UnitAr) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (dto.RespDataListDto, error) {
    var data []base_model.Unit
    var count int64
    
	switch unitDto.ModuleName {        
	case "admin_plat":
        // ... 平台组织数据库查询逻辑
		data, count, err = base_ar.GetUnitListByUserId(unitDto, &models.Plat{}, &models.PlatUser{})
	case "admin_mchnt":
        // ... 商户组织数据库查询逻辑
		data, count, err = base_ar.GetUnitListByUserId(unitDto, &models.Mchnt{}, &models.MchntUser{})
	default:
		err = errors.New("模块名称错误")
	}
    
    return helper.GetRespDataListDto(unitDto.PageSize, unitDto.CurrentPage, count, data)
}
```
2.若不管多用户体系与否，curd逻辑共用，那么`apps/<平台>/services/xxx/xxx.go`直接调用共用ActiveRecord 层（`apps\common\models_ar\file.go`）的方法返回数据即可

3. services层，如果存在多个数据库操作，必须使用事务, 如：
```go 

func (s *MenuService) Save(baseParamDto dto.BaseParamDto, data1Dto xxx1, data2Dto xxx2) (data map[string]string, err error) {
    // other code ...
	err = global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		err := xxxx1.Save(tx, baseParamDto, &data1Dto)
		if err != nil {
			return err
		}
		err = = xxxx2.Save(tx, baseParamDto, &data2Dto)
		return err
	})
	return
}
```


### 步骤 7：编写单元测试
测试文件位于 `tests/xxx_test.go`，调用 Service 层方法验证。
> ✅ **检查点**：每个步骤完成后，确认文件路径、包名、导入路径正确。

### 步骤 8：结束
- 提示用户执行命令`swag init -g cmd/http/main.go  --output apps/swagger `生成文档
- 提示用户执行命令`go mod tidy`更新依赖包
- 提示用户执行命令`go build -o app`编译项目
- 提示用户，添加权限菜单




## 新增MQ任务流程
### 1. 创建任务名称：
    在`apps\common\global\constant\mq.go`，添加任务名称为常量`MQ_XXXXXX`
### 2. 调用生产者：
    ```go
    args := []tasks.Arg{{Name: "action", Type: "string", Value: dataStr}}
    result, err := (&MqClient{}).SendTask(constant.MQ_XXXXXX, args)
    ```
### 4. 创建消费者：
    `apps\mq_task`目录结构是mvc类型，在`apps\mq_task\controllers`目录下创建消费者文件，并实现`ActionXXXXXX`方法，返回`error`类型。    

### 3. 添加消费者路由（方便统一管理）：
    在`routers\mq_task_router.go` GetMqTasks()函数内，添加路由切片。
    ```
    import "WenBeego/apps/mq_task/controllers"
    TasksList = append(TasksList, MqTasks{Name: string(constant.MQ_XXXXXX), CallBack: (*controllers.XXXXXX).ActionXXXXXX})
    ```



## 新增crontab任务流程
### 1. 创建任务名称：
    在`apps\common\global\constant\crontab.go`，添加任务名称为常量`CRON_XXXXXX`
### 2. 创建任务：
    `apps\cron_task`目录结构是mvc类型，在`apps\cron_task\controllers`目录下创建任务文件，并实现`ActionXXXXXX`方法。
### 3. 添加任务路由（方便统一管理）：
    在`routers\crontab_task\router.go` GetCronTasks()函数内，添加路由切片。
    ```go
    TasksList = append(TasksList, CronTasks{
		Name:     string(constant.CRON_XXXXXX),
		NameText: "生日提醒",
		CallBack: func() {
			(&controllers.XXX{}).ActionXXXXXX()
		},
	})
    ```


## 框架全局变量与工具
全局变量集中在 `apps/common/global`，可按需导入。

### 常量与路径
```go
import "WenBeego/apps/common/global/constant"
import "WenBeego/apps/common/global"

fmt.Println(constant.MODULE_NAME)                // 常量
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
import "WenBeego/apps/common/global/constant"
args := []tasks.Arg{{Name: "action", Type: "string", Value: dataStr}}
result, err := (&MqClient{}).SendTask(constant.MQ_XXXXXX, args)
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

## 常用DTO 
### 基本参数dto 
- `apps\common\dto\base_param.go`：包含http请求`ReqDataListDto`、响应的参数`Response`；包含系统业务必须基本参数`BaseParamDto`, BaseParamDto包含：请求头Host，应用模块名ModuleName，组织单位UnitId, 主体用户UserId，组织单位下独立的用户UnitUserId
- 自定义参数dto（如：`apps\common\dto\page_dto\system-user.go`）, system-user.go里面结构体直接继承base_param.go 即可，自行阅读源码


## ⚠️ 安全与注意事项
- **数据库写操作**（如 `global.GetWriteDb().Create`）应确保在事务中或数据已验证。
- **文件操作**（删除、覆盖）前必须先请求用户确认。
- **生产环境**避免打印敏感信息（密码、Token）。
- **权限校验**：涉及组织单位或角色的接口，调用 `helper.CheckUserHasUnit` 等函数进行鉴权。

## 外部参考文件指引
- 完整目录结构：`doc/目录结构.md`（按需读取）
- 数据库表定义：`doc/数据库表结构.sql`（按需读取特定表）
- 助手函数源码：`apps/common/helper/`（按需查看）
