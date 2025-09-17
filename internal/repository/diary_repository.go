package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

// DiaryRepository 日记仓库接口
type DiaryRepository interface {
	// GetDiaries 获取所有日记，支持分页、按标题、内容和标签关键字搜索
	GetDiaries(page, pageSize int, keyword string, tagIds []uuid.UUID, permissionId *uuid.UUID, userId uuid.UUID) ([]models.Diary, int64, error)

	// GetDiaryById 根据ID获取日记
	GetDiaryById(diaryId uuid.UUID, userId uuid.UUID) (*models.Diary, error)

	// GetDiaryWithDetails 获取日记详情，包括标签、权限、图片和视频
	GetDiaryWithDetails(diaryId uuid.UUID, userId uuid.UUID) (*models.Diary, []models.Tag, *models.DPermission, []models.DiaryImage, []models.DiaryVideo, error)

	// CreateDiary 创建日记
	CreateDiary(diary *models.Diary, userId uuid.UUID, permissionId uuid.UUID, tagIds []uuid.UUID, imageUrls []string, videoUrls []string) (*models.Diary, error)

	// UpdateDiary 更新日记
	UpdateDiary(diaryId uuid.UUID, userId uuid.UUID, updateData map[string]interface{}, permissionId *uuid.UUID, tagIds []uuid.UUID, imageUrls []string, videoUrls []string) error

	// DeleteDiary 删除日记
	DeleteDiary(diaryId uuid.UUID, userId uuid.UUID) error

	// AddLike 用户点赞日记
	AddLike(diaryId uuid.UUID, userId uuid.UUID) error

	// RemoveLike 取消点赞
	RemoveLike(diaryId uuid.UUID, userId uuid.UUID) error

	// CheckUserLike 检查用户是否点赞
	CheckUserLike(diaryId uuid.UUID, userId uuid.UUID) (bool, error)

	// ShareDiary 分享日记给其他用户
	ShareDiary(diaryId uuid.UUID, shareUserId uuid.UUID, currentUserId uuid.UUID) error
}

// diaryRepository 日记仓库实现
type diaryRepository struct {
	db *gorm.DB
}

// NewDiaryRepository 创建一个新的日记仓库
func NewDiaryRepository() DiaryRepository {
	return &diaryRepository{
		db: database.GetDB(),
	}
}

