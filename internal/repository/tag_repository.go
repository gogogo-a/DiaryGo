package repository

import (
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

// TagRepository 标签仓库接口
type TagRepository interface {
	// Create 创建标签
	Create(tag *models.Tag) error

	// GetByID 根据ID获取标签
	GetByID(id uuid.UUID) (*models.Tag, error)

	// GetAll 获取所有标签
	GetAll() ([]models.Tag, error)

	// Update 更新标签
	Update(tag *models.Tag) error

	// Delete 删除标签
	Delete(id uuid.UUID) error
}

// tagRepository 标签仓库实现
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓库
func NewTagRepository() TagRepository {
	return &tagRepository{
		db: database.GetDB(),
	}
}

// Create 创建标签
func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

// GetByID 根据ID获取标签
func (r *tagRepository) GetByID(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetAll 获取所有标签
func (r *tagRepository) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

// Update 更新标签
func (r *tagRepository) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 删除标签
func (r *tagRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Tag{}, id).Error
}
