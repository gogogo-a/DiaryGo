package repository

import (
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	// GetUsers 获取用户列表，支持分页和搜索
	GetUsers(page, pageSize int, keyword string) ([]models.User, int64, error)

	// GetByID 根据ID获取用户
	GetByID(id uuid.UUID) (*models.User, error)
	//Update 更新用户
	Update(user *models.User) error
}

// userRepository 用户仓库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.GetDB(),
	}
}

// GetUsers 获取用户列表，支持分页和搜索
func (r *userRepository) GetUsers(page, pageSize int, keyword string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 构建查询
	query := r.db.Model(&models.User{})

	// 如果有关键词，添加搜索条件
	if keyword != "" {
		query = query.Where("user_name LIKE ? OR nick_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *models.User) error {
	return r.db.Model(&models.User{}).Where("id = ?", user.Id).Updates(user).Error
}