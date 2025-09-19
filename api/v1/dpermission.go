package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/haogeng/DiaryGo/internal/repository"
	"github.com/haogeng/DiaryGo/pkg/response"
)

// DPermissionHandler 权限处理器
type DPermissionHandler struct {
	repo repository.DPermissionRepository
}

// NewDPermissionHandler 创建权限处理器
func NewDPermissionHandler() *DPermissionHandler {
	return &DPermissionHandler{
		repo: repository.NewDPermissionRepository(),
	}
}

// RegisterRoutes 注册权限相关路由
func (h *DPermissionHandler) RegisterRoutes(router *gin.RouterGroup) {
	permissions := router.Group("/permissions")
	{
		permissions.GET("", h.GetAllPermissions) // 获取所有权限
	}
}

// GetAllPermissions 获取所有权限
// @Summary 获取所有权限
// @Description 获取系统中所有可用的权限配置
// @Tags 权限
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]models.DPermission}
// @Failure 500 {object} response.Response
// @Router /api/v1/permissions [get]
func (h *DPermissionHandler) GetAllPermissions(c *gin.Context) {


	// 调用仓库方法获取所有权限
	permissions, err := h.repo.GetAll()
	if err != nil {
		response.ServerError(c, "获取权限列表失败")
		return
	}

	response.Success(c, permissions)
}
