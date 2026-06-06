package blocker

import (
	"WenBeego/apps/common/dto_vo/mq_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/global/app_error"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware/mq"
	"WenBeego/apps/common/services/framework"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/tasks"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"golang.org/x/time/rate"
)

var (
	limiter            = rate.NewLimiter(10, 20)
	cacheApiStatistics = make([]interface{}, 0)
	cacheMutex         = sync.Mutex{}
	maxBatchSize       = 500
	flushCacheInterval = 60 * time.Second
)

func init() {
	// 启动定期检查和批量处理
	go func() {
		ticker := time.NewTicker(flushCacheInterval)
		defer ticker.Stop()
		for range ticker.C {
			flushCacheIfNeeded()
		}
	}()
}

// 限制访问次数,防止dos攻击
func LimitTimes() web.FilterFunc {
	return func(ctx *beecontext.Context) {
		if !limiter.Allow() {
			fmt.Println("Too Many Requests")
			global.Log.Warn("Too Many Requests")
			http.Error(ctx.ResponseWriter, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		uri := ctx.Request.URL.Path
		if uri != "" && strings.HasPrefix(uri, "/uploads/private/") {
			signStr := ctx.Request.URL.Query().Get("sign")
			_, err := helper.LocalFileSignCheck(signStr)
			if err != nil {
				http.Error(ctx.ResponseWriter, "sign error", http.StatusBadRequest)
				return
			}
		}
	}
}

// api 请求前
func RouterBefore() web.FilterFunc {
	return func(ctx *beecontext.Context) {
		moduleName := helper.ParseModuleFromRoute(ctx.Request.URL.Path)
		ctx.Input.SetData(constant.MODULE_NAME, moduleName)
		ctx.Input.SetData(constant.CTX_ERROR_KEY, nil)

		dealSignAndEncrypt(ctx)
	}
}

// api 请求后
func RouterAfter(whiteApiList *[]string, authApiList *[]string) web.FilterFunc {
	return func(ctx *beecontext.Context) {
		globalLogToFile(ctx)
		statisticsApiLog(ctx, whiteApiList, authApiList)
	}
}

// 非默认error类型，全局写入文件日志
func globalLogToFile(ctx *beecontext.Context) {
	errInfo := ctx.Input.GetData(constant.CTX_ERROR_KEY)
	if errInfo == nil {
		return
	}
	if be, ok := errInfo.(*app_error.BaseError); ok {
		logString := fmt.Sprintf("Error: %v\n  Trace:\n%v\n", be.Err, be.Trace)
		global.Log.Error(logString)
		return
	}
}

// api 统计
func statisticsApiLog(ctx *beecontext.Context, whiteApiList *[]string, authApiList *[]string) {
	modules, _ := global.GetConfigDiy("logToDbModules")
	moduleName := helper.ParseModuleFromRoute(ctx.Request.URL.Path)
	tmpIgnoreArr := helper.ArrayMerge(*whiteApiList, *authApiList)
	headerBaseInfo := getBaseInfo(ctx)
	isInArray, _ := helper.InArray(headerBaseInfo.url, tmpIgnoreArr)

	if modules != nil && !isInArray {
		if res, err := helper.InArray(moduleName, modules); err == nil && res {

			userId, unitId := "", ""
			if headerBaseInfo.token != "" {
				brancaData, _ := helper.BrancaDecode(headerBaseInfo.token, moduleName)
				userId = brancaData.Sub
				unitId = brancaData.SubUnit
			}
			data := mq_dto.ApiLogDto{Uri: headerBaseInfo.url, Host: headerBaseInfo.host, Ip: headerBaseInfo.ip, Method: headerBaseInfo.method, UserId: userId, UnitId: unitId}

			// 加锁操作共享变量
			cacheMutex.Lock()
			cacheApiStatistics = append(cacheApiStatistics, data)
			needFlush := len(cacheApiStatistics) >= maxBatchSize
			cacheCopy := make([]interface{}, len(cacheApiStatistics))
			copy(cacheCopy, cacheApiStatistics)
			cacheMutex.Unlock()

			// 如果达到批量大小则立即发送
			if needFlush {
				go mqSendTask(cacheCopy)
			}
		}
	}
}

// api 统计-清空缓存
func flushCacheIfNeeded() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if len(cacheApiStatistics) > 0 {
		mqSendTask(cacheApiStatistics)
		cacheApiStatistics = make([]interface{}, 0) // 清空缓存
	}
}

// api 统计-mq发送任务
func mqSendTask(data interface{}) {
	if data == nil {
		return
	}
	dataNew := data.([]interface{})
	dataStr, err := json.Marshal(dataNew)
	if err != nil {
		global.Log.Error("mqSendTask() json.Marshal err:", err)
		return
	}
	args := []tasks.Arg{{Name: "actionSaveToDbData", Type: "string", Value: dataStr}}
	res, err := (&mq.MqClient{}).SendTask(constant.MQ_API_LOG_SAVE_TO_DB, args)
	if err != nil {
		global.Log.Error("mqSendTask() SendTask err:", err)
		global.Log.Error("mqSendTask() SendTask res:", res)
		return
	}
}

// 处理body签名和body加密，【解密后覆盖body】
func dealSignAndEncrypt(ctx *beecontext.Context) error {
	err := (&framework.AccessMiddleware{}).DealSignAndEncrypt(ctx)
	if err != nil {
		setResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}
	return err
}
