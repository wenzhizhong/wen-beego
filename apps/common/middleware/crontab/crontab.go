package crontab

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar"
	"WenBeego/routers/crontab_task"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

type CronManager struct {
	cron  *cron.Cron
	tasks sync.Map // 存储任务ID映射
	// once  sync.Once
}

var cronManager *CronManager
var once sync.Once

func GetCronManager() *CronManager {
	once.Do(func() {
		cronManager = &CronManager{
			cron: cron.New(
				cron.WithSeconds(), // 支持秒级精度
				cron.WithChain(cron.Recover(cron.DefaultLogger)), // 清除默认链，禁用 cron 库的默认 recover 机制
			),
		}
	})
	fmt.Println("init crontab manager...")
	return cronManager
}
func (cm *CronManager) AddSafeTask(spec string, unsafeCmd func(), taskID string) error {
	safeCmd := func() {
		defer func() {
			result := true
			errMsg := ""
			if r := recover(); r != nil {
				result = false
				// 记录日志或做其他处理
				traceStr := cm.GetTraceStr()
				errMsg = fmt.Sprintf("Crontab Task %s panicked: %v\ntrace:\n%s\n", taskID, r, traceStr)
				global.Log.Error(errMsg)
			}

			err := (&models_ar.PlatCronLogAr{}).Insert(taskID, result, errMsg)
			if err != nil {
				global.Log.Error("insert crontab log error: %v", err)
			}
		}()
		unsafeCmd()
	}

	return cm.addTask(spec, safeCmd, taskID)
}

// 添加任务
func (cm *CronManager) addTask(spec string, cmd func(), taskID string) error {
	if _, ok := cm.tasks.Load(taskID); ok {
		return nil
	}
	entryID, err := cm.cron.AddFunc(spec, cmd)
	if err != nil {
		return err
	}
	cm.tasks.Store(taskID, entryID)
	return nil
}

// 删除任务
func (cm *CronManager) RemoveTask(taskID string) {
	if entryID, ok := cm.tasks.Load(taskID); ok {
		cm.cron.Remove(entryID.(cron.EntryID))
		cm.tasks.Delete(taskID)
	}
}

// 启动调度器
func (cm *CronManager) Start() {
	cm.LoadCrontabsFromDB()
	cm.cron.Start()
}

// 停止调度器
func (cm *CronManager) Stop() {
	cm.cron.Stop()
}

// 获取所有任务
func (cm *CronManager) ListTasks() []cron.Entry {
	var entries []cron.Entry
	// for _, entry := range cm.cron.Entries() {
	// 	entries = append(entries, entry)
	// }
	entries = append(entries, cm.cron.Entries()...)
	return entries
}

// 加载数据库中的启用任务
func (cm *CronManager) LoadCrontabsFromDB() error {
	crontabDbList, err := (&models_ar.PlatCronAr{}).RunProjectGetCronList()
	if err != nil || len(crontabDbList) == 0 {
		return err
	}

	crontabDbListMap := make(map[string]models.PlatCron)
	for _, item := range crontabDbList {
		crontabDbListMap[item.NameEn] = item
	}
	list := crontab_task.GetCronTasks()
	for _, item := range list {
		exist, ok := crontabDbListMap[item.Name]
		if !(ok && item.Name == crontabDbListMap[item.Name].NameEn) {
			continue
		}
		err = cm.AddSafeTask(exist.CronExpr, item.CallBack, item.Name)
		if err != nil {
			return err
		}
	}
	global.Log.Info("loaded crontab from database!")
	return nil
}

func (cm *CronManager) GetTraceStr() string {
	pcs := make([]uintptr, 100)
	n := runtime.Callers(0, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	index := 0
	traceStr := ""
	tmoRootPath := strings.ReplaceAll(global.RootPath, "\\", "/")
	for {
		index++
		frame, more := frames.Next()
		if !more {
			break
		}
		if index <= 4 {
			continue
		} else if index >= 100 {
			break
		}

		if strings.HasPrefix(frame.File, tmoRootPath) {
			traceStr += fmt.Sprintf("  %s:%d %s\n", frame.File, frame.Line, frame.Function)
		}
	}
	return traceStr
}
