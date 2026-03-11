---
name: wen-beego
description: 介绍 MyFramework 的核心设计思想、目录结构和基础约定，用于指导后续代码编写
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
    验证码：1234

### 运行
    - 运行http服务：go run cmd/http/main.go； 
    - 运行mq服务：go run cmd/mq/main.go；

### 项目目录结构：
- apps/admin_plat：平台端
- apps/admin_mchnt：商户端
- apps/api：api模块
- apps/cron_task：定时任务
- apps/common：公共模块
- apps/mq_task：mq任务消费模块
- apps/swagger：swagger文档
- routers/crontab_task/router.go：定时任务路由
- routers/admin_plat_router.go：平台端路由
- routers/admin_mchnt_router.go：商户端路由
- routers/mq_task_router.go：mq任务路由
- apps客户端的模块，都是mvc架构
- 详细目录，请读取项目根目录下的 `doc/目录结构.md` 文件。该文件包含了项目目录结构

### 项目数据库表结构
请阅读项目根目录下的 `doc/数据库表结构.sql` 文件。该文件包含了项目数据库表结构


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
        // @Success 200 {object} dto.RespDataListDto "返回结果"
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
            // 具体调用查询数据方法省略
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

        global.GetReadDb().Model(menuModel). Select(selectStr). Joins(joinStr). Where(tableMenuMap+".unit_id in ?", unitIds). Where(tableMenuMap+".deleted = 0"). Where("xxxxx"). Order("xxx asc"). Find(&dataList) // 获取读数据库，即gorm.DB实例
        global.GetWriteDb().Create(&data) // 获取写数据库，即gorm.DB实例
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

        // 是否开发环境
        func IsDevRunMode(runMode string) bool 

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
