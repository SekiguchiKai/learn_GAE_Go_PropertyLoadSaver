package util

import (
	"strconv"
	"context"
	"runtime"
	"google.golang.org/appengine/log"
	"github.com/gin-gonic/gin"
)

func RespondAndLog(c *gin.Context, code int, format string, values ...interface{}) {
	if code >= 500 {
		ErrorLog(c, format, values...)
	} else if code >= 400 {
		InfoLog(c, format, values...)
	}
	c.String(code, format, values...)
}


func CriticalLog(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "nofile"
		line = -1
	}
	log.Criticalf(c, file+":"+strconv.Itoa(line)+":"+format, args...)
}

func DebugLog(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "nofile"
		line = -1
	}
	log.Debugf(c, file+":"+strconv.Itoa(line)+":"+format, args...)
}

func ErrorLog(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "nofile"
		line = -1
	}
	log.Errorf(c, file+":"+strconv.Itoa(line)+":"+format, args...)
}

func InfoLog(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "nofile"
		line = -1
	}
	log.Infof(c, file+":"+strconv.Itoa(line)+":"+format, args...)
}

func WarningLog(c context.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "nofile"
		line = -1
	}
	log.Warningf(c, file+":"+strconv.Itoa(line)+":"+format, args...)
}