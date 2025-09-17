package scripts

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

// 预定义的标签数据
var predefinedTags = []struct {
	TagName  string
	Type     string
	Category string
}{
	// 账单标签
	{TagName: "工资", Type: "收入", Category: "bill"},
	{TagName: "奖金", Type: "收入", Category: "bill"},
	{TagName: "兼职", Type: "收入", Category: "bill"},
	{TagName: "投资", Type: "收入", Category: "bill"},
	{TagName: "餐饮", Type: "支出", Category: "bill"},
	{TagName: "购物", Type: "支出", Category: "bill"},
	{TagName: "交通", Type: "支出", Category: "bill"},
	{TagName: "住宿", Type: "支出", Category: "bill"},
	{TagName: "医疗", Type: "支出", Category: "bill"},
	{TagName: "教育", Type: "支出", Category: "bill"},
	{TagName: "娱乐", Type: "支出", Category: "bill"},
	{TagName: "其他", Type: "支出", Category: "bill"},

	// 日记标签
	{TagName: "旅行", Type: "生活", Category: "diary"},
	{TagName: "美食", Type: "生活", Category: "diary"},
	{TagName: "工作", Type: "职场", Category: "diary"},
	{TagName: "学习", Type: "教育", Category: "diary"},
	{TagName: "健身", Type: "健康", Category: "diary"},
	{TagName: "感悟", Type: "思考", Category: "diary"},
	{TagName: "电影", Type: "娱乐", Category: "diary"},
	{TagName: "读书", Type: "教育", Category: "diary"},
	{TagName: "心情", Type: "情感", Category: "diary"},
	{TagName: "家庭", Type: "生活", Category: "diary"},
}

// 预定义的权限类型
var predefinedPermissions = []struct {
	PermissionName string
	Description    string
}{
	{PermissionName: "private", Description: "私密，仅创建者可见"},
	{PermissionName: "public", Description: "公开，所有人可见"},
	{PermissionName: "shared_read", Description: "共享，特定用户可查看"},
	{PermissionName: "shared_edit", Description: "共享，特定用户可编辑"},
}

// GenerateCommonData 生成与用户无关的公共数据
func GenerateCommonData() error {
	db := database.GetDB()

	// 生成标签数据
	if err := generateTags(db); err != nil {
		return fmt.Errorf("生成标签数据失败: %w", err)
	}

	// 生成权限类型数据
	if err := generatePermissions(db); err != nil {
		return fmt.Errorf("生成权限类型数据失败: %w", err)
	}

	log.Println("公共数据生成完成")
	return nil
}

// 生成标签数据
func generateTags(db *gorm.DB) error {
	log.Println("开始生成标签数据...")

	for _, tagData := range predefinedTags {
		// 检查标签是否已存在
		var count int64
		db.Model(&models.Tag{}).Where("tag_name = ? AND category = ?", tagData.TagName, tagData.Category).Count(&count)

		if count == 0 {
			tag := models.Tag{
				Id:        uuid.New(),
				TagName:   tagData.TagName,
				Type:      tagData.Type,
				Category:  tagData.Category,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := db.Create(&tag).Error; err != nil {
				return err
			}

			log.Printf("创建标签: %s (%s)", tag.TagName, tag.Category)
		} else {
			log.Printf("标签已存在: %s (%s)", tagData.TagName, tagData.Category)
		}
	}

	log.Println("标签数据生成完成")
	return nil
}

// 生成权限类型数据
func generatePermissions(db *gorm.DB) error {
	log.Println("开始生成权限类型数据...")

	for _, permData := range predefinedPermissions {
		// 检查权限是否已存在
		var count int64
		db.Model(&models.DPermission{}).Where("permission_name = ?", permData.PermissionName).Count(&count)

		if count == 0 {
			perm := models.DPermission{
				Id:             uuid.New(),
				PermissionName: permData.PermissionName,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			if err := db.Create(&perm).Error; err != nil {
				return err
			}

			log.Printf("创建权限类型: %s", perm.PermissionName)
		} else {
			log.Printf("权限类型已存在: %s", permData.PermissionName)
		}
	}

	log.Println("权限类型数据生成完成")
	return nil
}
