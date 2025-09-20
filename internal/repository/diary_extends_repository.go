package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

// DiaryExtendsRepository 日记扩展仓库接口，处理日记相关的图片和视频
type DiaryExtendsRepository interface {
	// 图片相关操作
	AddImage(diaryId uuid.UUID, imageUrl string) (*models.DiaryImage, error)
	DeleteImage(imageId uuid.UUID) error
	GetDiaryImages(diaryId uuid.UUID) ([]models.DiaryImage, error)
	GetImageByID(imageId uuid.UUID) (*models.DiaryImage, error)

	// 视频相关操作
	AddVideo(diaryId uuid.UUID, videoUrl string) (*models.DiaryVideo, error)
	DeleteVideo(videoId uuid.UUID) error
	GetDiaryVideos(diaryId uuid.UUID) ([]models.DiaryVideo, error)
	GetVideoByID(videoId uuid.UUID) (*models.DiaryVideo, error)

	// 验证权限
	CheckDiaryPermission(diaryId uuid.UUID, userId uuid.UUID) error
}

// diaryExtendsRepository 日记扩展仓库实现
type diaryExtendsRepository struct {
	db *gorm.DB
}

// NewDiaryExtendsRepository 创建日记扩展仓库
func NewDiaryExtendsRepository() DiaryExtendsRepository {
	return &diaryExtendsRepository{
		db: database.GetDB(),
	}
}

// CheckDiaryPermission 验证用户是否有权限操作此日记
func (r *diaryExtendsRepository) CheckDiaryPermission(diaryId uuid.UUID, userId uuid.UUID) error {
	var diaryUser models.DiaryUser
	if err := r.db.Where("diary_id = ? AND user_id = ?", diaryId, userId).First(&diaryUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("您没有权限操作此日记")
		}
		return err
	}
	return nil
}

// AddImage 添加日记图片
func (r *diaryExtendsRepository) AddImage(diaryId uuid.UUID, imageUrl string) (*models.DiaryImage, error) {
	// 创建图片记录
	diaryImage := &models.DiaryImage{
		DiaryId:  diaryId,
		ImageUrl: imageUrl,
	}

	if err := r.db.Create(diaryImage).Error; err != nil {
		return nil, err
	}

	return diaryImage, nil
}

// DeleteImage 删除日记图片
func (r *diaryExtendsRepository) DeleteImage(imageId uuid.UUID) error {
	return r.db.Delete(&models.DiaryImage{}, "id = ?", imageId).Error
}

// GetDiaryImages 获取日记的所有图片
func (r *diaryExtendsRepository) GetDiaryImages(diaryId uuid.UUID) ([]models.DiaryImage, error) {
	var images []models.DiaryImage
	if err := r.db.Where("diary_id = ?", diaryId).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

// AddVideo 添加日记视频
func (r *diaryExtendsRepository) AddVideo(diaryId uuid.UUID, videoUrl string) (*models.DiaryVideo, error) {
	// 创建视频记录
	diaryVideo := &models.DiaryVideo{
		DiaryId:  diaryId,
		VideoUrl: videoUrl,
	}

	if err := r.db.Create(diaryVideo).Error; err != nil {
		return nil, err
	}

	return diaryVideo, nil
}

// DeleteVideo 删除日记视频
func (r *diaryExtendsRepository) DeleteVideo(videoId uuid.UUID) error {
	return r.db.Delete(&models.DiaryVideo{}, "id = ?", videoId).Error
}

// GetDiaryVideos 获取日记的所有视频
func (r *diaryExtendsRepository) GetDiaryVideos(diaryId uuid.UUID) ([]models.DiaryVideo, error) {
	var videos []models.DiaryVideo
	if err := r.db.Where("diary_id = ?", diaryId).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// GetImageByID 根据ID获取图片
func (r *diaryExtendsRepository) GetImageByID(imageId uuid.UUID) (*models.DiaryImage, error) {
	var image models.DiaryImage
	if err := r.db.Where("id = ?", imageId).First(&image).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

// GetVideoByID 根据ID获取视频
func (r *diaryExtendsRepository) GetVideoByID(videoId uuid.UUID) (*models.DiaryVideo, error) {
	var video models.DiaryVideo
	if err := r.db.Where("id = ?", videoId).First(&video).Error; err != nil {
		return nil, err
	}
	return &video, nil
}
