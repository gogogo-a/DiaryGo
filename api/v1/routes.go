package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/haogeng/DiaryGo/internal/middleware"
)

// RegisterRoutes 注册API v1版本的所有路由
func RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	// 注册公开路由（不需要认证）
	// 用户登录相关路由
	userLoginHandler := NewUserLoginHandler()
	userLoginHandler.RegisterRoutes(v1)

	// 权限配置相关路由（不需要认证，便于客户端初始化）
	permissionHandler := NewDPermissionHandler()
	permissionHandler.RegisterRoutes(v1)

	// 注册需要认证的路由
	// 使用Auth中间件保护以下路由
	protected := v1.Group("")
	protected.Use(middleware.Auth())
	{
		// 日记相关路由（需要认证）
		diaryHandler := NewDiaryHandler()
		diaryHandler.RegisterRoutes(protected)

		// 日记扩展相关路由（图片、视频管理）
		diaryExtendsHandler := NewDiaryExtendsHandler()
		diaryExtendsHandler.RegisterRoutes(protected)

		// 账本相关路由（需要认证）
		accountBookHandler := NewAccountBookHandler()
		accountBookHandler.RegisterRoutes(protected)

		// 账本用户权限相关路由（需要认证）
		accountBookUserHandler := NewAccountBookUserHandler()
		accountBookUserHandler.RegisterRoutes(protected)

		// 用户管理相关路由（需要认证）
		userHandler := NewUserHandler()
		userHandler.RegisterRoutes(protected)

		// 账单相关路由（需要认证）
		billHandler := NewBillHandler()
		billHandler.RegisterRoutes(protected)

		// 标签相关路由（需要认证）
		tagHandler := NewTagHandler()
		tagHandler.RegisterRoutes(protected)
	}
}
