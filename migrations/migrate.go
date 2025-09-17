package migrations

import (
	"fmt"

	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/logger"
	"gorm.io/gorm"
)

// MigrateDatabase 使用分组的模型切片迁移所有数据库表
func MigrateDatabase(db *gorm.DB) error {
	// 按功能分组的模型
	modelGroups := []struct {
		name   string
		models []interface{}
	}{
		{
			name: "用户",
			models: []interface{}{
				&models.User{},
			},
		},
		{
			name: "日记核心",
			models: []interface{}{
				&models.Diary{},
				&models.DiaryUser{},
			},
		},
		{
			name: "日记媒体",
			models: []interface{}{
				&models.DiaryImage{},
				&models.DiaryVideo{},
			},
		},
		{
			name: "日记分类和权限",
			models: []interface{}{
				&models.Tag{},
				&models.DiaryTag{},
				&models.DPermission{},
				&models.DiaryDPermission{},
			},
		},
		{
			name: "交互",
			models: []interface{}{
				&models.DiaryLike{},
			},
		},
		{
			name: "账本",
			models: []interface{}{
				&models.AccountBook{},
				&models.AccountBookUser{},
				&models.Bill{},
				&models.BillTag{},
			},
		},
	}

	logger.Info("开始数据库迁移...")

	// 迁移每个分组的模型
	for _, group := range modelGroups {
		logger.Info("正在迁移%s相关表...", group.name)

		if err := db.AutoMigrate(group.models...); err != nil {
			logger.Error("迁移%s表失败: %v", group.name, err)
			return fmt.Errorf("迁移%s表失败: %w", group.name, err)
		}

		logger.Info("成功迁移%s相关表", group.name)
	}

	logger.Info("数据库迁移完成")
	return nil
}
