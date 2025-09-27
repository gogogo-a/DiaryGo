package repository

import (
	"fmt"
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

// BillWithTags 账单及其标签
type BillWithTags struct {
	Bill models.Bill `json:"bill"`
	Tags []TagInfo   `json:"tags"`
}

// TagInfo 标签信息
type TagInfo struct {
	ID      string `json:"id"`
	TagName string `json:"tag_name"`
	Type    string `json:"type"`
}

// BillRepository 账单仓库接口
type BillRepository interface {
	// Create 创建账单
	Create(bill *models.Bill, tagIDs []uuid.UUID) error

	// GetByID 根据ID获取账单
	GetByID(id uuid.UUID) (*models.Bill, error)

	// GetBillWithTags 获取账单及其标签
	GetBillWithTags(id uuid.UUID) (*models.Bill, []models.Tag, error)

	// GetBills 获取账单列表，支持分页和多条件搜索
	GetBills(
		accountBookID uuid.UUID,
		page, pageSize int,
		billType string,
		tagIDs []uuid.UUID,
		startTime, endTime time.Time,
		minAmount, maxAmount float64,
		keyword string,
	) ([]BillWithTags, int64, error)

	// Update 更新账单
	Update(bill *models.Bill, tagIDs []uuid.UUID) error

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
func (r *billRepository) Create(bill *models.Bill, tagIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建账单
		if err := tx.Create(bill).Error; err != nil {
			return err
		}

		// 创建账单与标签的关联
		for _, tagID := range tagIDs {
			billTag := models.BillTag{
				BillId: bill.Id,
				TagId:  tagID,
			}
			if err := tx.Create(&billTag).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetByID 根据ID获取账单
func (r *billRepository) GetByID(id uuid.UUID) (*models.Bill, error) {
	var bill models.Bill
	err := r.db.Where("id = ?", id).
		Preload("AccountBook").
		Preload("User").
		First(&bill).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

// GetBillWithTags 获取账单及其标签
func (r *billRepository) GetBillWithTags(id uuid.UUID) (*models.Bill, []models.Tag, error) {
	var bill models.Bill
	// 添加Preload预加载关联表数据
	if err := r.db.Where("id = ?", id).
		Preload("AccountBook").
		Preload("User").
		First(&bill).Error; err != nil {
		return nil, nil, err
	}

	var tags []models.Tag
	if err := r.db.Table("tags").
		Joins("JOIN bill_tags ON tags.id = bill_tags.tag_id").
		Where("bill_tags.bill_id = ?", id).
		Find(&tags).Error; err != nil {
		return nil, nil, err
	}

	// 验证并确保ID值正确
	if bill.Id != id {
		// 如果ID不匹配，说明存在数据异常，记录日志并使用传入的ID
		bill.Id = id
	}

	// 额外查询确保账单的账本和用户关系是正确的
	var originalBill models.Bill
	if err := r.db.Select("account_book_id", "user_id").
		Where("id = ?", id).
		First(&originalBill).Error; err == nil {
		// 确保使用数据库中的原始值
		bill.AccountBookId = originalBill.AccountBookId
		bill.UserId = originalBill.UserId
	}

	return &bill, tags, nil
}

// GetBills 获取账单列表，支持分页和多条件搜索
func (r *billRepository) GetBills(
	accountBookID uuid.UUID,
	page, pageSize int,
	billType string,
	tagIDs []uuid.UUID,
	startTime, endTime time.Time,
	minAmount, maxAmount float64,
	keyword string,
) ([]BillWithTags, int64, error) {
	var bills []models.Bill
	var total int64

	// 构建查询
	query := r.db.Model(&models.Bill{}).Where("account_book_id = ?", accountBookID)

	// 添加过滤条件
	if billType != "" {
		query = query.Where("type = ?", billType)
	}

	// 标签过滤
	if len(tagIDs) > 0 {
		// 使用子查询查找包含所有指定标签的账单
		subQuery := r.db.Table("bill_tags").
			Select("bill_id").
			Where("tag_id IN ?", tagIDs).
			Group("bill_id").
			Having("COUNT(DISTINCT tag_id) = ?", len(tagIDs))

		query = query.Where("bills.id IN (?)", subQuery)
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

	// 分页查询，添加预加载
	offset := (page - 1) * pageSize
	if err := query.Preload("AccountBook").
		Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	var billWithTagsList []BillWithTags
	for _, bill := range bills {
		// 获取每个账单的标签
		var tags []models.Tag
		if err := r.db.Table("tags").
			Joins("JOIN bill_tags ON tags.id = bill_tags.tag_id").
			Where("bill_tags.bill_id = ?", bill.Id).
			Find(&tags).Error; err != nil {
			return nil, 0, err
		}

		// 转换标签格式
		var tagInfos []TagInfo
		for _, tag := range tags {
			tagInfos = append(tagInfos, TagInfo{
				ID:      tag.Id.String(),
				TagName: tag.TagName,
				Type:    tag.Type,
			})
		}

		billWithTagsList = append(billWithTagsList, BillWithTags{
			Bill: bill,
			Tags: tagInfos,
		})
	}

	return billWithTagsList, total, nil
}

// Update 更新账单
func (r *billRepository) Update(bill *models.Bill, tagIDs []uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 保存原始的账本ID和用户ID
		// originalAccountBookId := bill.AccountBookId
		// originalUserId := bill.UserId

		fmt.Println("更新前的bill", bill.AccountBookId, bill.UserId)
		// 更新账单基本信息，只更新允许修改的字段
		tx.Model(bill).UpdateColumn("amount", bill.Amount)
		tx.Model(bill).UpdateColumn("type", bill.Type)
		tx.Model(bill).UpdateColumn("remark", bill.Remark)
		tx.Model(bill).UpdateColumn("image_url", bill.ImageUrl)

		// 删除现有的标签关联
		if err := tx.Where("bill_id = ?", bill.Id).Delete(&models.BillTag{}).Error; err != nil {
			return err
		}

		// 创建新的标签关联
		for _, tagID := range tagIDs {
			billTag := models.BillTag{
				BillId: bill.Id,
				TagId:  tagID,
			}
			if err := tx.Create(&billTag).Error; err != nil {
				return err
			}
		}

		fmt.Println("更新后的bill", bill.AccountBookId, bill.UserId)
		return nil
	})
}

// Delete 删除账单
func (r *billRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除账单标签关联
		if err := tx.Where("bill_id = ?", id).Delete(&models.BillTag{}).Error; err != nil {
			return err
		}

		// 删除账单
		return tx.Delete(&models.Bill{}, id).Error
	})
}

// GetStats 获取账单统计
func (r *billRepository) GetStats(accountBookID uuid.UUID, startTime, endTime time.Time, groupBy string) (*BillStats, error) {
	// 创建统计结果
	stats := &BillStats{
		TagStats: make(map[string]float64),
	}

	// 为收入和支出分别创建基础查询
	incomeQuery := r.db.Model(&models.Bill{}).Where("account_book_id = ? AND type = ?", accountBookID, "income")
	expenseQuery := r.db.Model(&models.Bill{}).Where("account_book_id = ? AND type = ?", accountBookID, "expense")

	// 添加时间范围
	if !startTime.IsZero() {
		incomeQuery = incomeQuery.Where("created_at >= ?", startTime)
		expenseQuery = expenseQuery.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		incomeQuery = incomeQuery.Where("created_at <= ?", endTime)
		expenseQuery = expenseQuery.Where("created_at <= ?", endTime)
	}

	// 计算总收入
	if err := incomeQuery.Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalIncome).Error; err != nil {
		return nil, err
	}

	// 计算总支出
	if err := expenseQuery.Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalExpense).Error; err != nil {
		return nil, err
	}

	// 计算净额
	stats.NetAmount = stats.TotalIncome - stats.TotalExpense

	// 按标签统计 - 区分收入和支出
	var tagStats []struct {
		TagID  uuid.UUID
		Name   string
		Type   string
		Amount float64
	}

	if err := r.db.Table("bills").
		Select("tags.id as tag_id, tags.tag_name as name, bills.type as type, SUM(bills.amount) as amount").
		Joins("JOIN bill_tags ON bills.id = bill_tags.bill_id").
		Joins("JOIN tags ON bill_tags.tag_id = tags.id").
		Where("bills.account_book_id = ?", accountBookID).
		Group("tags.id, tags.tag_name, bills.type").
		Scan(&tagStats).Error; err != nil {
		return nil, err
	}

	// 处理标签统计，区分收入和支出
	for _, ts := range tagStats {
		if ts.Type == "income" {
			stats.TagStats[ts.Name+"收入"] = ts.Amount
		} else {
			stats.TagStats[ts.Name+"支出"] = ts.Amount
		}
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
			// 改用两个独立的查询，然后合并结果
			type GroupResult struct {
				GroupKey string
				Amount   float64
			}

			var incomeResults []GroupResult
			var expenseResults []GroupResult

			// 收入查询
			incomeSubQuery := r.db.Model(&models.Bill{}).
				Select(groupFormat+" as group_key, COALESCE(SUM(amount), 0) as amount").
				Where("account_book_id = ? AND type = ?", accountBookID, "income")

			if !startTime.IsZero() {
				incomeSubQuery = incomeSubQuery.Where("created_at >= ?", startTime)
			}
			if !endTime.IsZero() {
				incomeSubQuery = incomeSubQuery.Where("created_at <= ?", endTime)
			}

			if err := incomeSubQuery.Group(groupFormat).Scan(&incomeResults).Error; err != nil {
				return nil, err
			}

			// 支出查询
			expenseSubQuery := r.db.Model(&models.Bill{}).
				Select(groupFormat+" as group_key, COALESCE(SUM(amount), 0) as amount").
				Where("account_book_id = ? AND type = ?", accountBookID, "expense")

			if !startTime.IsZero() {
				expenseSubQuery = expenseSubQuery.Where("created_at >= ?", startTime)
			}
			if !endTime.IsZero() {
				expenseSubQuery = expenseSubQuery.Where("created_at <= ?", endTime)
			}

			if err := expenseSubQuery.Group(groupFormat).Scan(&expenseResults).Error; err != nil {
				return nil, err
			}

			// 合并结果
			groupResultMap := make(map[string]BillGroupStats)

			// 处理收入
			for _, result := range incomeResults {
				groupStats := groupResultMap[result.GroupKey]
				groupStats.GroupKey = result.GroupKey
				groupStats.Income = result.Amount
				groupResultMap[result.GroupKey] = groupStats
			}

			// 处理支出
			for _, result := range expenseResults {
				groupStats := groupResultMap[result.GroupKey]
				groupStats.GroupKey = result.GroupKey
				groupStats.Expense = result.Amount
				groupResultMap[result.GroupKey] = groupStats
			}

			// 计算净额并添加到结果中
			for key, stat := range groupResultMap {
				stat.NetAmount = stat.Income - stat.Expense
				groupResultMap[key] = stat
				stats.GroupStats = append(stats.GroupStats, stat)
			}
		}
	}

	return stats, nil
}
