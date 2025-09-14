package main

import (
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认的gin引擎
	r := gin.Default()

	// 定义路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
