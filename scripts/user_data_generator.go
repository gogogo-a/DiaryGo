package scripts

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/pkg/database"
	"gorm.io/gorm"
)

// 随机内容生成所需的数据
var (
	diaryTitles = []string{
		"美好的一天", "工作笔记", "旅行记录", "读书感悟",
		"健身日志", "美食探店", "电影评论", "心情随笔",
		"家庭聚会", "学习笔记", "生活点滴", "周末回顾",
	}

	diaryContents = []string{
		"今天天气很好，心情也很愉快。早上起床后，我决定去公园散步，呼吸新鲜空气。在公园里遇到了几个朋友，我们一起聊天，分享近况。",
		"工作会议上，讨论了项目进度。团队成员各自汇报了工作情况，我们遇到了一些技术难题，但通过讨论找到了解决方案。下周将进入项目冲刺阶段。",
		"这次旅行给我留下了深刻的印象。当地的风景如画，人们热情友好。我参观了著名景点，品尝了当地美食，还学习了一些当地文化习俗。",
		"这本书讲述了一个关于成长与挑战的故事。主人公经历了许多困难，但最终找到了自己的人生方向。书中有许多值得思考的观点。",
		"今天的健身计划包括30分钟有氧运动和45分钟力量训练。感觉状态不错，比上周有进步。坚持健身已经一个月了，身体素质有明显提升。",
		"今天尝试了一家新开的餐厅。环境优雅，服务态度很好。点了几道招牌菜，味道确实不错。尤其推荐他们的招牌甜点，口感细腻。",
		"这部电影的剧情紧凑，演员表演出色。导演通过精心的镜头语言讲述了一个动人的故事。音乐配乐也很到位，为电影增添了不少情感。",
		"今天心情有些低落，可能是因为工作压力和天气的原因。决定给自己一些放松的时间，听了喜欢的音乐，做了一些简单的冥想，心情好多了。",
		"周末和家人聚在一起，我们一起做了饭，看了电影。和家人在一起的时光总是那么宝贵，让我感到温暖和幸福。",
		"今天学习了一个新的编程概念，虽然一开始有些困难，但通过实践和查阅资料，最终理解了它的原理和应用场景。",
		"生活中的小确幸：清晨的第一缕阳光，路边盛开的花朵，偶遇的小动物，以及朋友的一条问候信息。",
		"这个周末过得充实而愉快。周六参加了朋友的聚会，周日在家休息并规划了下周的工作计划。",
	}

	accountBookNames = []string{
		"个人账本", "家庭开支", "旅行预算", "学习投资",
		"饮食账本", "娱乐消费", "健康支出", "交通记录",
	}



	billRemarks = []string{
		"超市购物", "午餐", "晚餐", "公交车费",
		"电影票", "健身房会费", "书籍", "衣服",
		"水果", "医药费", "网购", "手机充值",
	}

	imageUrls = []string{
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQFcC8nduaP9ZUI4x_rEQ7QUw-TugycJFudKg&s",
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRyAYHvvVteq1nQELSx5bzbjSUyrem6tIqQCA&s",
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQS2mBibPhxp0MiqmmqHX3YAjqn_phgcR1eEA&s",
		"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQSGemYo8ReW_v7jbwf9L2Xi94qh5YVVkJaJg&s",
		"https://i.pinimg.com/474x/a4/96/00/a496008c58815f0e1677ba24013f19bb.jpg",
	}

	videoUrls = []string{
		"https://example.com/video1.mp4",
		"https://example.com/video2.mp4",
		"https://example.com/video3.mp4",
	}
)

// UserDataConfig 用户数据生成配置
type UserDataConfig struct {
	DiaryCount       int  // 要生成的日记数量
	AccountBookCount int  // 要生成的账本数量
	BillsPerBook     int  // 每个账本生成的账单数量
	WithImages       bool // 是否为日记生成图片关联
	WithVideos       bool // 是否为日记生成视频关联
}

// DefaultUserDataConfig 默认配置
func DefaultUserDataConfig() UserDataConfig {
	return UserDataConfig{
		DiaryCount:       5,
		AccountBookCount: 2,
		BillsPerBook:     10,
		WithImages:       true,
		WithVideos:       true,
	}
}

