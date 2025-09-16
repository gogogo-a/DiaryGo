package repository

import (
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

type UserLoginRepository interface {
	// Login 用户登录
	Login(user *models.User) error

	// FindByOpenID 根据OpenID查找用户
	FindByOpenID(openID string) (*models.User, error)

	// Create 创建新用户
	Create(user *models.User) error

	// GetByID 根据ID获取用户
	GetByID(id uuid.UUID) (*models.User, error)
}

// userLoginRepository 用户登录仓库实现
type userLoginRepository struct {
	db *gorm.DB
}

// NewUserLoginRepository 创建用户登录仓库
func NewUserLoginRepository() UserLoginRepository {
	return &userLoginRepository{
		db: database.GetDB(),
	}
}

// Login 用户登录
func (r *userLoginRepository) Login(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByOpenID 根据OpenID查找用户
func (r *userLoginRepository) FindByOpenID(openID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("plant_id = ? AND plant_form = ?", openID, "wechat").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *userLoginRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userLoginRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
