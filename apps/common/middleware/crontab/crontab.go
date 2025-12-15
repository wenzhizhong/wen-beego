package crontab

import (
	"WenBeego/apps/common/global"
	"fmt"
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
			cron: cron.New(cron.WithSeconds()), // 支持秒级精度
		}
	})
	fmt.Println("init crontab manager")
	return cronManager
}
func (cm *CronManager) AddSafeTask(spec string, unsafeCmd func(), taskID string) error {
	safeCmd := func() {
		defer func() {
			if r := recover(); r != nil {
				// 记录日志或做其他处理
				fmt.Printf("Task %s panicked: %v\n", taskID, r)
				global.Log.Error("Task %s panicked: %v\n", taskID, r)
			}
		}()
		unsafeCmd()
	}

	return cm.AddTask(spec, safeCmd, taskID)
}

// 添加任务
func (cm *CronManager) AddTask(spec string, cmd func(), taskID string) error {
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
func (cm *CronManager) LoadJobsFromDB() error {
	// jobs, err := models.GetAllActiveJobs()
	// if err != nil {
	//     return err
	// }

	// for _, job := range jobs {
	//     // 根据任务类型确定执行函数
	//     var taskFunc func()
	//     switch job.InvokeTarget {
	//     case "backupDatabase":
	//         taskFunc = func() { /* 执行备份 */ }
	//     case "cleanLogs":
	//         taskFunc = func() { /* 清理日志 */ }
	//     // ... 其他任务类型
	//     default:
	//         continue // 跳过未知任务类型
	//     }

	//     // 添加到调度器
	//     cm.AddTask(job.CronExpression, taskFunc, strconv.Itoa(job.Id))
	// }

	return nil
}
