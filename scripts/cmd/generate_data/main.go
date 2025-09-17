package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/haogeng/DiaryGo/pkg/database"
	"github.com/haogeng/DiaryGo/scripts"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Println("警告: 未找到.env文件或加载失败")
    }
	// 定义命令行参数
	userID := flag.String("user", "", "用户ID (UUID格式)，指定要为其生成数据的用户")
	onlyCommonData := flag.Bool("common-only", false, "仅生成通用数据，不生成用户数据")
	diaryCount := flag.Int("diaries", 5, "为用户生成的日记数量")
	accountBookCount := flag.Int("books", 2, "为用户生成的账本数量")
	billsPerBook := flag.Int("bills", 10, "每个账本生成的账单数量")
	noImages := flag.Bool("no-images", false, "不为日记生成图片关联")
	noVideos := flag.Bool("no-videos", false, "不为日记生成视频关联")

	// 解析命令行参数
	flag.Parse()

	// 初始化数据库连接
	_, err := database.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 用户配置
	config := scripts.DefaultUserDataConfig()
	config.DiaryCount = *diaryCount
	config.AccountBookCount = *accountBookCount
	config.BillsPerBook = *billsPerBook
	config.WithImages = !*noImages
	config.WithVideos = !*noVideos

	// 根据参数执行数据生成
	if *onlyCommonData {
		// 仅生成通用数据
		err = scripts.GenerateCommonData()
		if err != nil {
			log.Fatalf("生成通用数据失败: %v", err)
		}
		fmt.Println("通用数据生成完成")
	} else if *userID != "" {
		// 生成所有数据（通用数据 + 用户数据）
		err = scripts.GenerateData(*userID, &config)
		if err != nil {
			log.Fatalf("数据生成失败: %v", err)
		}
		fmt.Printf("已成功为用户 %s 生成数据\n", *userID)
	} else {
		// 未指定用户ID且未指定仅生成通用数据，则显示帮助信息
		fmt.Println("错误: 必须指定用户ID或使用 -common-only 标志")
		fmt.Println("\n使用说明:")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
