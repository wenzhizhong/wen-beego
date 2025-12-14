package main

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/middleware"
	cmdCommon "WenBeego/cmd/common"
	"WenBeego/routers"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	_ "github.com/beego/beego/v2/core/config/yaml"
	// beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// 注册自己资源服务
	cmdCommon.RunBefore()
	cmdCommon.InitCommonSource("pathMq")

	tasks := routers.GetMqTasks()
	if len(tasks) == 0 {
		fmt.Println("No mq task……")
		return
	}

	// 启动mq服务
	mqServer, err := (&middleware.MqServer{}).NewMq()
	if err != nil {
		global.Log.Error("MqServer.NewMq() error:", err)
		panic(err)
	}
	// 注册任务
	for _, task := range tasks {
		mqServer.RegisterTask(task.Name, func(args ...interface{}) error {
			return reflectCallback(task.CallBack, args)
		})
	}

	fmt.Println("Mq starting")
	// 启动Worker"服务"
	// worker := mqServer.Server.NewWorker("worker_1", 10)
	DefaultQueue := mqServer.Server.GetConfig().DefaultQueue
	worker := mqServer.Server.NewCustomQueueWorker("worker_"+DefaultQueue, 10, DefaultQueue)
	if err := worker.Launch(); err != nil {
		global.Log.Error("MqServer.Launch() error:", err)
		panic(err)
	}
	fmt.Println("Mq started！")
}

/**
 * 反射获取回调函数
 * @param f 订阅回调函数
 * @param args 参数
 * @return callbackValue 获取反射回调函数
 * @return paramTypes 参数类型
 * @return callbackToString 方法名（包含路径）
 */
func reflectCallback(f interface{}, args interface{}) error {
	callbackValue := reflect.ValueOf(f)
	// 检查是否是函数类型
	if callbackValue.Kind() != reflect.Func {
		err := errors.New("callback is not a function")
		return err
	}
	funcObj := runtime.FuncForPC(callbackValue.Pointer())
	if funcObj == nil {
		err := errors.New("cannot find the method")
		return err
	}

	// 获取方法的类型信息
	callbackToString := funcObj.Name()
	methodType := callbackValue.Type()
	funcNameSli := strings.Split(callbackToString, ".")
	lFuncSli := len(funcNameSli)
	if lFuncSli == 0 {
		err := errors.New("invalid method full name: " + callbackToString)
		return err
	}

	method := funcNameSli[lFuncSli-1]
	if len(method) == 0 {
		err := errors.New("method name is empty")
		return err
	} else if method[0] > 96 || method[0] < 65 {
		err := fmt.Errorf("%s is not a public method", method)
		return err
	}
	// check only one param which is the method receiver
	if numIn := methodType.NumIn(); numIn < 2 {
		err := errors.New("invalid number of param in, more than one param is allowed")
		return err
	}

	// 方法的第一个参数是接收者类型（例如 *controllers.ActionXxx）
	controllerType := methodType.In(0)
	var receiver reflect.Value
	// 动态创建接收者实例
	if controllerType.Kind() == reflect.Ptr {
		// 如果接收者是指针类型（如 *controllers.ActionXxx），创建新实例并获取指针
		elemType := controllerType.Elem() // 获取指针指向的元素类型，即 controllers.ActionXxx
		receiver = reflect.New(elemType)  // 创建该元素类型的新实例，并返回其指针的reflect.Value
	} else {
		// 如果接收者是值类型，创建值的实例
		receiver = reflect.New(controllerType).Elem()
	}

	// 准备调用参数：第一个是接收者，后面是方法本身的参数
	in := make([]reflect.Value, 1)
	in[0] = receiver
	if args != nil {
		argsSlice, ok := args.([]interface{})
		if !ok {
			return errors.New("args is not a slice of interface{}")
		}

		if len(argsSlice) != methodType.NumIn()-1 {
			return fmt.Errorf("%s 参数数量不匹配，期望 %d，实际 %d", callbackToString, methodType.NumIn()-1, len(argsSlice))
		}

		for _, arg := range argsSlice {
			in = append(in, reflect.ValueOf(arg))
		}
	}

	res := callbackValue.Call(in)
	if len(res) > 0 {
		if err, ok := res[0].Interface().(error); ok {
			return err
		}
	}
	return nil
}
