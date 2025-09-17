package scripts

import (
	"fmt"
	"log"
)

// ExampleUsage 展示如何使用数据生成器
func ExampleUsage() {
	// 示例1: 仅生成通用数据（标签、权限等）
	fmt.Println("\n=== 示例1: 生成通用数据 ===")
	if err := GenerateCommonData(); err != nil {
		log.Fatalf("生成通用数据失败: %v", err)
	}

	// 示例2: 使用默认配置为用户生成数据
	fmt.Println("\n=== 示例2: 使用默认配置为用户生成数据 ===")
	userID := "02b67436-1ec9-4635-94cf-50f61eaba009" // 这里替换为实际用户的UUID
	if err := GenerateData(userID, nil); err != nil {
		log.Fatalf("生成数据失败: %v", err)
	}

	// 示例3: 使用自定义配置为用户生成数据
	fmt.Println("\n=== 示例3: 使用自定义配置为用户生成数据 ===")
	customConfig := UserDataConfig{
		DiaryCount:       10,   // 生成10条日记
		AccountBookCount: 3,    // 生成3个账本
		BillsPerBook:     15,   // 每个账本15条账单
		WithImages:       true, // 包含图片
		WithVideos:       true, // 包含视频
	}
	if err := GenerateData(userID, &customConfig); err != nil {
		log.Fatalf("使用自定义配置生成数据失败: %v", err)
	}
}

/*
使用命令行工具生成数据的示例:

1. 生成通用数据:
   go run scripts/cmd/generate_data/main.go -common-only

2. 使用默认配置为用户生成数据:
   go run scripts/cmd/generate_data/main.go -user 02b67436-1ec9-4635-94cf-50f61eaba009

3. 使用自定义配置为用户生成数据:
   go run scripts/cmd/generate_data/main.go -user 02b67436-1ec9-4635-94cf-50f61eaba009 -diaries 10 -books 3 -bills 15
*/
