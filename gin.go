package gin

import (
	"github.com/gin-gonic/gin"
	"io"
)

const EnvGinMode = "GIN_MODE"

const (
	// DebugMode indicates gin mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates gin mode is release.
	ReleaseMode = "release"
	// TestMode indicates gin mode is test.
	TestMode = "test"
)

func ForceConsoleColor() {
	gin.ForceConsoleColor()
}
func IsDebugging() bool {
	return gin.IsDebugging()
}
func SetMode(mode string) {
	gin.SetMode(mode)
}
func DefaultWriter() io.Writer {
	return gin.DefaultWriter
}
func DebugPrintRouteFunc(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	gin.DebugPrintRouteFunc(httpMethod, absolutePath, handlerName, nuHandlers)
}
