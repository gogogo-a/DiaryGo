package scripts

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

// GenerateData 生成应用程序所需的所有数据
// 参数:
// - userId: 用户ID，如果为空则只生成通用数据
// - config: 用户数据生成配置，如果为nil则使用默认配置
func GenerateData(userId string, config *UserDataConfig) error {
	// 生成通用数据（与用户无关的数据，如标签、权限等）
	log.Println("开始生成通用数据...")
	if err := GenerateCommonData(); err != nil {
		return fmt.Errorf("生成通用数据失败: %w", err)
	}

	// 如果提供了用户ID，则为该用户生成数据
	if userId != "" {
		log.Printf("开始为用户 %s 生成数据...", userId)

		// 解析用户ID
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			return fmt.Errorf("无效的用户ID格式: %w", err)
		}

		// 使用默认配置或提供的配置
		userConfig := DefaultUserDataConfig()
		if config != nil {
			userConfig = *config
		}

		// 生成用户相关数据
		if err := GenerateUserData(userUUID, userConfig); err != nil {
			return fmt.Errorf("为用户 %s 生成数据失败: %w", userId, err)
		}

		log.Printf("用户 %s 的数据生成完成", userId)
	}

	log.Println("所有数据生成完成")
	return nil
}
