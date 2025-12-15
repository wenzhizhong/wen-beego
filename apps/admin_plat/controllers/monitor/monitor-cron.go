package monitor

import (
	paltMonitorService "WenBeego/apps/admin_plat/services/monitor"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto/cron_dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/helper"
	// 	"encoding/json"
	// 	"strconv"
	// 	"your-project/models"
	// 	"your-project/utils"
)

type CronController struct {
	commonControllers.AdminBaseController
	cronService paltMonitorService.CronService
}

// 定时任务列表
func (c *CronController) Get() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	reqDataListDto, err1 := helper.GetReqDataListDto(&c.Controller)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}

	reqDto := &page_dto.MonitorCronListReqDto{}
	reqDto.BaseParamDto = baseParamDto
	reqDto.ReqDataListDto = reqDataListDto
	data, err := c.cronService.GetCronList(reqDto)

	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

func (c *CronController) Add() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[cron_dto.UnitCronDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.cronService.AddTask(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}

func (c *CronController) Edit() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[cron_dto.UnitCronDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.cronService.UpdateTask(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}

func (c *CronController) StartTask() {

	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[cron_dto.UnitCronDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.cronService.StartTask(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}
func (c *CronController) StopTask() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[cron_dto.UnitCronDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.cronService.StopTask(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}

// // @router /task/:id/stop [post]
// func (c *CronController) StopTask() {
// 	idStr := c.Ctx.Input.Param(":id")

// 	job := &models.SysJob{Id: int(id)}
// 	if models.GetJobById(job) != nil {
// 		c.Ctx.Output.SetStatus(404)
// 		c.Data["json"] = map[string]interface{}{"error": "Job not found"}
// 		c.ServeJSON()
// 		return
// 	}

// 	job.Status = 0
// 	if err := models.UpdateJob(job); err != nil {
// 		c.Ctx.Output.SetStatus(500)
// 		c.Data["json"] = map[string]interface{}{"error": "Failed to stop job"}
// 		c.ServeJSON()
// 		return
// 	}

// 	// 从调度器移除
// 	taskManager := utils.GetTaskManager()
// 	taskManager.RemoveTask(idStr)

// 	c.Data["json"] = map[string]interface{}{"message": "Task stopped successfully"}
// 	c.ServeJSON()
// }

// // @router /task/execute-job/:jobType [post]
// func (c *CronController) ExecuteJob() {
// 	jobType := c.Ctx.Input.Param(":jobType")

// 	// 验证请求来源（来自可信的调度器）
// 	if !c.isTrustedSource() {
// 		c.Ctx.Output.SetStatus(403)
// 		return
// 	}

// 	switch jobType {
// 	case "backup":
// 		c.backupDatabase()
// 	case "cleanup":
// 		c.cleanupLogs()
// 	}

// 	c.Data["json"] = map[string]string{"status": "success"}
// 	c.ServeJSON()
// }

// // 辅助方法：添加任务到调度器
// func (c *CronController) addTaskToScheduler(taskID, cronExpr, target string) {
// 	taskManager := utils.GetTaskManager()

// 	// 根据 target 确定要执行的具体任务
// 	var taskFunc func()
// 	switch target {
// 	case "backupDatabase":
// 		taskFunc = c.backupDatabase
// 	case "cleanLogs":
// 		taskFunc = c.cleanLogs
// 	case "sendReports":
// 		taskFunc = c.sendReports
// 	default:
// 		taskFunc = func() {
// 			// 默认任务处理逻辑
// 		}
// 	}

func (c *CronController) Del() {

}

// // @router /task/:id [post]
// func (c *CronController) RemoveTask() {
//     taskID := c.Ctx.Input.Param(":id")
//     taskManager := utils.GetTaskManager()
//     taskManager.RemoveTask(taskID)

//     c.Data["json"] = map[string]interface{}{"message": "Task removed successfully"}
//     c.ServeJSON()
// }

// // @router /tasks [get]
// func (c *CronController) ListTasks() {
//     taskManager := utils.GetTaskManager()
//     entries := taskManager.ListTasks()

//     var tasks []map[string]interface{}
//     for _, entry := range entries {
//         taskInfo := map[string]interface{}{
//             "id":       entry.ID,
//             "next":     entry.Next,
//             "prev":     entry.Prev,
//             "schedule": entry.Schedule,
//         }
//         tasks = append(tasks, taskInfo)
//     }

//     c.Data["json"] = tasks
//     c.ServeJSON()
// }
// 	// 添加到调度器
// 	taskManager.AddTask(cronExpr, taskFunc, taskID)
// }

// // 具体的任务实现
// func (c *CronController) backupDatabase() {
// 	// 实现数据库备份逻辑
// }

// func (c *CronController) cleanLogs() {
// 	// 实现日志清理逻辑
// }

// func (c *CronController) sendReports() {
// 	// 实现报告发送逻辑
// }
