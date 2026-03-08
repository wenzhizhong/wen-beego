---
name: wen-beego
description: wen-beego框架开发。
license: MIT
compatibility: OpenCode >= 1.0
---

# wen-beego
## 项目简介
wen-beego框架， 是基于beego框架的web开发框架。系统包含平台端（apps/admin_plat）、商户端（apps/admin_mchnt）,用于二次开发


## 项目二次开发：

### 测试参数
    账号：15912345678
    密码：G72+shD3^6

### 运行
    - 运行http服务：go run cmd/http/main.go； 
    - 运行mq服务：go run cmd/mq/main.go；

### 新增接口举例
    - 可以打开路由文件：
        routers\admin_mchnt_router.go
        routers\admin_plat_router.go
        routers\mq_task_router.go
        routers\crontab_task\router.go
    - 移动端等api模块，尚未新增；已经新增，则直接编辑
    - 以新增平台端'获取组织单位'接口`/admin_plat/system-unit/get`为例：
        1. 创建、编辑函数，在函数里面，把路由`system-unit/get`放入切片，再返回所有路由切片。
        2. 这个函数命名要求：（平台类型+模块类型+slices）;
        切片元素要求：beego.NSCtrl<请求类型：Get、Post、Delete、Put>("/system-unit/get", (*控制器包名.组织单位Controller).操作Action)；
        示例：
        ```
            func platSystemSlices() []beego.LinkNamespace {
                return []beego.LinkNamespace{
                    beego.NSCtrlGet("/system-unit/get", (*adminSystem.UnitController).Get),
                }
            }
        ```
        3. 创建/修改mvc代码文件，并放入对应目录下
        
        3-1. controller控制器：`apps\admin_plat\controllers\system\system-unit.go`
        控制器命名要求：`apps\<系统、api类型>_<主体类型>\controllers\模块名称\<前缀>-模块名称.go`
        控制器采用‘类’写法： type 模块Controller struct{继承其他类}
        控制器方法注释：swagger注释
        控制器方法内：请求参数，可以接受参数到dto， 返回参数，如果返回json, 获取返回数据格式c.Data["json"] = helper.Response(200, "success", data)，然后调用c.ServeJSON()

        示例：
        ```
        package system

        // 系统管理-内部组织管理
        import (
            systemService "WenBeego/apps/admin_plat/services/system"
            commonControllers "WenBeego/apps/common/controller"
            "WenBeego/apps/common/dto/page_dto"
            "WenBeego/apps/common/dto/unit_dto"
            "WenBeego/apps/common/helper"
            "errors"
        )

        type UnitController struct {
            commonControllers.AdminBaseController
            UnitService systemService.UnitService
        }

        // 系统管理-获取内部组织管理
        // @Summary 获取内部组织管理
        // @Description 获取内部组织管理
        // @Tags 系统管理-内部组织管理
        // @Accept application/json
        // @Produce application/json
        // @Param parentUnitId query string true "父级ID"
        // @Success 200 {object} dto.RespDataListDto
        // @Router /admin_plat/system-unit/get [get]
        func (c *UnitController) Get() {
            baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
            reqDataListDto, err2 := helper.GetReqDataListDto(&c.Controller)
            if err != nil {
                c.Data["json"] = helper.Response(500, err.Error(), nil)
                c.ServeJSON()
                return
            }
            if err2 != nil {
                c.Data["json"] = helper.Response(500, err2.Error(), nil)
                c.ServeJSON()
                return
            }
            unitDto := page_dto.SystemUnitListReqDto{}
            unitDto.BaseParamDto = baseParamDto
            unitDto.ReqDataListDto = reqDataListDto
            unitDto.Name = c.GetString("name")
            unitDto.Code = c.GetString("code")
            unitDto.Status, _ = c.GetInt("status")

            data, err := c.UnitService.GetUnitList(unitDto)
            if err != nil {
                c.Data["json"] = helper.Response(500, err.Error(), nil)
                c.ServeJSON()
                return
            }
            c.Data["json"] = helper.Response(200, "success", data)
            c.ServeJSON()
        }
        ```

        3-2. service服务：`apps\admin_plat\services\system\system-unit.go`
        servcie文件命名要求：`<前缀>-模块名称.go`
        service文件也是采用‘类’写法： type <模块>Service struct {继承其他类}
        service方法返回的数据，在控制器直接传递到helper.Response(200, "success", data)函数；最终前端json为：{"code":200,"message": "success", "data": data}，所以服务端data一定要规范为一个map、结构体；比如 data := helper.GetRespDataListDto(unitDto.PageSize, unitDto.CurrentPage, count, data)，那么返回data就是一个包含分页+数据的结构体
        示例：
        ```
        package system

        import (
            "WenBeego/apps/common/dto"
            "WenBeego/apps/common/dto/page_dto"
            "WenBeego/apps/common/dto/unit_dto"
            CommonSystem "WenBeego/apps/common/services/system"
        )

        type UnitService struct {
            commonSystemUnit CommonSystem.Unit
        }

        // 系统管理- 获取内部组织列表
        func (s *UnitService) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (dto.RespDataListDto, error) {
            return s.commonSystemUnit.GetUnitList(unitDto)
        }
        ```
        3-3. model_ar数据访问层：`apps\common\services\system\system-unit.go`；如果平台端、商户端可以通过泛型共用，那么在`apps\common\services`目录下写service业务逻辑，反之需要平台端（`apps\admin_plat\services`）、商户端（`apps\admin_mchnt\services`）分别处理，其他客户端service设计逻辑也是这个规则。
        示例：
        ```
        package system

        import (
            "WenBeego/apps/common/dto"
            "WenBeego/apps/common/dto/page_dto"
            "WenBeego/apps/common/dto/unit_dto"
            "WenBeego/apps/common/global"
            "WenBeego/apps/common/helper"
            "WenBeego/apps/common/middleware/business_store"
            "WenBeego/apps/common/models"
            "WenBeego/apps/common/models/base_model"
            "WenBeego/apps/common/models/itf"
            "WenBeego/apps/common/models_ar/base_ar"
            "errors"

            "gorm.io/gorm"
        )

        type Unit struct {
        }

        // 系统管理-获取内部组织列表
        func (s *Unit) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (resultDto dto.RespDataListDto, err error) {
            data := make([]base_model.Unit, 0)
            var count int64 = 0

            switch unitDto.ModuleName {
            case "admin_plat":
                data, count, err = base_ar.GetUnitListByUserId(unitDto, &models.Plat{}, &models.PlatUser{})
            case "mchnt_plat":
                data, count, err = base_ar.GetUnitListByUserId(unitDto, &models.Mchnt{}, &models.MchntUser{})
            default:
                err = errors.New("GetUnitList:模块名称错误")
            }
            if err != nil {
                return
            }
            for k, v := range data {
                tmpLogo, err1 := helper.LocalFileSign(unitDto.Host, v.Logo)
                tmpPath, err2 := helper.LocalFileSign(unitDto.Host, v.License)
                if err1 != nil || err2 != nil {
                    errStr := helper.Ternary(err1 != nil, err1.Error(), err2.Error())
                    global.Log.Error(errStr)
                    continue
                }
                data[k].LogoLink = tmpLogo
                data[k].LicenseLink = tmpPath
            }
            resultDto, err = helper.GetRespDataListDto(unitDto.PageSize, unitDto.CurrentPage, count, data)
            return
        }
        ```
    - test单元测试：`.\tests\xxx_test.go`，调用sevice层的函数、方法

