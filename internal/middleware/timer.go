package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haogeng/DiaryGo/pkg/logger"
)

// Timer 计时中间件，记录请求处理时间
func Timer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前记录时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 请求后计算耗时
		latency := time.Since(startTime)

		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		// 记录请求日志
		logger.RequestLog(method, path, statusCode, latency)
	}
}
