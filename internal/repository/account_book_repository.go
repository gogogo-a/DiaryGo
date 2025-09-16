package repository

import (
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

type AccountBookRepository interface {
	// Create 创建新的账本
	Create(accountBook *models.AccountBook) error

	// GetByID 根据ID获取账本
	GetByID(id uuid.UUID) (*models.AccountBook, error)

	// GetAll 获取用户的所有账本
	GetAll(userID uuid.UUID) ([]models.AccountBook, error)

	// Update 更新账本信息
	Update(accountBook *models.AccountBook) error

	// Delete 删除账本
	Delete(id uuid.UUID) error
}

type accountBookRepository struct {
	db *gorm.DB
}

func NewAccountBookRepository() AccountBookRepository {
	return &accountBookRepository{
		db: database.GetDB(),
	}
}

func (r *accountBookRepository) Create(accountBook *models.AccountBook) error {
	return r.db.Create(accountBook).Error
}

func (r *accountBookRepository) GetByID(id uuid.UUID) (*models.AccountBook, error) {
	var accountBook models.AccountBook
	err := r.db.Where("id = ?", id).First(&accountBook).Error
	if err != nil {
		return nil, err
	}
	return &accountBook, nil
}

func (r *accountBookRepository) GetAll(userID uuid.UUID) ([]models.AccountBook, error) {
	var accountBooks []models.AccountBook
	err := r.db.Where("user_id = ?", userID).Find(&accountBooks).Error
	if err != nil {
		return nil, err
	}
	return accountBooks, nil
}

func (r *accountBookRepository) Update(accountBook *models.AccountBook) error {
	return r.db.Save(accountBook).Error
}

func (r *accountBookRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.AccountBook{}, id).Error
}
