package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/internal/repository"
	"github.com/haogeng/DiaryGo/pkg/response"
)

// DiaryHandler 日记处理器
type DiaryHandler struct {
	repo repository.DiaryRepository
}

// NewDiaryHandler 创建日记处理器
func NewDiaryHandler() *DiaryHandler {
	return &DiaryHandler{
		repo: repository.NewDiaryRepository(),
	}
}

// RegisterRoutes 注册日记相关路由
func (h *DiaryHandler) RegisterRoutes(router *gin.RouterGroup) {
	diary := router.Group("/diary")
	{
		diary.POST("", h.CreateDiary)            // 创建日记
		diary.GET("", h.ListDiaries)             // 获取日记列表
		diary.PUT("/:id", h.UpdateDiary)         // 更新日记
		diary.DELETE("/:id", h.DeleteDiary)      // 删除日记
		diary.GET("/:id", h.GetDiary)            // 获取日记
		diary.POST("/:id/share", h.ShareDiary)   // 分享给朋友
		diary.POST("/:id/like", h.LikeDiary)     // 点赞日记
		diary.DELETE("/:id/like", h.UnlikeDiary) // 取消点赞
		diary.GET("/:id/like", h.CheckLikeDiary) // 检查是否点赞
	}
}

// 日记请求参数
type DiaryRequest struct {
	Title        string   `json:"title" binding:"required"`
	Content      string   `json:"content" binding:"required"`
	Address      string   `json:"address"`
	PermissionId string   `json:"permission_id" binding:"required"`
	TagIds       []string `json:"tag_ids" binding:"required"`
	ImageUrls    []string `json:"image_urls"`
	VideoUrls    []string `json:"video_urls"`
}

