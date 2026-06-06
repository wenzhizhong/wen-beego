package base_helper

import (
	"fmt"
	"runtime"
	"strings"
)

func GetTraceStr() string {
	pcs := make([]uintptr, 100)
	n := runtime.Callers(0, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	index := 0
	traceStr := ""
	lastFuncName := ""
	// tmoRootPath := strings.ReplaceAll(global.RootPath, "\\", "/")
	for {
		index++
		frame, more := frames.Next()
		if !more {
			break
		}
		if index < 4 {
			lastFuncName = frame.Function
			continue
		} else if index >= 100 {
			break
		}

		// if strings.HasPrefix(frame.File, tmoRootPath) {
		// curLineFuncName := formatFuncName(lastFuncName)
		// traceStr += fmt.Sprintf("  %s:%d %s\n", frame.File, frame.Line, curLineFuncName)
		// }
		curLineFuncName := formatFuncName(lastFuncName)
		traceStr += fmt.Sprintf("  %s:%d %s\n", frame.File, frame.Line, curLineFuncName)
		lastFuncName = frame.Function
	}
	// traceStr = Ternary(traceStr != "", strings.ReplaceAll(traceStr, tmoRootPath+"/", ""), traceStr)
	return traceStr
}
func formatFuncName(fullName string) string {
	if fullName == "" {
		return ""
	}

	parts := strings.Split(fullName, "/")
	partsLen := len(parts)
	lastPart := ""
	if partsLen > 0 {
		lastPart = parts[partsLen-1] + "()"
	}
	return lastPart
}
