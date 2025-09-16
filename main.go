package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	v1 "github.com/haogeng/DiaryGo/api/v1"
	"github.com/haogeng/DiaryGo/internal/middleware"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"github.com/haogeng/DiaryGo/pkg/logger"
	"github.com/haogeng/DiaryGo/pkg/response"
	"github.com/joho/godotenv"
)

func main() {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到.env文件或加载失败")
	}

	// 初始化日志记录器
	if err := logger.Init(); err != nil {
		log.Fatalf("初始化日志记录器失败: %v", err)
	}
	defer logger.Close()

	// 初始化数据库
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
		logger.Error("数据库连接失败: %v", err)
	}

	// 自动迁移模型 - 先迁移 User 模型，再迁移 Diary 模型（因为有外键关系）
	if err := db.AutoMigrate(&models.User{}, &models.Diary{}, &models.DiaryImage{}, &models.DiaryVideo{}, &models.DiaryTag{}, &models.DiaryDPermission{}, &models.Tag{}, &models.DPermission{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
		logger.Error("数据库迁移失败: %v", err)
	}

	// 创建gin引擎 (不使用默认中间件，因为我们要自定义)
	r := gin.New()

	// 使用Recovery中间件从任何panic恢复
	r.Use(gin.Recovery())

	// 使用我们的计时中间件
	r.Use(middleware.Timer())

	// 基本路由
	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, "pong")
	})

	// 注册API v1版本的路由
	v1.RegisterRoutes(r)

	// 获取端口号，如果环境变量不存在则使用默认端口8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("服务器将在端口 %s 上启动", port)
	logger.Info("服务器将在端口 %s 上启动", port)

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
		logger.Error("启动服务器失败: %v", err)
	}
}
