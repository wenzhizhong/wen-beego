package main

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/middleware/mq"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar"
	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/tasks"
	cmdCommon "WenBeego/cmd/common"
	"WenBeego/routers"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	_ "github.com/beego/beego/v2/core/config/yaml"
)

func main() {
	cmdCommon.RunBefore()
	cmdCommon.InitCommonSource("pathMq")
	cmdCommon.InitMqClient()

	taskDefs := routers.GetMqTasks()
	if len(taskDefs) == 0 {
		fmt.Println("No mq task……")
		return
	}

	mqServer, err := (&mq.MqServer{}).NewMq()
	if err != nil {
		global.Log.Error("MqServer.NewMq() error:", err)
		panic(err)
	}
	for _, task := range taskDefs {
		mqServer.RegisterTask(task.Name, func(args ...interface{}) error {
			return reflectCallback(task.CallBack, args)
		})
	}

	fmt.Println("Mq starting")
	DefaultQueue := mqServer.Server.GetConfig().DefaultQueue

	if len(os.Args) > 1 && os.Args[1] == "dlx" {
		dlxQueue := DefaultQueue + ".dlq"
		fmt.Printf("DLX mode, consuming from: %s\n", dlxQueue)
		if err := mqServer.Server.StartDLQConsuming(dlxQueue, func(sig *tasks.Signature) error {
			cb := findCallback(taskDefs, sig.Name)
			if cb == nil {
				global.Log.Error("DLQ task not registered:", sig.Name)
				return nil
			}
			args := make([]interface{}, len(sig.Args))
			for i, a := range sig.Args {
				args[i] = a.Value
			}
			err := reflectCallback(cb, args)
			if err != nil {
				global.Log.Error("DLQ task failed:", sig.UUID, sig.Name, err)
				if sig.RetryCount <= 0 {
					saveDLXFailedMsg(sig, err)
				}
				return err
			}
			global.Log.Info("DLQ task processed:", sig.UUID, sig.Name)
			return nil
		}); err != nil {
			global.Log.Error("StartDLQConsuming error:", err)
			panic(err)
		}
		fmt.Println("DLX consumer stopped")
		return
	}

	worker := mqServer.Server.NewCustomQueueWorker("worker_"+DefaultQueue, 10, DefaultQueue)
	if err := worker.Launch(); err != nil {
		global.Log.Error("MqServer.Launch() error:", err)
		panic(err)
	}
	fmt.Println("Mq stoped!")
}

func saveDLXFailedMsg(sig *tasks.Signature, taskErr error) {
	record := &models.QueueDlxFailedLog{
		TaskUUID:   sig.UUID,
		TaskName:   sig.Name,
		ErrorMsg:   taskErr.Error(),
		CreateTime: time.Now(),
	}
	if argsJSON, err := json.Marshal(sig.Args); err == nil {
		record.TaskArgs = string(argsJSON)
	}
	ar := &models_ar.QueueDlxFailedLogAR{}
	if err := ar.Insert(record); err != nil {
		global.Log.Error("saveDLXFailedMsg db error:", err)
	}
}

func findCallback(taskDefs []routers.MqTasks, name string) interface{} {
	for _, t := range taskDefs {
		if t.Name == name {
			return t.CallBack
		}
	}
	return nil
}

func reflectCallback(f interface{}, args interface{}) error {
	callbackValue := reflect.ValueOf(f)
	if callbackValue.Kind() != reflect.Func {
		return errors.New("callback is not a function")
	}
	funcObj := runtime.FuncForPC(callbackValue.Pointer())
	if funcObj == nil {
		return errors.New("cannot find the method")
	}

	callbackToString := funcObj.Name()
	methodType := callbackValue.Type()
	funcNameSli := strings.Split(callbackToString, ".")
	lFuncSli := len(funcNameSli)
	if lFuncSli == 0 {
		return errors.New("invalid method full name: " + callbackToString)
	}

	method := funcNameSli[lFuncSli-1]
	if len(method) == 0 {
		return errors.New("method name is empty")
	} else if method[0] > 96 || method[0] < 65 {
		return fmt.Errorf("%s is not a public method", method)
	}
	if numIn := methodType.NumIn(); numIn < 2 {
		return errors.New("invalid number of param in, more than one param is allowed")
	}

	controllerType := methodType.In(0)
	var receiver reflect.Value
	if controllerType.Kind() == reflect.Ptr {
		elemType := controllerType.Elem()
		receiver = reflect.New(elemType)
	} else {
		receiver = reflect.New(controllerType).Elem()
	}

	in := make([]reflect.Value, 1)
	in[0] = receiver
	if args != nil {
		argsSlice, ok := args.([]interface{})
		if !ok {
			return errors.New("args is not a slice of interface{}")
		}
		if len(argsSlice) != methodType.NumIn()-1 {
			return fmt.Errorf("%s param count mismatch, expected %d got %d",
				callbackToString, methodType.NumIn()-1, len(argsSlice))
		}
		for _, arg := range argsSlice {
			in = append(in, reflect.ValueOf(arg))
		}
	}

	res := callbackValue.Call(in)
	if len(res) > 0 {
		if err, ok := res[0].Interface().(error); ok {
			fmt.Println("mq exec error, callback:", err)
			return err
		}
	}
	return nil
}
