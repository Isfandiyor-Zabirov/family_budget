package logger

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
)

const logFileName = "internal/logger/logs.log"

var customLogger = &lumberjack.Logger{
	Filename:   logFileName,
	MaxSize:    20,   // megabytes
	MaxBackups: 90,   //quantity of files
	MaxAge:     30,   //days
	Compress:   true, // disabled by default
}

// InitLogger - запускает логирование в файле logs.log
func InitLogger() {
	log.SetOutput(customLogger)
	log.Println("Logging to file: " + logFileName)
}

func GetLogger() *lumberjack.Logger {
	return customLogger
}

func FormatLogger(router *gin.Engine) {
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[GIN] %s - [%s] \"%s %s %s %d %s \"%s\" %s \n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
}
