package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/internal/repository"
	"github.com/haogeng/DiaryGo/pkg/response"
)

// DiaryHandler 处理日记相关的请求
type DiaryHandler struct {
	repo repository.DiaryRepository
}

// NewDiaryHandler 创建日记处理器
func NewDiaryHandler() *DiaryHandler {
	return &DiaryHandler{
		repo: repository.NewDiaryRepository(),
	}
}

// RegisterRoutes 注册日记相关的路由
func (h *DiaryHandler) RegisterRoutes(router *gin.RouterGroup) {
	diaries := router.Group("/diaries")
	{
		diaries.GET("", h.List)          // 获取所有日记
		diaries.GET("/:id", h.Get)       // 获取单个日记
		diaries.POST("", h.Create)       // 创建日记
		diaries.PUT("/:id", h.Update)    // 更新日记
		diaries.DELETE("/:id", h.Delete) // 删除日记
	}
}

// List 获取所有日记
// @Summary 获取所有日记
// @Description 获取当前用户的所有日记
// @Tags diaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]models.Diary}
// @Router /api/v1/diaries [get]
func (h *DiaryHandler) List(c *gin.Context) {
	// 从上下文中获取用户ID（由Auth中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	diaries, err := h.repo.GetAll(userID.(uuid.UUID))
	if err != nil {
		response.ServerError(c, "获取日记失败")
		return
	}

	response.Success(c, diaries)
}

// Get 获取单个日记
// @Summary 获取单个日记
// @Description 根据ID获取单个日记
// @Tags diaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日记ID"
// @Success 200 {object} response.Response{data=models.Diary}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/diaries/{id} [get]
func (h *DiaryHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的ID格式")
		return
	}

	diary, err := h.repo.GetByID(id)
	if err != nil {
		response.NotFound(c, "日记不存在")
		return
	}

	response.Success(c, diary)
}

// Create 创建日记
// @Summary 创建日记
// @Description 创建新的日记
// @Tags diaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param diary body models.Diary true "日记信息"
// @Success 200 {object} response.Response{data=models.Diary}
// @Failure 400 {object} response.Response
// @Router /api/v1/diaries [post]
func (h *DiaryHandler) Create(c *gin.Context) {
	var diary models.Diary
	if err := c.ShouldBindJSON(&diary); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 从上下文中获取用户ID（由Auth中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 设置用户ID
	diary.UserId = userID.(uuid.UUID)

	if err := h.repo.Create(&diary); err != nil {
		response.ServerError(c, "创建日记失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", diary)
}

// Update 更新日记
// @Summary 更新日记
// @Description 根据ID更新日记
// @Tags diaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日记ID"
// @Param diary body models.Diary true "日记信息"
// @Success 200 {object} response.Response{data=models.Diary}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/diaries/{id} [put]
func (h *DiaryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的ID格式")
		return
	}

	// 从上下文中获取用户ID（由Auth中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	existingDiary, err := h.repo.GetByID(id)
	if err != nil {
		response.NotFound(c, "日记不存在")
		return
	}

	// 验证日记所有权
	if existingDiary.UserId != userID.(uuid.UUID) {
		response.Forbidden(c, "您没有权限修改此日记")
		return
	}

	var updatedDiary models.Diary
	if err := c.ShouldBindJSON(&updatedDiary); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 保留ID和用户ID
	updatedDiary.Id = existingDiary.Id
	updatedDiary.UserId = existingDiary.UserId

	if err := h.repo.Update(&updatedDiary); err != nil {
		response.ServerError(c, "更新日记失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", updatedDiary)
}

// Delete 删除日记
// @Summary 删除日记
// @Description 根据ID删除日记
// @Tags diaries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "日记ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diaries/{id} [delete]
func (h *DiaryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的ID格式")
		return
	}

	// 从上下文中获取用户ID（由Auth中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 获取日记以验证所有权
	existingDiary, err := h.repo.GetByID(id)
	if err != nil {
		response.NotFound(c, "日记不存在")
		return
	}

	// 验证日记所有权
	if existingDiary.UserId != userID.(uuid.UUID) {
		response.Forbidden(c, "您没有权限删除此日记")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		response.ServerError(c, "删除日记失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
