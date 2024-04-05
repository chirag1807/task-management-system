package utils

import (
	"fmt"
	"runtime"
)

func CreateErrorMessage() string {
	pc, file, line, _ := runtime.Caller(1)
	functionName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("Error in file %s, function %s, line %d", file, functionName, line)
}