// 日记查询参数
type DiaryQuery struct {
	Page         int      `form:"page" binding:"omitempty,min=1"`
	PageSize     int      `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword      string   `form:"keyword"`
	TagIds       []string `form:"tag_ids"`
	PermissionId string   `form:"permission_id"`
}

// 分享日记请求
type ShareDiaryRequest struct {
	UserId string `json:"user_id" binding:"required"`
}

// CreateDiary 创建日记
// @Summary 创建日记
// @Description 创建新的日记，包括标题、内容、地址、权限、标签、图片和视频
// @Tags 日记
// @Accept json
// @Produce json
// @Param diary body DiaryRequest true "日记信息"
// @Success 201 {object} response.Response{data=models.Diary}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary [post]
func (h *DiaryHandler) CreateDiary(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析请求
	var req DiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 转换权限ID
	permissionId, err := uuid.Parse(req.PermissionId)
	if err != nil {
		response.ParamError(c, "无效的权限ID格式")
		return
	}

	// 转换标签ID
	var tagIds []uuid.UUID
	for _, tagIdStr := range req.TagIds {
		tagId, err := uuid.Parse(tagIdStr)
		if err != nil {
			response.ParamError(c, "无效的标签ID: "+tagIdStr)
			return
		}
		tagIds = append(tagIds, tagId)
	}

	// 创建日记对象
	diary := &models.Diary{
		Title:   req.Title,
		Content: req.Content,
		Address: req.Address,
	}

	// 调用仓库方法创建日记
	createdDiary, err := h.repo.CreateDiary(diary, userID.(uuid.UUID), permissionId, tagIds, req.ImageUrls, req.VideoUrls)
	if err != nil {
		response.ServerError(c, "创建日记失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建日记成功", createdDiary)
}

// ListDiaries 获取日记列表
// @Summary 获取日记列表
// @Description 获取日记列表，支持分页和多条件搜索
// @Tags 日记
// @Accept json
// @Produce json
// @Param page query int false "页码，默认为1"
// @Param page_size query int false "每页数量，默认为10，最大100"
// @Param keyword query string false "搜索关键词"
// @Param tag_ids query []string false "标签ID列表"
// @Param permission_id query string false "权限ID"
// @Success 200 {object} response.Response{data=[]models.Diary}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary [get]
func (h *DiaryHandler) ListDiaries(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析查询参数
	var query DiaryQuery
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

	// 转换标签ID
	var tagIds []uuid.UUID
	for _, tagIdStr := range query.TagIds {
		tagId, err := uuid.Parse(tagIdStr)
		if err != nil {
			response.ParamError(c, "无效的标签ID: "+tagIdStr)
			return
		}
		tagIds = append(tagIds, tagId)
	}

	// 转换权限ID
	var permissionId *uuid.UUID
	if query.PermissionId != "" {
		pId, err := uuid.Parse(query.PermissionId)
		if err != nil {
			response.ParamError(c, "无效的权限ID格式")
			return
		}
		permissionId = &pId
	}

	// 调用仓库方法获取日记列表
	diaries, total, err := h.repo.GetDiaries(query.Page, query.PageSize, query.Keyword, tagIds, permissionId, userID.(uuid.UUID))
	if err != nil {
		response.ServerError(c, "获取日记列表失败: "+err.Error())
		return
	}

	// 构造分页响应
	pagedData := response.PagedData{
		List:     diaries,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	response.Success(c, pagedData)
}

// UpdateDiary 更新日记
// @Summary 更新日记
// @Description 更新日记信息，包括标题、内容、地址、权限、标签、图片和视频
// @Tags 日记
// @Accept json
// @Produce json
// @Param id path string true "日记ID"
// @Param diary body DiaryRequest true "日记信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary/{id} [put]
func (h *DiaryHandler) UpdateDiary(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析日记ID
	idStr := c.Param("id")
	diaryId, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 解析请求
	var req DiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 构造更新数据
	updateData := make(map[string]interface{})
	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.Content != "" {
		updateData["content"] = req.Content
	}
	if req.Address != "" {
		updateData["address"] = req.Address
	}

	// 转换权限ID
	var permissionId *uuid.UUID
	if req.PermissionId != "" {
		pId, err := uuid.Parse(req.PermissionId)
		if err != nil {
			response.ParamError(c, "无效的权限ID格式")
			return
		}
		permissionId = &pId
	}

	// 转换标签ID
	var tagIds []uuid.UUID
	for _, tagIdStr := range req.TagIds {
		tagId, err := uuid.Parse(tagIdStr)
		if err != nil {
			response.ParamError(c, "无效的标签ID: "+tagIdStr)
			return
		}
		tagIds = append(tagIds, tagId)
	}

	// 调用仓库方法更新日记
	err = h.repo.UpdateDiary(diaryId, userID.(uuid.UUID), updateData, permissionId, tagIds, req.ImageUrls, req.VideoUrls)
	if err != nil {
		if err.Error() == "无权更新该日记" {
			response.Forbidden(c, err.Error())
			return
		}
		if err.Error() == "record not found" {
			response.NotFound(c, "日记不存在")
			return
		}
		response.ServerError(c, "更新日记失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新日记成功", nil)
}

// DeleteDiary 删除日记
// @Summary 删除日记
// @Description 删除日记（只有创建者可以删除）
// @Tags 日记
// @Accept json
// @Produce json
// @Param id path string true "日记ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary/{id} [delete]
func (h *DiaryHandler) DeleteDiary(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析日记ID
	idStr := c.Param("id")
	diaryId, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 调用仓库方法删除日记
	err = h.repo.DeleteDiary(diaryId, userID.(uuid.UUID))
	if err != nil {
		if err.Error() == "只有创建者可以删除日记" {
			response.Forbidden(c, err.Error())
			return
		}
		if err.Error() == "record not found" {
			response.NotFound(c, "日记不存在")
			return
		}
		response.ServerError(c, "删除日记失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除日记成功", nil)
}

// GetDiary 获取日记详情
// @Summary 获取日记详情
// @Description 获取日记详情，包括日记内容、标签、权限、图片和视频
// @Tags 日记
// @Accept json
// @Produce json
// @Param id path string true "日记ID"
// @Success 200 {object} response.Response{data=models.Diary}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary/{id} [get]
func (h *DiaryHandler) GetDiary(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析日记ID
	idStr := c.Param("id")
	diaryId, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 调用仓库方法获取日记详情
	diary, tags, permission, images, videos, err := h.repo.GetDiaryWithDetails(diaryId, userID.(uuid.UUID))
	if err != nil {
		if err.Error() == "无权访问该日记" {
			response.Forbidden(c, err.Error())
			return
		}
		if err.Error() == "record not found" {
			response.NotFound(c, "日记不存在")
			return
		}
		response.ServerError(c, "获取日记详情失败: "+err.Error())
		return
	}

	// 检查当前用户是否已点赞
	isLiked, err := h.repo.CheckUserLike(diaryId, userID.(uuid.UUID))
	if err != nil {
		response.ServerError(c, "检查点赞状态失败: "+err.Error())
		return
	}

	// 构造响应
	result := map[string]interface{}{
		"diary":      diary,
		"tags":       tags,
		"permission": permission,
		"images":     images,
		"videos":     videos,
		"is_liked":   isLiked,
	}

	response.Success(c, result)
}

// ShareDiary 分享日记
// @Summary 分享日记给其他用户
// @Description 分享日记给其他用户
// @Tags 日记
// @Accept json
// @Produce json
// @Param id path string true "日记ID"
// @Param request body ShareDiaryRequest true "分享请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary/{id}/share [post]
func (h *DiaryHandler) ShareDiary(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析日记ID
	idStr := c.Param("id")
	diaryId, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 解析请求
	var req ShareDiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 解析分享目标用户ID
	shareUserId, err := uuid.Parse(req.UserId)
	if err != nil {
		response.ParamError(c, "无效的用户ID格式")
		return
	}

	// 调用仓库方法分享日记
	err = h.repo.ShareDiary(diaryId, shareUserId, userID.(uuid.UUID))
	if err != nil {
		if err.Error() == "您没有权限分享此日记" {
			response.Forbidden(c, err.Error())
			return
		}
		if err.Error() == "该用户已经有此日记的权限" {
			response.ParamError(c, err.Error())
			return
		}
		response.ServerError(c, "分享日记失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "分享日记成功", nil)
}

// LikeDiary 点赞日记
// @Summary 给日记点赞
// @Description 给日记点赞，增加点赞数
// @Tags 日记
// @Accept json
// @Produce json
// @Param id path string true "日记ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary/{id}/like [post]
func (h *DiaryHandler) LikeDiary(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析日记ID
	idStr := c.Param("id")
	diaryId, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 调用仓库方法增加点赞
	err = h.repo.AddLike(diaryId, userID.(uuid.UUID))
	if err != nil {
		if err.Error() == "日记不存在" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "您已经点赞过该日记" {
			response.ParamError(c, err.Error())
			return
		}
		response.ServerError(c, "点赞失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "点赞成功", nil)
}

// UnlikeDiary 取消点赞
// @Summary 取消日记点赞
// @Description 取消日记点赞，减少点赞数
// @Tags 日记
// @Accept json
// @Produce json
// @Param id path string true "日记ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary/{id}/like [delete]
func (h *DiaryHandler) UnlikeDiary(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析日记ID
	idStr := c.Param("id")
	diaryId, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 调用仓库方法取消点赞
	err = h.repo.RemoveLike(diaryId, userID.(uuid.UUID))
	if err != nil {
		if err.Error() == "日记不存在" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "您尚未点赞该日记" {
			response.ParamError(c, err.Error())
			return
		}
		response.ServerError(c, "取消点赞失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "取消点赞成功", nil)
}

// CheckLikeDiary 检查是否点赞
// @Summary 检查用户是否点赞日记
// @Description 检查当前用户是否已经点赞该日记
// @Tags 日记
// @Accept json
// @Produce json
// @Param id path string true "日记ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary/{id}/like [get]
func (h *DiaryHandler) CheckLikeDiary(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析日记ID
	idStr := c.Param("id")
	diaryId, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 调用仓库方法检查点赞状态
	isLiked, err := h.repo.CheckUserLike(diaryId, userID.(uuid.UUID))
	if err != nil {
		response.ServerError(c, "检查点赞状态失败: "+err.Error())
		return
	}

	response.Success(c, map[string]bool{"is_liked": isLiked})
}
