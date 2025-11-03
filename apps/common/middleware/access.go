package middleware

import (
	"WenBeego/apps/common/dto/mq_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/RichardKnop/machinery/v1/tasks"
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

type AccessMiddleware struct {
}

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
func (m *AccessMiddleware) LimitTimes() web.FilterFunc {
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
				http.Error(ctx.ResponseWriter, "sign error", http.StatusUnauthorized)
				return
			}
		}
	}
}

// api 请求前
func (m *AccessMiddleware) RouterBefore() web.FilterFunc {
	return func(ctx *beecontext.Context) {
		// url, host, shceme, method, token, ip := m.getBaseInfo(ctx)
		// timeStr := time.Now().Format("2006-01-02 15:04:05")
	}
}

// api 请求后
func (m *AccessMiddleware) RouterAfter(whiteApiList *[]string, authApiList *[]string) web.FilterFunc {
	return func(ctx *beecontext.Context) {
		m.statisticsApiLog(ctx, whiteApiList, authApiList)
	}
}

func (m *AccessMiddleware) getBaseInfo(ctx *beecontext.Context) (url, host, shceme, method, token, ip string) {
	url = ctx.Request.URL.Path
	host = ctx.Request.Host
	shceme = ctx.Request.URL.Scheme
	method = ctx.Request.Method
	token = ctx.Request.Header.Get("Authorization")
	ip = ctx.Request.RemoteAddr
	return
}

// api 统计
func (m *AccessMiddleware) statisticsApiLog(ctx *beecontext.Context, whiteApiList *[]string, authApiList *[]string) {
	modules, _ := global.GetConfigDiy("logToDbModules")
	moduleName := helper.ParseModuleFromRoute(ctx.Request.URL.Path)
	tmpIgnoreArr := helper.ArrayMerge(*whiteApiList, *authApiList)
	url, host, _, method, token, ip := m.getBaseInfo(ctx)
	isInArray, _ := helper.InArray(url, tmpIgnoreArr)

	if modules != nil && !isInArray {
		if res, err := helper.InArray(moduleName, modules); err == nil && res {

			userId, unitId := "", ""
			if token != "" {
				brancaData, _ := helper.BrancaDecode(token, moduleName)
				userId = brancaData.Sub
				unitId = brancaData.SubUnit
			}
			data := mq_dto.ApiLogDto{Uri: url, Host: host, Ip: ip, Method: method, UserId: userId, UnitId: unitId}

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

func flushCacheIfNeeded() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if len(cacheApiStatistics) > 0 {
		mqSendTask(cacheApiStatistics)
		cacheApiStatistics = make([]interface{}, 0) // 清空缓存
	}
}

// mq发送任务
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
	(&MqClient{}).SendTask("ApiLog.ActionSaveToDb", args)
}
