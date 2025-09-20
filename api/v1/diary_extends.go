package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/repository"
	"github.com/haogeng/DiaryGo/pkg/response"
	"gorm.io/gorm"
)

// DiaryExtendsHandler 日记扩展处理器，处理日记相关的图片和视频
type DiaryExtendsHandler struct {
	repo repository.DiaryExtendsRepository
}

// NewDiaryExtendsHandler 创建日记扩展处理器
func NewDiaryExtendsHandler() *DiaryExtendsHandler {
	return &DiaryExtendsHandler{
		repo: repository.NewDiaryExtendsRepository(),
	}
}

// RegisterRoutes 注册日记扩展相关路由
func (h *DiaryExtendsHandler) RegisterRoutes(router *gin.RouterGroup) {
	diaryExtends := router.Group("/diary-extends")
	{
		// 图片相关路由
		diaryExtends.POST("/images", h.AddImage)                // 添加图片
		diaryExtends.DELETE("/images/:id", h.DeleteImage)       // 删除图片
		diaryExtends.GET("/images/diary/:diaryId", h.GetImages) // 获取日记图片列表

		// 视频相关路由
		diaryExtends.POST("/videos", h.AddVideo)                // 添加视频
		diaryExtends.DELETE("/videos/:id", h.DeleteVideo)       // 删除视频
		diaryExtends.GET("/videos/diary/:diaryId", h.GetVideos) // 获取日记视频列表
	}
}

// AddImageRequest 添加图片请求参数
type AddImageRequest struct {
	DiaryId  string `json:"diary_id" binding:"required"`
	ImageUrl string `json:"image_url" binding:"required"`
}

// AddImage 添加日记图片
// @Summary 添加日记图片
// @Description 向指定日记添加图片URL
// @Tags 日记扩展
// @Accept json
// @Produce json
// @Param request body AddImageRequest true "图片信息"
// @Success 201 {object} response.Response{data=models.DiaryImage}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary-extends/images [post]
func (h *DiaryExtendsHandler) AddImage(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析请求
	var req AddImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 转换日记ID
	diaryId, err := uuid.Parse(req.DiaryId)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 检查用户是否有权限操作此日记
	if err := h.repo.CheckDiaryPermission(diaryId, userID.(uuid.UUID)); err != nil {
		response.Forbidden(c, err.Error())
		return
	}

	// 添加图片
	image, err := h.repo.AddImage(diaryId, req.ImageUrl)
	if err != nil {
		response.ServerError(c, "添加图片失败")
		return
	}

	// 创建精简响应对象，避免返回空的日记对象
	result := map[string]interface{}{
		"id":        image.Id,
		"diary_id":  image.DiaryId,
		"image_url": image.ImageUrl,
	}

	response.SuccessWithMessage(c, "添加图片成功", result)
}

// DeleteImage 删除日记图片
// @Summary 删除日记图片
// @Description 删除指定ID的图片
// @Tags 日记扩展
// @Accept json
// @Produce json
// @Param id path string true "图片ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary-extends/images/{id} [delete]
func (h *DiaryExtendsHandler) DeleteImage(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析图片ID
	idStr := c.Param("id")
	imageId, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的图片ID格式")
		return
	}

	// 查询图片信息，验证权限
	image, err := h.repo.GetImageByID(imageId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "图片不存在")
		} else {
			response.ServerError(c, "查询图片失败")
		}
		return
	}

	// 检查用户是否有权限删除此图片
	if err := h.repo.CheckDiaryPermission(image.DiaryId, userID.(uuid.UUID)); err != nil {
		response.Forbidden(c, err.Error())
		return
	}

	// 删除图片
	if err := h.repo.DeleteImage(imageId); err != nil {
		response.ServerError(c, "删除图片失败")
		return
	}

	response.SuccessWithMessage(c, "删除图片成功", nil)
}

