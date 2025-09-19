package repository

import (
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

// DPermissionRepository 权限仓库接口
type DPermissionRepository interface {
	// GetAll 获取所有权限
	GetAll() ([]models.DPermission, error)

	// GetByID 根据ID获取权限
	GetByID(id uuid.UUID) (*models.DPermission, error)
}

// dpermissionRepository 权限仓库实现
type dpermissionRepository struct {
	db *gorm.DB
}

// NewDPermissionRepository 创建权限仓库
func NewDPermissionRepository() DPermissionRepository {
	return &dpermissionRepository{
		db: database.GetDB(),
	}
}

// GetAll 获取所有权限
func (r *dpermissionRepository) GetAll() ([]models.DPermission, error) {
	var permissions []models.DPermission
	if err := r.db.Order("created_at").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetByID 根据ID获取权限
func (r *dpermissionRepository) GetByID(id uuid.UUID) (*models.DPermission, error) {
	var permission models.DPermission
	if err := r.db.Where("id = ?", id).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}
