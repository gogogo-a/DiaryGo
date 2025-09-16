package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/repository"
	"github.com/haogeng/DiaryGo/pkg/response"
)

// UserHandler 处理用户相关的请求
// 提供用户管理功能，包括查询用户列表（支持分页和搜索）
// URL:
// - GET /api/v1/users - 获取用户列表，支持分页和搜索
// - GET /api/v1/users/:id - 获取特定用户详情
type UserHandler struct {
	repo repository.UserRepository
}

// NewUserHandler 创建用户处理器
// 初始化处理器并注入所需的仓库依赖
func NewUserHandler() *UserHandler {
	return &UserHandler{
		repo: repository.NewUserRepository(),
	}
}

// RegisterRoutes 注册用户相关的路由
// 注册以下路由:
// - GET /users - 获取用户列表，支持分页和搜索
// - GET /users/:id - 获取特定用户详情
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("", h.ListUsers)
		users.GET("/:id", h.GetUser)
	}
}

// UserQuery 用户查询参数
type UserQuery struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword" binding:"omitempty"`
}

// ListUsers 获取用户列表
// URL: GET /api/v1/users
// 功能: 获取用户列表，支持分页和关键词搜索
// 权限: 需要用户登录
// 查询参数:
//   - page: 页码，默认为1
//   - page_size: 每页数量，默认为10，最大100
//   - keyword: 搜索关键词，可选
//
// 返回:
//   - 成功: 200 状态码，包含用户列表和分页信息
//   - 失败: 400 (参数错误), 401 (未授权), 500 (服务器错误)
//
// @Summary 获取用户列表
// @Description 获取用户列表，支持分页和关键词搜索
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码，默认为1"
// @Param page_size query int false "每页数量，默认为10，最大100"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=response.PagedData{list=[]models.User}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 获取当前用户ID
	_, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析查询参数
	var query UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	} else if query.PageSize > 100 {
		query.PageSize = 100
	}

	// 获取用户列表
	users, total, err := h.repo.GetUsers(query.Page, query.PageSize, query.Keyword)
	if err != nil {
		response.ServerError(c, "获取用户列表失败")
		return
	}

	// 构造分页响应
	pagedData := response.PagedData{
		List:     users,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	response.Success(c, pagedData)
}

// GetUser 获取特定用户详情
// URL: GET /api/v1/users/:id
// 功能: 获取特定用户的详细信息
// 权限: 需要用户登录
// 参数:
//   - id: 用户ID (UUID格式，通过URL路径参数传递)
//
// 返回:
//   - 成功: 200 状态码，包含用户详情的JSON数据
//   - 失败: 400 (参数错误), 401 (未授权), 404 (用户不存在)
//
// @Summary 获取用户详情
// @Description 获取特定用户的详细信息
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "用户ID"
// @Success 200 {object} response.Response{data=models.User}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	// 获取当前用户ID
	_, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析用户ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的用户ID格式")
		return
	}

	// 获取用户详情
	user, err := h.repo.GetByID(id)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}
