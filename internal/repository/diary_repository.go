package repository

import (
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

// DiaryRepository 日记仓库接口
type DiaryRepository interface {
	Create(diary *models.Diary) error
	GetByID(id uuid.UUID) (*models.Diary, error)
	GetAll(userID uuid.UUID) ([]models.Diary, error)
	Update(diary *models.Diary) error
	Delete(id uuid.UUID) error
}

// diaryRepository 日记仓库实现
type diaryRepository struct {
	db *gorm.DB
}

// NewDiaryRepository 创建日记仓库
func NewDiaryRepository() DiaryRepository {
	return &diaryRepository{
		db: database.GetDB(),
	}
}

// Create 创建日记
func (r *diaryRepository) Create(diary *models.Diary) error {
	return r.db.Create(diary).Error
}

// GetByID 根据ID获取日记
func (r *diaryRepository) GetByID(id uuid.UUID) (*models.Diary, error) {
	var diary models.Diary
	err := r.db.First(&diary, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &diary, nil
}

// GetAll 获取用户的所有日记
func (r *diaryRepository) GetAll(userID uuid.UUID) ([]models.Diary, error) {
	var diaries []models.Diary
	err := r.db.Where("user_id = ?", userID).Find(&diaries).Error
	return diaries, err
}

// Update 更新日记
func (r *diaryRepository) Update(diary *models.Diary) error {
	return r.db.Save(diary).Error
}

// Delete 删除日记
func (r *diaryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Diary{}, "id = ?", id).Error
}
