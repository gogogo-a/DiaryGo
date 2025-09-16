package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	// 默认日志
	defaultLogger *Logger
)

// Logger 日志记录器
type Logger struct {
	logger *log.Logger
	file   *os.File
}

// Init 初始化默认日志记录器
func Init() error {
	// 创建logs目录
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return err
	}

	// 创建日志文件，按日期命名
	logFileName := fmt.Sprintf("logs/app-%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// 创建日志记录器
	logger := log.New(file, "", log.LstdFlags)
	defaultLogger = &Logger{
		logger: logger,
		file:   file,
	}

	return nil
}

// Close 关闭日志文件
func Close() {
	if defaultLogger != nil && defaultLogger.file != nil {
		defaultLogger.file.Close()
	}
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	if defaultLogger != nil {
		message := fmt.Sprintf(format, v...)
		defaultLogger.logger.Printf("[INFO] %s", message)
	}
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	if defaultLogger != nil {
		message := fmt.Sprintf(format, v...)
		defaultLogger.logger.Printf("[ERROR] %s", message)
	}
}

// RequestLog 记录请求日志
func RequestLog(method, path string, statusCode int, latency time.Duration) {
	if defaultLogger != nil {
		defaultLogger.logger.Printf("[REQUEST] %s %s %d %s", method, path, statusCode, latency)
	}
}
