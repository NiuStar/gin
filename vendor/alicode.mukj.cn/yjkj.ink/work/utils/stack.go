package utils

import (
	"fmt"
	"github.com/json-iterator/go"
	"runtime"
)

type CodeStack struct {
	MethodName string
	FileName   string
	Line       int
}

func (stack *CodeStack) StackTrace() string {
	return fmt.Sprintf("%s [%s:%d]", stack.FileName, stack.MethodName, stack.Line)
}

func StackTrace() []*CodeStack {
	return call(0)
}

func StackTraceJson() string {
	stacks := call(0)
	var list []string
	for _, s := range stacks {
		list = append(list, s.StackTrace())
	}
	stacksJson, _ := jsoniter.MarshalToString(list)
	return stacksJson
}

func call(skip int) (results []*CodeStack) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		lists := call(skip + 1)
		results = append(results, lists...)
	}
	pcName := runtime.FuncForPC(pc).Name() //获取函数名
	stack := &CodeStack{MethodName: pcName, FileName: file, Line: line}
	results = append(results, stack)
	return results
}
