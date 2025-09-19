package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

// 定义标签分类常量
const (
	TagCategoryBill  = "bill"
	TagCategoryDiary = "diary"
)

// TagRepository 标签仓库接口
type TagRepository interface {
	// Create 创建标签
	Create(tag *models.Tag) error

	// GetByID 根据ID获取标签
	GetByID(id uuid.UUID) (*models.Tag, error)

	// GetAll 获取所有标签，可按分类过滤
	GetAll(category string) ([]models.Tag, error)

	// GetByName 根据名称和分类获取标签
	GetByName(name string, category string) (*models.Tag, error)

	// Update 更新标签
	Update(tag *models.Tag) error

	// Delete 删除标签
	Delete(id uuid.UUID) error

	// // BatchCreate 批量创建标签
	// BatchCreate(tags []*models.Tag) error
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
	// 验证标签分类是否有效
	if !isValidCategory(tag.Category) {
		return errors.New("无效的标签分类")
	}

	// 检查同名标签是否已存在
	var existingTag models.Tag
	if err := r.db.Where("tag_name = ? AND category = ?", tag.TagName, tag.Category).First(&existingTag).Error; err == nil {
		return errors.New("同名标签已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.Create(tag).Error
}

// GetByID 根据ID获取标签
func (r *tagRepository) GetByID(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("标签不存在")
		}
		return nil, err
	}
	return &tag, nil
}

// GetAll 获取所有标签，可按分类过滤
func (r *tagRepository) GetAll(category string) ([]models.Tag, error) {
	var tags []models.Tag
	query := r.db.Order("created_at DESC")

	// 按分类过滤
	if category != "" {
		if !isValidCategory(category) {
			return nil, errors.New("无效的标签分类")
		}
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetByName 根据名称和分类获取标签
func (r *tagRepository) GetByName(name string, category string) (*models.Tag, error) {
	if !isValidCategory(category) {
		return nil, errors.New("无效的标签分类")
	}

	var tag models.Tag
	err := r.db.Where("tag_name = ? AND category = ?", name, category).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("标签不存在")
		}
		return nil, err
	}
	return &tag, nil
}

// Update 更新标签
func (r *tagRepository) Update(tag *models.Tag) error {
	// 验证标签分类是否有效
	if !isValidCategory(tag.Category) {
		return errors.New("无效的标签分类")
	}

	// 检查标签是否存在
	var existingTag models.Tag
	if err := r.db.Where("id = ?", tag.Id).First(&existingTag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return err
	}

	// 检查更新后的标签名是否与其他标签冲突
	if existingTag.TagName != tag.TagName {
		var duplicateTag models.Tag
		if err := r.db.Where("tag_name = ? AND category = ? AND id <> ?", tag.TagName, tag.Category, tag.Id).First(&duplicateTag).Error; err == nil {
			return errors.New("同名标签已存在")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	return r.db.Save(tag).Error
}

// Delete 删除标签
func (r *tagRepository) Delete(id uuid.UUID) error {
	// 检查标签是否存在
	var tag models.Tag
	if err := r.db.Where("id = ?", id).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("标签不存在")
		}
		return err
	}

	// 检查标签是否被使用
	var billTagCount int64
	if err := r.db.Model(&models.BillTag{}).Where("tag_id = ?", id).Count(&billTagCount).Error; err != nil {
		return err
	}

	var diaryTagCount int64
	if err := r.db.Model(&models.DiaryTag{}).Where("tag_id = ?", id).Count(&diaryTagCount).Error; err != nil {
		return err
	}

	if billTagCount > 0 || diaryTagCount > 0 {
		return errors.New("标签已被使用，无法删除")
	}

	// 执行删除
	return r.db.Delete(&tag).Error
}

// // BatchCreate 批量创建标签
// func (r *tagRepository) BatchCreate(tags []*models.Tag) error {
// 	return r.db.Transaction(func(tx *gorm.DB) error {
// 		for _, tag := range tags {
// 			// 验证标签分类是否有效
// 			if !isValidCategory(tag.Category) {
// 				return errors.New("无效的标签分类: " + tag.Category)
// 			}

// 			// 检查同名标签是否已存在
// 			var existingTag models.Tag
// 			if err := tx.Where("tag_name = ? AND category = ?", tag.TagName, tag.Category).First(&existingTag).Error; err == nil {
// 				// 如果已存在，跳过这个标签
// 				continue
// 			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
// 				return err
// 			}

// 			// 创建标签
// 			if err := tx.Create(tag).Error; err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

// 验证标签分类是否有效
func isValidCategory(category string) bool {
	return category == TagCategoryBill || category == TagCategoryDiary
}
