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
	err := r.db.
		Table("account_books").
		Joins("JOIN account_book_users ON account_books.id = account_book_users.account_book_id").
		Where("account_book_users.user_id = ?", userID).
		Find(&accountBooks).Error

	if err != nil {
		return nil, err
	}
	return accountBooks, nil
}

func (r *accountBookRepository) Update(accountBook *models.AccountBook) error {
	// 确保有明确的WHERE条件，只更新特定ID的账本
	// return r.db.Model(&models.AccountBook{}).Where("id = ?", accountBook.Id).UpdateColumn("name", accountBook.Name).Error
	return r.db.Save(accountBook).Error
}

func (r *accountBookRepository) Delete(id uuid.UUID) error {
	// 使用事务确保操作的原子性
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除关联的account_book_users记录
		if err := tx.Where("account_book_id = ?", id).Delete(&models.AccountBookUser{}).Error; err != nil {
			return err
		}

		// 然后删除账本记录
		if err := tx.Delete(&models.AccountBook{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}
