package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

// BillStats 账单统计结果
type BillStats struct {
	TotalIncome  float64            `json:"total_income"`
	TotalExpense float64            `json:"total_expense"`
	NetAmount    float64            `json:"net_amount"`
	GroupStats   []BillGroupStats   `json:"group_stats,omitempty"`
	TagStats     map[string]float64 `json:"tag_stats,omitempty"`
}

// BillGroupStats 分组统计结果
type BillGroupStats struct {
	GroupKey  string  `json:"group_key"`
	Income    float64 `json:"income"`
	Expense   float64 `json:"expense"`
	NetAmount float64 `json:"net_amount"`
}

// BillRepository 账单仓库接口
type BillRepository interface {
	// Create 创建账单
	Create(bill *models.Bill) error

	// GetByID 根据ID获取账单
	GetByID(id uuid.UUID) (*models.Bill, error)

	// GetBills 获取账单列表，支持分页和多条件搜索
	GetBills(
		accountBookID uuid.UUID,
		page, pageSize int,
		billType string,
		tagID uuid.UUID,
		startTime, endTime time.Time,
		minAmount, maxAmount float64,
		keyword string,
	) ([]models.Bill, int64, error)

	// Update 更新账单
	Update(bill *models.Bill) error

	// Delete 删除账单
	Delete(id uuid.UUID) error

	// GetStats 获取账单统计
	GetStats(accountBookID uuid.UUID, startTime, endTime time.Time, groupBy string) (*BillStats, error)
}

// billRepository 账单仓库实现
type billRepository struct {
	db *gorm.DB
}

// NewBillRepository 创建账单仓库
func NewBillRepository() BillRepository {
	return &billRepository{
		db: database.GetDB(),
	}
}

// Create 创建账单
func (r *billRepository) Create(bill *models.Bill) error {
	return r.db.Create(bill).Error
}

// GetByID 根据ID获取账单
func (r *billRepository) GetByID(id uuid.UUID) (*models.Bill, error) {
	var bill models.Bill
	err := r.db.Where("id = ?", id).First(&bill).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

// GetBills 获取账单列表，支持分页和多条件搜索
func (r *billRepository) GetBills(
	accountBookID uuid.UUID,
	page, pageSize int,
	billType string,
	tagID uuid.UUID,
	startTime, endTime time.Time,
	minAmount, maxAmount float64,
	keyword string,
) ([]models.Bill, int64, error) {
	var bills []models.Bill
	var total int64

	// 构建查询
	query := r.db.Model(&models.Bill{}).Where("account_book_id = ?", accountBookID)

	// 添加过滤条件
	if billType != "" {
		query = query.Where("type = ?", billType)
	}

	// 标签过滤
	if tagID != uuid.Nil {
		query = query.Where("tag_id = ?", tagID)
	}

	// 时间范围
	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	// 金额范围
	if minAmount > 0 {
		query = query.Where("amount >= ?", minAmount)
	}
	if maxAmount > 0 {
		query = query.Where("amount <= ?", maxAmount)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("remark LIKE ?", "%"+keyword+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, total, nil
}

// Update 更新账单
func (r *billRepository) Update(bill *models.Bill) error {
	return r.db.Save(bill).Error
}

// Delete 删除账单
func (r *billRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Bill{}, id).Error
}

// GetStats 获取账单统计
func (r *billRepository) GetStats(accountBookID uuid.UUID, startTime, endTime time.Time, groupBy string) (*BillStats, error) {
	// 创建统计结果
	stats := &BillStats{
		TagStats: make(map[string]float64),
	}

	// 基础查询
	query := r.db.Model(&models.Bill{}).Where("account_book_id = ?", accountBookID)

	// 添加时间范围
	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	// 计算总收入
	if err := query.Where("type = ?", "income").Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalIncome).Error; err != nil {
		return nil, err
	}

	// 计算总支出
	if err := query.Where("type = ?", "expense").Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalExpense).Error; err != nil {
		return nil, err
	}

	// 计算净额
	stats.NetAmount = stats.TotalIncome - stats.TotalExpense

	// 按标签统计
	var tagStats []struct {
		TagId  uuid.UUID
		Name   string
		Amount float64
	}
	if err := r.db.Table("bills").
		Select("bills.tag_id, tags.name, SUM(bills.amount) as amount").
		Joins("JOIN tags ON bills.tag_id = tags.id").
		Where("bills.account_book_id = ?", accountBookID).
		Group("bills.tag_id, tags.name").
		Scan(&tagStats).Error; err != nil {
		return nil, err
	}

	for _, ts := range tagStats {
		stats.TagStats[ts.Name] = ts.Amount
	}

	// 按时间分组统计
	if groupBy != "" {
		var groupFormat string
		switch groupBy {
		case "day":
			groupFormat = "DATE(created_at)"
		case "week":
			groupFormat = "YEARWEEK(created_at, 1)"
		case "month":
			groupFormat = "DATE_FORMAT(created_at, '%Y-%m')"
		case "year":
			groupFormat = "YEAR(created_at)"
		}

		if groupFormat != "" {
			var groupStats []struct {
				GroupKey string
				Income   float64
				Expense  float64
			}

			// 按时间分组查询收支情况
			subQuery := r.db.Table("(?) as income_table, (?) as expense_table",
				r.db.Model(&models.Bill{}).
					Select(groupFormat+" as group_key, COALESCE(SUM(amount), 0) as income").
					Where("account_book_id = ? AND type = ?", accountBookID, "income").
					Group(groupFormat),
				r.db.Model(&models.Bill{}).
					Select(groupFormat+" as group_key, COALESCE(SUM(amount), 0) as expense").
					Where("account_book_id = ? AND type = ?", accountBookID, "expense").
					Group(groupFormat),
			).Select("income_table.group_key, income_table.income, expense_table.expense").
				Where("income_table.group_key = expense_table.group_key")

			// 添加时间范围
			if !startTime.IsZero() {
				subQuery = subQuery.Where("income_table.created_at >= ?", startTime)
			}
			if !endTime.IsZero() {
				subQuery = subQuery.Where("income_table.created_at <= ?", endTime)
			}

			if err := subQuery.Scan(&groupStats).Error; err != nil {
				return nil, err
			}

			// 转换为返回格式
			for _, gs := range groupStats {
				stats.GroupStats = append(stats.GroupStats, BillGroupStats{
					GroupKey:  gs.GroupKey,
					Income:    gs.Income,
					Expense:   gs.Expense,
					NetAmount: gs.Income - gs.Expense,
				})
			}
		}
	}

	return stats, nil
}