### 框架全局变量使用
    - 全局变量包路径：`apps\common\global`
    - 使用示例：
        1. 常量示例
        ```
        import "WenBeego/apps/common/global/constant"
        import "fmt"

        fmt.println(constant.ADMIN) // 输出：admin
        ```
        2. 路径示例
        ```
        import "fmt"
        import "WenBeego/apps/common/global"

        fmt.Println(RootPath)
        fmt.Println(AppDir)
        fmt.Println(ConfigDir)
        fmt.Println(StaticDir)
        fmt.Println(RoutersDir)
        fmt.Println(TempDir)
        fmt.Println(UploadsDir)
        ```

        2. 配置文件示例
        ```
        import "WenBeego/apps/common/global"

        key, err := global.GetConfigDiy(xxxKey)
        ```

        3. 缓存示例
        ```
        import "WenBeego/apps/common/global"
        
        data, err = global.RedisCache.Get(ctx, key).Result()
        ```

        4. 日志示例
        ```
        import "WenBeego/apps/common/global"
        Log := global.Log
        Log.Info("info log")
        Log.Error("error log")
        Log.Debug("debug log")
        Log.Warn("warn log")
        Log.Trace("trace log")
        ````

        5. 数据库示例
        ```
        import "WenBeego/apps/common/global"

        global.GetReadDb() // 获取读数据库，gorm.DB实例
        global.GetWriteDb() // 获取写数据库，gorm.DB实例
        ```

        6. mq示例, 直接调用已封装好的函数：
        ```
        args := []tasks.Arg{{Name: "actionXxxx", Type: "string", Value: dataStr}}
        result, err := (&MqClient{}).SendTask("ApiLog.ActionSaveToDb", args)
        ```
### 助手函数
    - 助手函数包路径：`apps\common\helper`
    - apps\common\helper\array.go:
        // 判断字符串是否在数组中
        func InArray(val string, array interface{}) (bool, error) 

        // 数组合并
        func ArrayMerge(arr1, arr2 []string) []string 

    - apps\common\helper\branca.go:
        // encode
        func BrancaEncode(data BrancaData, moduleName string) (string, error) 

        // decode
        func BrancaDecode(needDecodeString string, moduleName string) (BrancaData, error) 

        // Brance key
        func getBranceKey(moduleName string) (string, error) 

    - apps\common\helper\conver.go:
        // interface{} 转 map[string]interface{}
        func Interface2MapInterface(i interface{}) (map[string]interface{}, error) 

        // map[string]interface{} 转 map[string]string
        func MapInterface2MapString(i map[string]interface{}) (map[string]string, error) 

        // interface{} 转 int64
        func Interface2Int64(value interface{}) (int64, error) 

        func String2Int64(value string) (int64, error) 

        func String2Int(value string) (int, error) 

        func Int2String(value int) string 

        func Int642String(value int64) string 

        func Float642String(value float64) string 

    - apps\common\helper\datetime.go:
        // 获取时间戳
        func GetTimestamp(datetime ...string) int64 

        // 获取日期时间戳
        func GetDateStamp(date ...string) int64 

        // 获取时间
        func GetTime() time.Time 

        // 获取时间字符串
        func GetTimeString(format ...string) string 

        // 时间戳转时间
        func TimestampToTime(timestamp int64, format ...string) string 

    - apps\common\helper\filepath.go:
        // 路径存在
        func PathIsExist(path string) bool 
        // 创建目录
        func MkdirAll(path string) error 

        // 本地文件访问签名
        func LocalFileSign(host, filePath string) (string, error) 

        // 本地文件访问签名校验
        func LocalFileSignCheck(sign string) (bool, error) 

    - apps\common\helper\framework.go:
        // ASCII 艺术字符
        func Get_ASCII_ArtisticCharacters() (string, error) 

        // 输出艺术字logo到终端
        func Output_ASCII_ArtisticCharacters() 

        // 获取系统应用名称
        func AppName() (string, error) 

        // 获取框架版本
        func AppVersion() (string, error) 

        // 获取框架runmode
        func AppRunmode() (string, error) 

        // 检查运行模式
        func CheckRunMode(runMode string) bool 

        // 从路由路径解析模块名
        func ParseModuleFromRoute(path string) string 

    - apps\common\helper\gorm.go:
        // 数据库记录不存在
        func DbNotFound(err error) bool 

    - apps\common\helper\map.go:
        // 获取map的keys
        func GetMapKeys[K comparable, V any](m map[K]V) []K 

        //  获取map的values
        func GetMapValues[K comparable, V any](m map[K]V) []any 

    - apps\common\helper\password.go:
        // 获取随机密码
        func GetRandomPassword(length int) (string, error) 

        // 校验密码规则
        func CheckPasswordRule(password string) error 

        // 密码强度校验，顺序字符(相邻字符ascii相差1)超过5个，则返回弱密码错误
        func CheckPasswordRuleSafe(password string) error 

        // 生成加密密码
        func GenerateCryptPassword(password string) (cryptPassword string, err error) 

        // 验证密码
        func CompareCryptPassword(password string, cryptPassword string) bool 

    - apps\common\helper\redis.go:
        // redis key
        func GetCustomRedisKey(key string) (string, error) 

        // redis put value to redis
        func RedisPut(key string, value interface{}, timeoutAfter int) error 

        // redis get value from redis
        func RedisGet(key string) (string, error) 

        // redis delete value from redis
        func RedisDel(key string) error 

    - apps\common\helper\request.go:
        // 获取请求的token
        func GetReqToken(ctx beecontext.Context) string 

        // 获取请求接口
        func GetReqUrl(ctx beecontext.Context) string 

        // 获取请求方法
        func GetReqMethod(ctx beecontext.Context) string 

        // 获取请求协议
        func GetReqScheme(ctx beecontext.Context) string 

        // 获取请求主机
        func GetReqHost(ctx beecontext.Context) string 

        // 请求IP
        func GetReqIp(ctx beecontext.Context) string 

        // 请求body
        func GetReqBody[T any](ctx *beecontext.Context) (T, error) 

        // 设置并获取基本参数
        func GetBaseParamDto(ctx *beecontext.Context, moduleName string) (dto.BaseParamDto, error) 

        // 设置并返回请求页面数据列表参数
        func GetReqDataListDto(ctxCtrl *web.Controller) (dto.ReqDataListDto, error) 

    - apps\common\helper\response.go:
        // 接口响应函数
        func Response(code int, message string, data interface{}) dto.Response 

        // 设置并返回页面列表数据
        func GetRespDataListDto(pageSize int, currentPage int, total int64, list interface{}) (dto.RespDataListDto, error) 

    - apps\common\helper\string.go:
        // 删除字符串中的空格
        func DeleteSpace(str string) string 

    - apps\common\helper\utils.go:
        // Ternary 泛型三元运算符
        func Ternary[T any](condition bool, trueVal, falseVal T) T 

        // 指针版本（避免值复制）
        func TernaryPtr[T any](condition bool, trueVal, falseVal *T) *T 

        // 字符串专用
        func StringTernary(condition bool, trueVal, falseVal string) string 

        // 整数专用
        func IntTernary(condition bool, trueVal, falseVal int) int 

        // 布尔专用
        func BoolTernary(condition bool, trueVal, falseVal bool) bool 

        // 判断手机号
        func IsCellPhone(phone string) bool 

        // 判断邮箱
        func IsEmail(email string) bool 

        // 判断中文
        func IsChinese(str string) bool 

        // 解析字符串模板
        func ParseStringTpl(tpl string, data any) (str string, err error) 

        // 判断是否是官方平台id
        func IsOfficial(moduleName, unitId string) (bool, error) 

        // 判断是否是管理员
        func IsAdmin(moduleName string, unitUserId string) (bool, error) 

        // 用户是否有组织单位权限
        func CheckUserHasUnit(moduleName, userId string, requiredUnitIds []string) (bool, error) 

        // 获取管理员用户
        func getAdminData[RoleClassifyModel itf.RoleClassifyItf, RoleModel itf.RoleItf, UserRoleModel itf.UserRoleItf](unitUserId string, roleClassify RoleClassifyModel, role RoleModel, userRoleModel UserRoleModel) (bool, error) 

        // 获取uuid
        func GetUuid() (string, error) 

        // 获取md5值
        func Md5(str string) string 

        // 获取refresh token 缓存key
        func getRefreshTokenKey(brancaToken string, refreshToken string) string 

        // 获取refresh token
        func GetRefreshToken(moduleName string, brancaToken string, userId string) (string, error) 

        // 验证refresh token
        func VerifyRefreshToken(brancaToken string, refreshToken string) (result bool, userId string, err error) 

        // 删除refresh token
        func DelRefreshToken(brancaToken string, refreshToken string) 


    - apps\common\helper\tree\tree-converter.go:
        // WithRootParentID 设置根节点的父ID值
        func (tc *TreeConverter) WithRootParentID(rootParentID interface{}) *TreeConverter 

        // WithKeepOrder 设置是否保持原始顺序
        func (tc *TreeConverter) WithKeepOrder(keepOrder bool) *TreeConverter 

        // ConvertToTree 将数组转换为树形结构
        func (tc *TreeConverter) ConvertToTree(nodes []TreeNode) ([]TreeNode, error) 

        // isRootNode 判断是否为根节点
        func (tc *TreeConverter) isRootNode(node TreeNode) bool 

        // buildTree 递归构建树
        func (tc *TreeConverter) buildTree(node TreeNode, childrenMap map[interface{}][]TreeNode, processed map[interface{}]bool) 

        // sortRootsByOriginalOrder 按照原始顺序排序根节点
        func (tc *TreeConverter) sortRootsByOriginalOrder(roots []TreeNode, originalNodes []TreeNode) []TreeNode 

        // sortChildrenByOriginalOrder 按照原始顺序排序子节点
        func (tc *TreeConverter) sortChildrenByOriginalOrder(children []TreeNode, originalOrder []TreeNode) []TreeNode 

        // ConvertSliceToTree 通用切片转树形结构（使用反射，不需要实现TreeNode接口）
        func ConvertSliceToTree(slice interface{}, idField, parentIdField, childrenField string) (interface{}, error) 

        // buildTreeWithReflection 递归构建树（反射版本）
        func buildTreeWithReflection(node reflect.Value, childrenMap map[interface{}][]reflect.Value, processed map[interface{}]bool, childrenField string) 

        // 辅助函数：获取字段值（处理指针）
        func getFieldValue(v reflect.Value, fieldName string) reflect.Value 

        // isZeroValue 判断是否为零值
        func isZeroValue(v reflect.Value) bool 

    - apps\common\helper\tree\tree.go:
        // BuildTree 快速构建树形结构（使用TreeNode接口）
        func BuildTree(nodes []TreeNode) ([]TreeNode, error) 

        // BuildTreeWithRootID 使用指定的根节点父ID构建树
        func BuildTreeWithRootID(nodes []TreeNode, rootParentID interface{}) ([]TreeNode, error) 

        // NewTreeConverter 创建树形转换器
        func NewTreeConverter() *TreeConverter 


## 项目目录结构：
```
├── LICENSE
├── apps
│   ├── admin_mchnt 商户管理mvc模块
│   │   ├── controllers 控制器
│   │   │   ├── auth 认证相关模块
│   │   │   │   ├── auth-mchnt.go
│   │   │   │   ├── auth-menu.go
│   │   │   │   ├── auth.go
│   │   │   │   └── params.go
│   │   │   ├── system 系统管理模块
│   │   │   │   ├── system-dept.go
│   │   │   │   ├── system-role.go
│   │   │   │   ├── system-unit.go
│   │   │   │   └── system-user.go
│   │   │   └── upload 上传模块
│   │   │       └── upload.go
│   │   ├── models 商户特有数据模型目录
│   │   ├── services 商户业务逻辑目录
│   │   │   ├── auth 认证相关模块
│   │   │   │   ├── auth-mchnt.go
│   │   │   │   ├── auth-menu.go
│   │   │   │   └── auth.go
│   │   │   ├── system 系统管理模块
│   │   │   │   ├── system-dept.go
│   │   │   │   ├── system-role.go
│   │   │   │   ├── system-unit.go
│   │   │   │   └── system-user.go
│   │   │   └── upload 上传模块
│   │   │       └── upload.go
│   │   └── views 视图目录
│   │       └── auth
│   │           └── index.html
│   ├── admin_plat 平台管理mvc模块
│   │   ├── controllers 控制器
│   │   │   ├── auth 认证相关模块
│   │   │   │   ├── auth-menu.go
│   │   │   │   ├── auth-plat.go
│   │   │   │   ├── auth.go
│   │   │   │   └── params.go
│   │   │   ├── monitor 系统监控模块
│   │   │   │   ├── monitor-cron-log.go
│   │   │   │   └── monitor-cron.go
│   │   │   ├── system 平台系统管理模块
│   │   │   │   ├── system-dept.go
│   │   │   │   ├── system-menu-mchnt.go
│   │   │   │   ├── system-menu-plat.go
│   │   │   │   ├── system-role.go
│   │   │   │   ├── system-unit.go
│   │   │   │   └── system-user.go
│   │   │   ├── system_mchnt 平台管理商户端的系统模块
│   │   │   │   ├── system-dept.go
│   │   │   │   ├── system-role.go
│   │   │   │   ├── system-unit.go
│   │   │   │   └── system-user.go
│   │   │   └── upload 平台端上传模块
│   │   │       └── upload.go
│   │   ├── models 平台数据模型目录
│   │   ├── models_ar 平台数据模型AR目录, 模型AR是单独从model抽出来的，curd在AR中实现
│   │   │   ├── monitor-cron-log.go
│   │   │   └── monitor-cron.go
│   │   ├── services  平台servers层目录
│   │   │   ├── auth 认证服务模块
│   │   │   │   ├── auth-menu.go
│   │   │   │   ├── auth-plat.go
│   │   │   │   └── auth.go
│   │   │   ├── monitor 监控服务模块
│   │   │   │   ├── monitor-cron-log.go
│   │   │   │   └── monitor-cron.go
│   │   │   ├── system 系统服务模块
│   │   │   │   ├── system-dept.go
│   │   │   │   ├── system-menu.go
│   │   │   │   ├── system-role.go
│   │   │   │   ├── system-unit.go
│   │   │   │   └── system-user.go
│   │   │   ├── system_mchnt 平台管理商户端的系统模块
│   │   │   │   ├── system-dept.go
│   │   │   │   ├── system-role.go
│   │   │   │   ├── system-unit.go
│   │   │   │   └── system-user.go
│   │   │   └── upload 文件上传模块
│   │   │       └── upload.go
│   │   └── views
│   │       └── auth
│   │           └── index.html
│   ├── common  公用目录
│   │   ├── controller 公用控制器基类
│   │   │   ├── admin_base_controller.go
│   │   │   ├── base_controller.go
│   │   │   └── index_base_controller.go
│   │   ├── dto 公用DTO目录
│   │   │   ├── auth_dto 认证模块dto
│   │   │   │   ├── auth.go
│   │   │   │   ├── menu.go
│   │   │   │   └── unit.go
│   │   │   ├── base_param.go 获取请求参数dto
│   │   │   ├── cron_dto 定时任务模块dto
│   │   │   │   └── unit_cron.go
│   │   │   ├── dept_dto 部门模块dto
│   │   │   │   └── unit_dept.go
│   │   │   ├── menu_dto 菜单模块dto
│   │   │   │   └── menu.go
│   │   │   ├── mq_dto mq模块dto
│   │   │   │   └── mq_api_log.go
│   │   │   ├── page_dto 后台页面数据列表dto
│   │   │   │   ├── monitor-cron-log.go
│   │   │   │   ├── monitor-cron.go
│   │   │   │   ├── system-dept.go
│   │   │   │   ├── system-menu.go
│   │   │   │   ├── system-role.go
│   │   │   │   ├── system-unit.go
│   │   │   │   └── system-user.go
│   │   │   ├── role_dto 角色模块dto
│   │   │   │   └── unit_role.go
│   │   │   ├── unit_dto 组织单位模块dto
│   │   │   │   ├── unit.go
│   │   │   │   └── unit_user.go
│   │   │   ├── upload_dto 文件上传模块dto
│   │   │   │   └── upload.go
│   │   │   └── user_dto 用户模块dto
│   │   │       └── user.go
│   │   ├── global 公用全局变量目录，系统启动、服务初始化时设置
│   │   │   ├── config.go 全局配置
│   │   │   ├── constant 常量
│   │   │   │   └── constant.go
│   │   │   ├── db.go 数据库实例全局变量
│   │   │   ├── frame_path.go 系统路径全局变量
│   │   │   ├── log.go 日志全局变量
│   │   │   ├── mq.go 消息队列实例全局变量
│   │   │   └── redis.go redis实例全局变量
│   │   ├── helper 公用助手目录
│   │   │   ├── array.go 数组助手
│   │   │   ├── branca.go branca加密
│   │   │   ├── conver.go 转换辅助函数
│   │   │   ├── datetime.go 时间助手
│   │   │   ├── filepath.go 路径助手
│   │   │   ├── framework.go 框架助手
│   │   │   ├── gorm.go gorm数据库助手
│   │   │   ├── map.go 映射助手
│   │   │   ├── password.go 密码助手
│   │   │   ├── redis.go redis助手
│   │   │   ├── request.go http请求助手
│   │   │   ├── response.go http响应助手
│   │   │   ├── string.go 字符串助手
│   │   │   ├── tree 树形结构助手
│   │   │   │   ├── base-node.go
│   │   │   │   ├── tree-converter.go
│   │   │   │   └── tree.go
│   │   │   └── utils.go 杂用函数
│   │   ├── middleware 公用中间件目录
│   │   │   ├── access.go 框架访问控制
│   │   │   ├── auth.go 框架认证
│   │   │   ├── business_store 业务存储中间件
│   │   │   │   └── check-auth.go
│   │   │   ├── captcha 验证码中间件
│   │   │   │   ├── base64CaptchaRedisStore.go
│   │   │   │   └── captcha.go
│   │   │   ├── crontab 定时任务中间件
│   │   │   │   ├── crontab.go
│   │   │   │   └── distributed_lock.go
│   │   │   ├── database 数据库中间件
│   │   │   │   ├── db.go
│   │   │   │   ├── logger-custom.go
│   │   │   │   └── logger-impl.go
│   │   │   ├── log.go 日志中间件 
│   │   │   ├── mqClient.go mqClient中间件
│   │   │   ├── mqServer.go mqServer中间件
│   │   │   └── redis.go redis中间件
│   │   ├── models 公用模型目录
│   │   │   ├── base_model 模型基础目录
│   │   │   ├── itf 模型接口目录
│   │   │   │   ├── api_statistics_interface.go
│   │   │   │   ├── dept_interface.go
│   │   │   │   ├── menu_interface.go
│   │   │   │   ├── menu_map_interface.go
│   │   │   │   ├── role_classify_interface.go
│   │   │   │   ├── role_interface.go
│   │   │   │   ├── role_menu_interface.go
│   │   │   │   ├── unit_cron_interface.go
│   │   │   │   ├── unit_cron_log_interface.go
│   │   │   │   ├── unit_interface.go
│   │   │   │   ├── unit_user_interface.go
│   │   │   │   ├── user_dept_interface.go
│   │   │   │   ├── user_profile_interface.go
│   │   │   │   └── user_role_interface.go
│   │   │   ├── mchnt.go 商户模型
│   │   │   ├── mchnt_api_statistics.go 商户接口统计
│   │   │   ├── mchnt_customer.go 商户部门
│   │   │   ├── mchnt_dept.go 商户部门
│   │   │   ├── mchnt_menu.go 商户菜单
│   │   │   ├── mchnt_menu_map.go 商户菜单映射
│   │   │   ├── mchnt_role.go 商户角色
│   │   │   ├── mchnt_role_classify.go 商户角色分类
│   │   │   ├── mchnt_role_menu.go 商户角色菜单
│   │   │   ├── mchnt_user.go 商户用户
│   │   │   ├── mchnt_user_dept.go 商户用户部门
│   │   │   ├── mchnt_user_profile.go 商户用户信息
│   │   │   ├── mchnt_user_role.go 商户用户角色
│   │   │   ├── plat.go 平台
│   │   │   ├── plat_api_statistics.go 平台接口统计
│   │   │   ├── plat_cron.go 平台定时任务
│   │   │   ├── plat_cron_log.go 平台定时任务日志
│   │   │   ├── plat_dept.go 平台部门
│   │   │   ├── plat_menu.go 平台菜单
│   │   │   ├── plat_menu_map.go 平台菜单权限
│   │   │   ├── plat_menu_map_view.go 平台菜单权限视图
│   │   │   ├── plat_menu_view.go 平台菜单视图
│   │   │   ├── plat_role.go 平台角色
│   │   │   ├── plat_role_classify.go 平台角色分类
│   │   │   ├── plat_role_menu.go 平台角色菜单
│   │   │   ├── plat_user.go 平台用户
│   │   │   ├── plat_user_dept.go 平台用户部门
│   │   │   ├── plat_user_profile.go 平台用户信息
│   │   │   ├── plat_user_role.go 平台用户角色
│   │   │   ├── user.go 用户
│   │   │   └── user_profile.go 用户信息
│   │   ├── models_ar 公用模型AR目录，从models分离出来，curd在这里实现
│   │   │   ├── base_ar 通用的模型AR目录
│   │   │   │   ├── common_api_statistics.go
│   │   │   │   ├── common_dept.go
│   │   │   │   ├── common_menu.go
│   │   │   │   ├── common_menu_map.go
│   │   │   │   ├── common_role.go
│   │   │   │   ├── common_role_classify.go
│   │   │   │   ├── common_role_menu.go
│   │   │   │   ├── common_unit.go
│   │   │   │   ├── common_unit_user.go
│   │   │   │   ├── common_unit_user_dept.go
│   │   │   │   ├── common_unit_user_profile.go
│   │   │   │   ├── common_unit_user_role.go
│   │   │   │   ├── common_user.go
│   │   │   │   └── common_user_profile.go
│   │   │   ├── file.go 文件模型AR
│   │   │   ├── file_slice.go 文件分片模型AR
│   │   │   ├── mchnt.go 商户模型AR
│   │   │   ├── mchnt_api_statistics.go 商户API统计模型AR
│   │   │   ├── mchnt_customer.go 商户客户模型AR
│   │   │   ├── mchnt_menu.go 商户菜单模型AR
│   │   │   ├── mchnt_role.go 商户角色模型AR
│   │   │   ├── mchnt_role_classify.go 商户角色分类模型AR
│   │   │   ├── mchnt_role_menu.go 商户角色菜单模型AR
│   │   │   ├── mchnt_user.go 商户用户模型AR
│   │   │   ├── mchnt_user_role.go 商户用户角色模型AR
│   │   │   ├── plat.go 平台模型AR
│   │   │   ├── plat_api_statistics.go 平台接口统计模型AR
│   │   │   ├── plat_cron.go 平台计划任务模型AR
│   │   │   ├── plat_cron_log.go 平台计划任务日志模型AR
│   │   │   ├── plat_menu.go 平台菜单模型AR
│   │   │   ├── plat_menu_map_view.go 平台菜单映射视图模型AR
│   │   │   ├── plat_menu_view.go 平台菜单视图模型AR
│   │   │   ├── plat_role.go 平台角色模型AR
│   │   │   ├── plat_role_classify.go 平台角色分类模型AR
│   │   │   ├── plat_role_menu.go 平台角色菜单模型AR
│   │   │   ├── plat_user.go 平台用户模型AR
│   │   │   ├── plat_user_role.go 平台用户角色模型AR
│   │   │   ├── user.go 用户模型AR
│   │   │   └── user_profile.go 用户信息模型AR
│   │   └── services 公用service目录
│   │       ├── auth 后台通用认证service目录
│   │       ├── framework 框架的service目录
│   │       ├── system 后台通用系统管理service目录
│   │       └── upload 通用上传文件service目录
│   ├── cron_task crontab定时任务mvc目录
│   │   ├── controllers
│   │   ├── models
│   │   └── services
│   ├── index web端首页路由mvc目录
│   ├── mq_task mq任务消费mvc目录
│   └── swagger  swagger接口文档
├── cmd
│   ├── common
│   │   ├── init_source_service.go 入口文件初始化数据库
│   │   └── run_fefore.go 入口运行前的函数
│   ├── http
│   │   └── main.go http服务入口文件
│   └── mq
│       └── mainMq.go mq服务入口文件
├── conf
│   ├── app.conf 弃用
│   ├── app.yaml 开发环境
│   ├── prod.app.yaml 生产环境
│   └── test.app.yaml 测试环境
├── doc
│   ├── 生产环境部署.md
│   └── 系统手册.md
├── favicon.ico
├── go.mod
├── go.sum
├── readme-zh.md
├── readme.md
├── routers
│   ├── admin_mchnt_router.go 商户后台路由
│   ├── admin_plat_router.go 平台后台路由
│   ├── crontab_task crontab任务路由
│   │   └── router.go
│   ├── index_router.go web端首页路由
│   ├── mq_task_router.go mq任务路由
│   └── router.go 弃用
├── static 静态资源目录
├── temp 临时目录
├── tests 测试目录
├── uploads 文件上传目录
│   ├── private 私有文件，请求会校验权限
│   └── public 公共文件
├── vendor
└── wen-admin-docker
    ├── docker-compose.yml
```

### 项目数据库表结构

