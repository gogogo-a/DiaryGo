package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/internal/repository"
	"github.com/haogeng/DiaryGo/pkg/response"
)

// TagHandler 标签处理器
type TagHandler struct {
	repo repository.TagRepository
}

// NewTagHandler 创建标签处理器
func NewTagHandler() *TagHandler {
	return &TagHandler{
		repo: repository.NewTagRepository(),
	}
}

// RegisterRoutes 注册标签相关路由
func (h *TagHandler) RegisterRoutes(router *gin.RouterGroup) {
	tags := router.Group("/tags")
	{
		tags.POST("", h.CreateTag)       // 创建标签
		tags.GET("", h.GetTags)          // 获取标签列表
		tags.GET("/:id", h.GetTag)       // 获取标签详情
		tags.PUT("/:id", h.UpdateTag)    // 更新标签
		tags.DELETE("/:id", h.DeleteTag) // 删除标签
		// tags.POST("/batch", h.BatchCreateTags) // 批量创建标签
	}
}

// TagRequest 标签请求参数
type TagRequest struct {
	TagName  string `json:"tag_name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Category string `json:"category" binding:"required"`
}

// BatchCreateTagsRequest 批量创建标签请求参数
type BatchCreateTagsRequest struct {
	Tags []TagRequest `json:"tags" binding:"required,min=1"`
}

// CreateTag 创建标签
// @Summary 创建标签
// @Description 创建新的标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param tag body TagRequest true "标签信息"
// @Success 201 {object} response.Response{data=models.Tag}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 创建标签对象
	tag := &models.Tag{
		TagName:  req.TagName,
		Type:     req.Type,
		Category: req.Category,
	}

	// 调用仓库方法创建标签
	if err := h.repo.Create(tag); err != nil {
		response.ServerError(c, "创建标签失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建标签成功", tag)
}

// GetTags 获取标签列表
// @Summary 获取标签列表
// @Description 获取标签列表，可按分类过滤
// @Tags 标签
// @Accept json
// @Produce json
// @Param category query string false "标签分类 (bill/diary)"
// @Success 200 {object} response.Response{data=[]models.Tag}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	category := c.Query("category")

	// 调用仓库方法获取标签列表
	tags, err := h.repo.GetAll(category)
	if err != nil {
		response.ServerError(c, "获取标签列表失败: "+err.Error())
		return
	}

	response.Success(c, tags)
}

// GetTag 获取标签详情
// @Summary 获取标签详情
// @Description 获取标签详情
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path string true "标签ID"
// @Success 200 {object} response.Response{data=models.Tag}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	// 解析标签ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的标签ID格式")
		return
	}

	// 调用仓库方法获取标签
	tag, err := h.repo.GetByID(id)
	if err != nil {
		if err.Error() == "标签不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.ServerError(c, "获取标签失败: "+err.Error())
		return
	}

	response.Success(c, tag)
}

// UpdateTag 更新标签
// @Summary 更新标签
// @Description 更新标签信息
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path string true "标签ID"
// @Param tag body TagRequest true "标签信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	// 解析标签ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的标签ID格式")
		return
	}

	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 获取原始标签信息
	originalTag, err := h.repo.GetByID(id)
	if err != nil {
		if err.Error() == "标签不存在" {
			response.NotFound(c, err.Error())
			return
		}
		response.ServerError(c, "获取标签失败: "+err.Error())
		return
	}

	// 更新标签信息
	originalTag.TagName = req.TagName
	originalTag.Type = req.Type
	originalTag.Category = req.Category

	// 调用仓库方法更新标签
	if err := h.repo.Update(originalTag); err != nil {
		response.ServerError(c, "更新标签失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新标签成功", nil)
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Description 删除标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param id path string true "标签ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	// 解析标签ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的标签ID格式")
		return
	}

	// 调用仓库方法删除标签
	if err := h.repo.Delete(id); err != nil {
		if err.Error() == "标签不存在" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "标签已被使用，无法删除" {
			response.ParamError(c, err.Error())
			return
		}
		response.ServerError(c, "删除标签失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除标签成功", nil)
}

// BatchCreateTags 批量创建标签
// @Summary 批量创建标签
// @Description 批量创建多个标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param tags body BatchCreateTagsRequest true "标签列表"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/tags/batch [post]
// func (h *TagHandler) BatchCreateTags(c *gin.Context) {
// 	var req BatchCreateTagsRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		response.ParamError(c, err.Error())
// 		return
// 	}

// 	// 转换请求数据到标签对象
// 	var tags []*models.Tag
// 	for _, t := range req.Tags {
// 		tag := &models.Tag{
// 			TagName:  t.TagName,
// 			Type:     t.Type,
// 			Category: t.Category,
// 		}
// 		tags = append(tags, tag)
// 	}

// 	// 调用仓库方法批量创建标签
// 	if err := h.repo.BatchCreate(tags); err != nil {
// 		response.ServerError(c, "批量创建标签失败: "+err.Error())
// 		return
// 	}

// 	response.SuccessWithMessage(c, "批量创建标签成功", nil)
// }