// GetImages 获取日记图片列表
// @Summary 获取日记图片列表
// @Description 获取指定日记的所有图片URL
// @Tags 日记扩展
// @Accept json
// @Produce json
// @Param diaryId path string true "日记ID"
// @Success 200 {object} response.Response{data=[]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary-extends/images/diary/{diaryId} [get]
func (h *DiaryExtendsHandler) GetImages(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析日记ID
	diaryIdStr := c.Param("diaryId")
	diaryId, err := uuid.Parse(diaryIdStr)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 检查用户是否有权限查看此日记
	if err := h.repo.CheckDiaryPermission(diaryId, userID.(uuid.UUID)); err != nil {
		response.Forbidden(c, err.Error())
		return
	}

	// 获取日记图片列表
	images, err := h.repo.GetDiaryImages(diaryId)
	if err != nil {
		response.ServerError(c, "获取图片列表失败")
		return
	}

	// 只提取图片URL列表
	var imageUrls []string
	for _, image := range images {
		imageUrls = append(imageUrls, image.ImageUrl)
	}

	// 返回图片URL列表
	response.Success(c, imageUrls)
}

// AddVideoRequest 添加视频请求参数
type AddVideoRequest struct {
	DiaryId  string `json:"diary_id" binding:"required"`
	VideoUrl string `json:"video_url" binding:"required"`
}

// AddVideo 添加日记视频
// @Summary 添加日记视频
// @Description 向指定日记添加视频URL
// @Tags 日记扩展
// @Accept json
// @Produce json
// @Param request body AddVideoRequest true "视频信息"
// @Success 201 {object} response.Response{data=models.DiaryVideo}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary-extends/videos [post]
func (h *DiaryExtendsHandler) AddVideo(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析请求
	var req AddVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 转换日记ID
	diaryId, err := uuid.Parse(req.DiaryId)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 检查用户是否有权限操作此日记
	if err := h.repo.CheckDiaryPermission(diaryId, userID.(uuid.UUID)); err != nil {
		response.Forbidden(c, err.Error())
		return
	}

	// 添加视频
	video, err := h.repo.AddVideo(diaryId, req.VideoUrl)
	if err != nil {
		response.ServerError(c, "添加视频失败")
		return
	}

	// 创建精简响应对象，避免返回空的日记对象
	result := map[string]interface{}{
		"id":        video.Id,
		"diary_id":  video.DiaryId,
		"video_url": video.VideoUrl,
	}

	response.SuccessWithMessage(c, "添加视频成功", result)
}

// DeleteVideo 删除日记视频
// @Summary 删除日记视频
// @Description 删除指定ID的视频
// @Tags 日记扩展
// @Accept json
// @Produce json
// @Param id path string true "视频ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary-extends/videos/{id} [delete]
func (h *DiaryExtendsHandler) DeleteVideo(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析视频ID
	idStr := c.Param("id")
	videoId, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的视频ID格式")
		return
	}

	// 查询视频信息，验证权限
	video, err := h.repo.GetVideoByID(videoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "视频不存在")
		} else {
			response.ServerError(c, "查询视频失败")
		}
		return
	}

	// 检查用户是否有权限删除此视频
	if err := h.repo.CheckDiaryPermission(video.DiaryId, userID.(uuid.UUID)); err != nil {
		response.Forbidden(c, err.Error())
		return
	}

	// 删除视频
	if err := h.repo.DeleteVideo(videoId); err != nil {
		response.ServerError(c, "删除视频失败")
		return
	}

	response.SuccessWithMessage(c, "删除视频成功", nil)
}

// GetVideos 获取日记视频列表
// @Summary 获取日记视频列表
// @Description 获取指定日记的所有视频URL
// @Tags 日记扩展
// @Accept json
// @Produce json
// @Param diaryId path string true "日记ID"
// @Success 200 {object} response.Response{data=[]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/diary-extends/videos/diary/{diaryId} [get]
func (h *DiaryExtendsHandler) GetVideos(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析日记ID
	diaryIdStr := c.Param("diaryId")
	diaryId, err := uuid.Parse(diaryIdStr)
	if err != nil {
		response.ParamError(c, "无效的日记ID格式")
		return
	}

	// 检查用户是否有权限查看此日记
	if err := h.repo.CheckDiaryPermission(diaryId, userID.(uuid.UUID)); err != nil {
		response.Forbidden(c, err.Error())
		return
	}

	// 获取日记视频列表
	videos, err := h.repo.GetDiaryVideos(diaryId)
	if err != nil {
		response.ServerError(c, "获取视频列表失败")
		return
	}

	// 只提取视频URL列表
	var videoUrls []string
	for _, video := range videos {
		videoUrls = append(videoUrls, video.VideoUrl)
	}

	// 返回视频URL列表
	response.Success(c, videoUrls)
}
