package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

func Log(text string) {
	prefix := "%s:%d "
	_, file, line, _ := runtime.Caller(1)
	str := fmt.Sprintf(prefix+text, file, line)
	fmt.Fprintf(gin.DefaultWriter, str)
}

func Logger(format string, texts ...interface{}) {
	fmt.Fprintf(gin.DefaultWriter, fmt.Sprint(format, texts))
}
