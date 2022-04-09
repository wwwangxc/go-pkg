package config

import (
	"fmt"
	"log"
)

const (
	packageName = "go-pkg/config"

	logStatusError = "[ERROR]"
	logStatusWarn  = "[WARN]"
	logStatusInfo  = "[INFO]"
)

func logErrorf(format string, args ...interface{}) {
	logf(logStatusError, format, args...)
}

func logWarn(format string, args ...interface{}) {
	logf(logStatusWarn, format, args...)
}

func logInfo(format string, args ...interface{}) {
	logf(logStatusInfo, format, args...)
}

func logf(logStatus, format string, args ...interface{}) {
	log.Printf("%s %s %s", packageName, logStatus, fmt.Sprintf(format, args...))
}
