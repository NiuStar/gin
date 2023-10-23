package utils

import (
	"runtime"
	"reflect"
	"strings"
	"fmt"
)

func GetFunctionName(i interface{}) string {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	fmt.Println("fn:",fn)

	name := fn[strings.LastIndex(fn,"."):strings.LastIndex(fn,"-")]


	return name
}