// GetDiaries 获取所有日记，支持分页、按标题、内容和标签关键字搜索
func (r *diaryRepository) GetDiaries(page, pageSize int, keyword string, tagIds []uuid.UUID, permissionId *uuid.UUID, userId uuid.UUID) ([]models.Diary, int64, error) {
	var diaries []models.Diary
	var total int64

	query := r.db.Model(&models.Diary{})

	// 关联日记用户表，确保只返回用户有权限查看的日记
	query = query.Joins("JOIN diary_users ON diaries.id = diary_users.diary_id AND diary_users.user_id = ?", userId)

	// 如果提供了权限ID，则过滤指定权限的日记
	if permissionId != nil {
		query = query.Joins("JOIN diary_dpermissions ON diaries.id = diary_dpermissions.diary_id AND diary_dpermissions.dpermission_id = ?", permissionId)
	}

	// 关键字搜索
	if keyword != "" {
		query = query.Where("diaries.title LIKE ? OR diaries.content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 标签过滤
	if len(tagIds) > 0 {
		query = query.Joins("JOIN diary_tags ON diaries.id = diary_tags.diary_id").
			Where("diary_tags.tag_id IN ?", tagIds).
			Group("diaries.id") // 去重
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize).Order("diaries.created_at DESC")

	// 执行查询
	if err := query.Find(&diaries).Error; err != nil {
		return nil, 0, err
	}

	return diaries, total, nil
}

// GetDiaryById 根据ID获取日记
func (r *diaryRepository) GetDiaryById(diaryId uuid.UUID, userId uuid.UUID) (*models.Diary, error) {
	var diary models.Diary

	// 获取日记并检查用户是否有权限查看
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 查询日记
		if err := tx.First(&diary, "id = ?", diaryId).Error; err != nil {
			return err
		}

		// 检查用户是否有权限查看该日记
		var diaryUser models.DiaryUser
		if err := tx.First(&diaryUser, "diary_id = ? AND user_id = ?", diaryId, userId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 检查日记是否是公开的
				var diaryPermission models.DiaryDPermission
				if err := tx.Joins("JOIN dpermissions ON diary_dpermissions.dpermission_id = dpermissions.id").
					Where("diary_dpermissions.diary_id = ? AND dpermissions.permission_name = ?", diaryId, "公开").
					First(&diaryPermission).Error; err != nil {
					return errors.New("无权访问该日记")
				}
			} else {
				return err
			}
		}

		// 增加阅读量
		if err := tx.Model(&diary).UpdateColumn("pageview", gorm.Expr("pageview + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &diary, nil
}

// GetDiaryWithDetails 获取日记详情，包括标签、权限、图片和视频
func (r *diaryRepository) GetDiaryWithDetails(diaryId uuid.UUID, userId uuid.UUID) (*models.Diary, []models.Tag, *models.DPermission, []models.DiaryImage, []models.DiaryVideo, error) {
	// 获取日记
	diary, err := r.GetDiaryById(diaryId, userId)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// 获取标签
	var tags []models.Tag
	if err := r.db.Joins("JOIN diary_tags ON tags.id = diary_tags.tag_id").
		Where("diary_tags.diary_id = ?", diaryId).Find(&tags).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// 获取权限
	var permission models.DPermission
	if err := r.db.Joins("JOIN diary_dpermissions ON dpermissions.id = diary_dpermissions.dpermission_id").
		Where("diary_dpermissions.diary_id = ?", diaryId).First(&permission).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// 获取图片
	var images []models.DiaryImage
	if err := r.db.Where("diary_id = ?", diaryId).Find(&images).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// 获取视频
	var videos []models.DiaryVideo
	if err := r.db.Where("diary_id = ?", diaryId).Find(&videos).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return diary, tags, &permission, images, videos, nil
}

// CreateDiary 创建日记
func (r *diaryRepository) CreateDiary(diary *models.Diary, userId uuid.UUID, permissionId uuid.UUID, tagIds []uuid.UUID, imageUrls []string, videoUrls []string) (*models.Diary, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 创建日记
		if err := tx.Create(diary).Error; err != nil {
			return err
		}

		// 创建日记用户关联（记录创建者）
		diaryUser := models.DiaryUser{
			DiaryId: diary.Id,
			UserId:  userId,
		}
		if err := tx.Create(&diaryUser).Error; err != nil {
			return err
		}

		// 创建日记权限关联
		diaryPermission := models.DiaryDPermission{
			DiaryId:       diary.Id,
			DPermissionId: permissionId,
		}
		if err := tx.Create(&diaryPermission).Error; err != nil {
			return err
		}

		// 创建日记标签关联
		for _, tagId := range tagIds {
			diaryTag := models.DiaryTag{
				DiaryId: diary.Id,
				TagId:   tagId,
			}
			if err := tx.Create(&diaryTag).Error; err != nil {
				return err
			}
		}

		// 创建日记图片关联
		for _, imageUrl := range imageUrls {
			diaryImage := models.DiaryImage{
				DiaryId:  diary.Id,
				ImageUrl: imageUrl,
			}
			if err := tx.Create(&diaryImage).Error; err != nil {
				return err
			}
		}

		// 创建日记视频关联
		for _, videoUrl := range videoUrls {
			diaryVideo := models.DiaryVideo{
				DiaryId:  diary.Id,
				VideoUrl: videoUrl,
			}
			if err := tx.Create(&diaryVideo).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return diary, nil
}

// UpdateDiary 更新日记
func (r *diaryRepository) UpdateDiary(diaryId uuid.UUID, userId uuid.UUID, updateData map[string]interface{}, permissionId *uuid.UUID, tagIds []uuid.UUID, imageUrls []string, videoUrls []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否有权限更新该日记
		var diaryUser models.DiaryUser
		if err := tx.First(&diaryUser, "diary_id = ? AND user_id = ?", diaryId, userId).Error; err != nil {
			return errors.New("无权更新该日记")
		}

		// 更新日记基本信息
		if err := tx.Model(&models.Diary{}).Where("id = ?", diaryId).Updates(updateData).Error; err != nil {
			return err
		}

		// 如果提供了新的权限ID，则更新权限
		if permissionId != nil {
			if err := tx.Model(&models.DiaryDPermission{}).Where("diary_id = ?", diaryId).
				Update("dpermission_id", permissionId).Error; err != nil {
				return err
			}
		}

		// 如果提供了新的标签IDs，则更新标签
		if len(tagIds) > 0 {
			// 删除旧标签
			if err := tx.Where("diary_id = ?", diaryId).Delete(&models.DiaryTag{}).Error; err != nil {
				return err
			}

			// 创建新标签关联
			for _, tagId := range tagIds {
				diaryTag := models.DiaryTag{
					DiaryId: diaryId,
					TagId:   tagId,
				}
				if err := tx.Create(&diaryTag).Error; err != nil {
					return err
				}
			}
		}

		// 如果提供了新的图片URLs，则更新图片
		if len(imageUrls) > 0 {
			// 删除旧图片
			if err := tx.Where("diary_id = ?", diaryId).Delete(&models.DiaryImage{}).Error; err != nil {
				return err
			}

			// 创建新图片关联
			for _, imageUrl := range imageUrls {
				diaryImage := models.DiaryImage{
					DiaryId:  diaryId,
					ImageUrl: imageUrl,
				}
				if err := tx.Create(&diaryImage).Error; err != nil {
					return err
				}
			}
		}

		// 如果提供了新的视频URLs，则更新视频
		if len(videoUrls) > 0 {
			// 删除旧视频
			if err := tx.Where("diary_id = ?", diaryId).Delete(&models.DiaryVideo{}).Error; err != nil {
				return err
			}

			// 创建新视频关联
			for _, videoUrl := range videoUrls {
				diaryVideo := models.DiaryVideo{
					DiaryId:  diaryId,
					VideoUrl: videoUrl,
				}
				if err := tx.Create(&diaryVideo).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// DeleteDiary 删除日记（只有创建者可以删除）
func (r *diaryRepository) DeleteDiary(diaryId uuid.UUID, userId uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 检查是否是创建者（查找最早创建的diary_user记录）
		var firstDiaryUser models.DiaryUser
		if err := tx.Where("diary_id = ?", diaryId).Order("created_at asc").First(&firstDiaryUser).Error; err != nil {
			return err
		}

		// 如果不是创建者，则拒绝删除
		if firstDiaryUser.UserId != userId {
			return errors.New("只有创建者可以删除日记")
		}

		// 删除日记（关联表的记录会通过外键级联删除）
		if err := tx.Delete(&models.Diary{}, "id = ?", diaryId).Error; err != nil {
			return err
		}

		return nil
	})
}

// CheckUserLike 检查用户是否点赞
func (r *diaryRepository) CheckUserLike(diaryId uuid.UUID, userId uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.DiaryLike{}).Where("diary_id = ? AND user_id = ?", diaryId, userId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// AddLike 用户点赞日记
func (r *diaryRepository) AddLike(diaryId uuid.UUID, userId uuid.UUID) error {
	// 在事务中执行
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 检查日记是否存在
		var diary models.Diary
		if err := tx.First(&diary, "id = ?", diaryId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("日记不存在")
			}
			return err
		}

		// 检查用户是否已经点赞
		var count int64
		if err := tx.Model(&models.DiaryLike{}).Where("diary_id = ? AND user_id = ?", diaryId, userId).Count(&count).Error; err != nil {
			return err
		}

		// 如果已点赞，则不重复处理
		if count > 0 {
			return errors.New("您已经点赞过该日记")
		}

		// 创建点赞记录
		diaryLike := models.DiaryLike{
			DiaryId:   diaryId,
			UserId:    userId,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&diaryLike).Error; err != nil {
			return err
		}

		// 增加日记点赞数
		if err := tx.Model(&models.Diary{}).Where("id = ?", diaryId).
			UpdateColumn("like", gorm.Expr("like + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
}

// RemoveLike 取消点赞
func (r *diaryRepository) RemoveLike(diaryId uuid.UUID, userId uuid.UUID) error {
	// 在事务中执行
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 检查日记是否存在
		var diary models.Diary
		if err := tx.First(&diary, "id = ?", diaryId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("日记不存在")
			}
			return err
		}

		// 检查用户是否已经点赞
		var diaryLike models.DiaryLike
		result := tx.Where("diary_id = ? AND user_id = ?", diaryId, userId).Delete(&diaryLike)
		if result.Error != nil {
			return result.Error
		}

		// 如果没有找到记录，说明用户没有点赞
		if result.RowsAffected == 0 {
			return errors.New("您尚未点赞该日记")
		}

		// 减少日记点赞数，但确保不会小于0
		if err := tx.Model(&models.Diary{}).Where("id = ? AND `like` > 0", diaryId).
			UpdateColumn("like", gorm.Expr("like - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
}

// ShareDiary 分享日记给其他用户
func (r *diaryRepository) ShareDiary(diaryId uuid.UUID, shareUserId uuid.UUID, currentUserId uuid.UUID) error {
	// 检查当前用户是否有权限分享
	var diaryUser models.DiaryUser
	if err := r.db.First(&diaryUser, "diary_id = ? AND user_id = ?", diaryId, currentUserId).Error; err != nil {
		return errors.New("您没有权限分享此日记")
	}

	// 检查被分享用户是否已经有权限
	var existingShare models.DiaryUser
	err := r.db.First(&existingShare, "diary_id = ? AND user_id = ?", diaryId, shareUserId).Error
	if err == nil {
		return errors.New("该用户已经有此日记的权限")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建新的分享记录
	newDiaryUser := models.DiaryUser{
		DiaryId: diaryId,
		UserId:  shareUserId,
	}
	return r.db.Create(&newDiaryUser).Error
}