// GenerateUserData 为指定用户生成数据
func GenerateUserData(userId uuid.UUID, config UserDataConfig) error {
	db := database.GetDB()

	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 获取所有标签
	var diaryTags []models.Tag
	if err := db.Where("category = ?", "diary").Find(&diaryTags).Error; err != nil {
		return fmt.Errorf("获取日记标签失败: %w", err)
	}

	var billTags []models.Tag
	if err := db.Where("category = ?", "bill").Find(&billTags).Error; err != nil {
		return fmt.Errorf("获取账单标签失败: %w", err)
	}

	// 获取所有权限类型
	var permissions []models.DPermission
	if err := db.Find(&permissions).Error; err != nil {
		return fmt.Errorf("获取权限类型失败: %w", err)
	}

	// 生成日记数据
	if err := generateDiaries(db, userId, config.DiaryCount, diaryTags, permissions, config.WithImages, config.WithVideos); err != nil {
		return fmt.Errorf("生成日记数据失败: %w", err)
	}

	// 生成账本和账单数据
	if err := generateAccountBooksAndBills(db, userId, config.AccountBookCount, config.BillsPerBook, billTags); err != nil {
		return fmt.Errorf("生成账本和账单数据失败: %w", err)
	}

	log.Printf("用户 %s 的数据生成完成", userId)
	return nil
}

