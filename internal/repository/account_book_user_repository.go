package repository

import (
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

type AccountBookUserRepository interface {
	// Create 创建账本用户关联（授予用户账本权限）
	Create(accountBookUser *models.AccountBookUser) error

	// GetByID 根据ID获取账本用户关联
	GetByID(id uuid.UUID) (*models.AccountBookUser, error)

	// GetByUserID 获取用户的账本关联
	GetByUserID(userID uuid.UUID) (*models.AccountBookUser, error)

	// GetByAccountBookID 获取账本的用户关联
	GetByAccountBookID(accountBookID uuid.UUID) (*models.AccountBookUser, error)

	// GetByAccountBookIDAndUserID 根据账本ID和用户ID获取特定的账本用户关联
	GetByAccountBookIDAndUserID(accountBookID, userID uuid.UUID) (*models.AccountBookUser, error)

	// GetAllUsersByAccountBookID 获取账本的所有用户
	GetAllUsersByAccountBookID(accountBookID uuid.UUID) ([]models.User, error)

	// Delete 根据ID删除账本用户关联
	Delete(id uuid.UUID) error

	// DeleteByUserID 删除用户的所有账本关联
	DeleteByUserID(userID uuid.UUID) error

	// DeleteByAccountBookID 删除账本的所有用户关联
	DeleteByAccountBookID(accountBookID uuid.UUID) error

	// DeleteByAccountBookIDAndUserID 删除特定用户的特定账本权限
	DeleteByAccountBookIDAndUserID(accountBookID, userID uuid.UUID) error
}

// accountBookUserRepository 账本用户关联仓库实现
type accountBookUserRepository struct {
	db *gorm.DB
}

// NewAccountBookUserRepository 创建账本用户关联仓库
func NewAccountBookUserRepository() AccountBookUserRepository {
	return &accountBookUserRepository{
		db: database.GetDB(),
	}
}

func (r *accountBookUserRepository) Create(accountBookUser *models.AccountBookUser) error {
	return r.db.Create(accountBookUser).Error
}

func (r *accountBookUserRepository) GetByID(id uuid.UUID) (*models.AccountBookUser, error) {
	var accountBookUser models.AccountBookUser
	err := r.db.Where("id = ?", id).First(&accountBookUser).Error
	if err != nil {
		return nil, err
	}
	return &accountBookUser, nil
}

func (r *accountBookUserRepository) GetByUserID(userID uuid.UUID) (*models.AccountBookUser, error) {
	var accountBookUser models.AccountBookUser
	err := r.db.Where("user_id = ?", userID).First(&accountBookUser).Error
	if err != nil {
		return nil, err
	}
	return &accountBookUser, nil
}

func (r *accountBookUserRepository) GetByAccountBookID(accountBookID uuid.UUID) (*models.AccountBookUser, error) {
	var accountBookUser models.AccountBookUser
	err := r.db.Where("account_book_id = ?", accountBookID).First(&accountBookUser).Error
	if err != nil {
		return nil, err
	}
	return &accountBookUser, nil
}

func (r *accountBookUserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.AccountBookUser{}, id).Error
}

func (r *accountBookUserRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Delete(&models.AccountBookUser{}, "user_id = ?", userID).Error
}

func (r *accountBookUserRepository) DeleteByAccountBookID(accountBookID uuid.UUID) error {
	return r.db.Delete(&models.AccountBookUser{}, "account_book_id = ?", accountBookID).Error
}

// GetAllUsersByAccountBookID 获取账本的所有用户
func (r *accountBookUserRepository) GetAllUsersByAccountBookID(accountBookID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.Table("users").
		Joins("JOIN account_book_users ON users.id = account_book_users.user_id").
		Where("account_book_users.account_book_id = ?", accountBookID).
		Find(&users).Error
	return users, err
}

// GetByAccountBookIDAndUserID 根据账本ID和用户ID获取账本用户关联
func (r *accountBookUserRepository) GetByAccountBookIDAndUserID(accountBookID, userID uuid.UUID) (*models.AccountBookUser, error) {
	var accountBookUser models.AccountBookUser
	err := r.db.Where("account_book_id = ? AND user_id = ?", accountBookID, userID).First(&accountBookUser).Error
	if err != nil {
		return nil, err
	}
	return &accountBookUser, nil
}

// DeleteByAccountBookIDAndUserID 根据账本ID和用户ID删除账本用户关联
func (r *accountBookUserRepository) DeleteByAccountBookIDAndUserID(accountBookID, userID uuid.UUID) error {
	return r.db.Where("account_book_id = ? AND user_id = ?", accountBookID, userID).Delete(&models.AccountBookUser{}).Error
}
