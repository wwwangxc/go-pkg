package redis

import (
	"fmt"
	"log"
)

const (
	packageName = "go-pkg/redis"

	logStatusError = "[ERROR]"
	logStatusWarn  = "[WARN]"
)

func logErrorf(format string, args ...interface{}) {
	logf(logStatusError, format, args...)
}

func logWarn(format string, args ...interface{}) {
	logf(logStatusWarn, format, args...)
}

func logf(logStatus, format string, args ...interface{}) {
	log.Printf("%s %s %s", packageName, logStatus, fmt.Sprintf(format, args...))
}