// 生成日记数据
func generateDiaries(db *gorm.DB, userId uuid.UUID, count int, tags []models.Tag, permissions []models.DPermission, withImages, withVideos bool) error {
	log.Printf("开始为用户 %s 生成 %d 条日记...", userId, count)

	for i := 0; i < count; i++ {
		// 创建日记
		diary := models.Diary{
			Id:        uuid.New(),
			Title:     getRandomItem(diaryTitles),
			Content:   getRandomItem(diaryContents),
			Like:      rand.Intn(10),
			CreatedAt: randomTime(30),
			UpdatedAt: time.Now(),
		}

		// 使用事务确保数据一致性
		err := db.Transaction(func(tx *gorm.DB) error {
			// 创建日记
			if err := tx.Create(&diary).Error; err != nil {
				return err
			}

			// 创建日记-用户关联
			diaryUser := models.DiaryUser{
				Id:        uuid.New(),
				DiaryId:   diary.Id,
				UserId:    userId,
				CreatedAt: diary.CreatedAt,
				UpdatedAt: diary.UpdatedAt,
			}
			if err := tx.Create(&diaryUser).Error; err != nil {
				return err
			}

			// 随机选择1-3个标签关联到日记
			tagCount := rand.Intn(3) + 1
			usedTagIds := make(map[uuid.UUID]bool)

			for j := 0; j < tagCount && j < len(tags); j++ {
				tagIndex := rand.Intn(len(tags))
				tagId := tags[tagIndex].Id

				// 确保标签不重复
				if usedTagIds[tagId] {
					continue
				}
				usedTagIds[tagId] = true

				diaryTag := models.DiaryTag{
					DiaryId: diary.Id,
					TagId:   tagId,
				}
				if err := tx.Create(&diaryTag).Error; err != nil {
					return err
				}
			}

			// 随机选择一个权限类型
			if len(permissions) > 0 {
				permIndex := rand.Intn(len(permissions))
				permId := permissions[permIndex].Id

				diaryPerm := models.DiaryDPermission{
					DiaryId:       diary.Id,
					DPermissionId: permId,
				}
				if err := tx.Create(&diaryPerm).Error; err != nil {
					return err
				}
			}

			// 随机添加图片
			if withImages && rand.Intn(2) == 0 { // 50%概率添加图片
				imageCount := rand.Intn(3) + 1
				for k := 0; k < imageCount; k++ {
					diaryImage := models.DiaryImage{
						Id:       uuid.New(),
						DiaryId:  diary.Id,
						ImageUrl: getRandomItem(imageUrls),
					}
					if err := tx.Create(&diaryImage).Error; err != nil {
						return err
					}
				}
			}

			// 随机添加视频
			if withVideos && rand.Intn(4) == 0 { // 25%概率添加视频
				diaryVideo := models.DiaryVideo{
					Id:       uuid.New(),
					DiaryId:  diary.Id,
					VideoUrl: getRandomItem(videoUrls),
				}
				if err := tx.Create(&diaryVideo).Error; err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}

		log.Printf("创建日记: %s", diary.Title)
	}

	log.Printf("用户 %s 的日记数据生成完成", userId)
	return nil
}

// 生成账本和账单数据
func generateAccountBooksAndBills(db *gorm.DB, userId uuid.UUID, bookCount, billsPerBook int, tags []models.Tag) error {
	log.Printf("开始为用户 %s 生成 %d 个账本...", userId, bookCount)

	for i := 0; i < bookCount; i++ {
		// 创建账本
		accountBook := models.AccountBook{
			Id:        uuid.New(),
			Name:      getRandomItem(accountBookNames),
			CreatedAt: randomTime(60),
			UpdatedAt: time.Now(),
		}

		// 使用事务确保数据一致性
		err := db.Transaction(func(tx *gorm.DB) error {
			// 创建账本
			if err := tx.Create(&accountBook).Error; err != nil {
				return err
			}

			// 创建账本-用户关联
			accountBookUser := models.AccountBookUser{
				Id:            uuid.New(),
				AccountBookId: accountBook.Id,
				UserId:        userId,
				CreatedAt:     accountBook.CreatedAt,
				UpdatedAt:     accountBook.UpdatedAt,
			}
			if err := tx.Create(&accountBookUser).Error; err != nil {
				return err
			}

			// 为账本创建账单
			for j := 0; j < billsPerBook; j++ {
				// 随机决定账单类型
				billType := "expense"
				if rand.Intn(4) == 0 { // 25%概率为收入
					billType = "income"
				}

				// 随机金额，支出一般小于收入
				var amount float64
				if billType == "expense" {
					amount = float64(rand.Intn(1000) + 10) // 10-1010
				} else {
					amount = float64(rand.Intn(5000) + 1000) // 1000-6000
				}

				// 创建账单
				bill := models.Bill{
					Id:            uuid.New(),
					AccountBookId: accountBook.Id,
					UserId:        userId,
					Amount:        amount,
					Type:          billType,
					Remark:        getRandomItem(billRemarks),
					ImageUrl:      "", // 大多数账单不需要图片
					CreatedAt:     randomTimeBetween(accountBook.CreatedAt, time.Now()),
					UpdatedAt:     time.Now(),
				}

				// 随机添加图片URL
				if rand.Intn(10) == 0 { // 10%概率添加图片
					bill.ImageUrl = getRandomItem(imageUrls)
				}

				// 创建账单
				if err := tx.Create(&bill).Error; err != nil {
					return err
				}

				// 根据账单类型选择标签
				var filteredTags []models.Tag
				for _, tag := range tags {
					if (billType == "income" && tag.Type == "收入") ||
						(billType == "expense" && tag.Type == "支出") {
						filteredTags = append(filteredTags, tag)
					}
				}

				// 如果找到了匹配的标签
				if len(filteredTags) > 0 {
					// 随机选择1-2个标签关联到账单
					tagCount := rand.Intn(2) + 1
					for k := 0; k < tagCount && k < len(filteredTags); k++ {
						tagIndex := rand.Intn(len(filteredTags))

						billTag := models.BillTag{
							BillId: bill.Id,
							TagId:  filteredTags[tagIndex].Id,
						}
						if err := tx.Create(&billTag).Error; err != nil {
							return err
						}
					}
				}
			}

			return nil
		})

		if err != nil {
			return err
		}

		log.Printf("创建账本: %s 包含 %d 条账单", accountBook.Name, billsPerBook)
	}

	log.Printf("用户 %s 的账本和账单数据生成完成", userId)
	return nil
}

// 辅助函数: 从切片中随机获取一个元素
func getRandomItem(items []string) string {
	return items[rand.Intn(len(items))]
}

// 辅助函数: 生成过去n天内的随机时间
func randomTime(maxDaysAgo int) time.Time {
	daysAgo := rand.Intn(maxDaysAgo)
	return time.Now().AddDate(0, 0, -daysAgo)
}

// 辅助函数: 生成两个时间点之间的随机时间
func randomTimeBetween(start, end time.Time) time.Time {
	delta := end.Sub(start)
	randDelta := time.Duration(rand.Int63n(int64(delta)))
	return start.Add(randDelta)
}
