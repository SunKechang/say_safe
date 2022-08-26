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
	prefix := "%s:%d "
	_, file, line, _ := runtime.Caller(1)
	str := fmt.Sprintf(prefix+format, file, line, texts)
	fmt.Fprintf(gin.DefaultWriter, str)
}
