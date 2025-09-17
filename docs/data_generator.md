# DiaryGo 数据生成器使用指南

## 1. 概述

DiaryGo数据生成器是一个专为开发和测试环境设计的工具，用于快速生成测试数据。它包括两个主要部分：

1. **通用数据生成器**：创建与用户无关的基础数据，如标签和权限
2. **用户数据生成器**：为指定用户创建个性化数据，包括日记、账本和账单

该工具可以帮助开发人员快速填充数据库，便于功能测试和演示，也可在新用户注册时自动生成示例数据。

## 2. 功能特点

### 2.1 通用数据生成

- 生成标签数据
  - 账单标签（收入/支出类别）
  - 日记标签（生活、工作、旅行等类别）
- 生成权限类型
  - 私密（仅创建者可见）
  - 公开（所有人可见）
  - 共享-只读（特定用户可查看）
  - 共享-可编辑（特定用户可编辑）

### 2.2 用户数据生成

- **日记数据**
  - 随机标题和内容
  - 关联多个标签
  - 设置权限
  - 可选关联图片和视频
  - 随机生成创建时间
- **账本数据**
  - 随机账本名称
  - 创建账本-用户关联
- **账单数据**
  - 收入和支出类型
  - 合理的金额范围（收入通常大于支出）
  - 关联对应类型的标签
  - 随机生成交易时间
  - 可选附加图片

## 3. 使用方法

### 3.1 通过命令行使用

数据生成器提供了命令行工具，可以方便地从终端执行：

```bash
# 1. 仅生成通用数据（标签、权限）
go run scripts/cmd/generate_data/main.go -common-only

# 2. 为特定用户生成数据（使用默认配置）
go run scripts/cmd/generate_data/main.go -user 02b67436-1ec9-4635-94cf-50f61eaba009

# 3. 为特定用户生成自定义数量的数据
go run scripts/cmd/generate_data/main.go -user 02b67436-1ec9-4635-94cf-50f61eaba009 -diaries 10 -books 3 -bills 15

# 4. 生成数据但不包含图片和视频
go run scripts/cmd/generate_data/main.go -user 02b67436-1ec9-4635-94cf-50f61eaba009 -no-images -no-videos
```

### 3.2 命令行参数说明

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-user` | string | "" | 用户ID（UUID格式），指定要为其生成数据的用户 |
| `-common-only` | bool | false | 仅生成通用数据，不生成用户数据 |
| `-diaries` | int | 5 | 要生成的日记数量 |
| `-books` | int | 2 | 要生成的账本数量 |
| `-bills` | int | 10 | 每个账本要生成的账单数量 |
| `-no-images` | bool | false | 设置为true时不为日记生成图片关联 |
| `-no-videos` | bool | false | 设置为true时不为日记生成视频关联 |

### 3.3 在代码中使用

数据生成器也可以直接在代码中调用，集成到应用程序的其他功能中：

```go
import (
    "github.com/google/uuid"
    "github.com/haogeng/DiaryGo/scripts"
)

func main() {
    // 1. 生成通用数据
    scripts.GenerateCommonData()

    // 2. 为用户生成数据（使用默认配置）
    userID, _ := uuid.Parse("02b67436-1ec9-4635-94cf-50f61eaba009")
    scripts.GenerateUserData(userID, scripts.DefaultUserDataConfig())

    // 3. 使用自定义配置
    customConfig := scripts.UserDataConfig{
        DiaryCount:      10,
        AccountBookCount: 3,
        BillsPerBook:    15,
        WithImages:      true,
        WithVideos:      false,
    }
    scripts.GenerateUserData(userID, customConfig)
}
```

## 4. 配置选项

### 4.1 用户数据配置

`UserDataConfig` 结构体提供了以下可配置选项：

```go
type UserDataConfig struct {
    DiaryCount      int  // 要生成的日记数量
    AccountBookCount int  // 要生成的账本数量
    BillsPerBook    int  // 每个账本生成的账单数量
    WithImages      bool // 是否为日记生成图片关联
    WithVideos      bool // 是否为日记生成视频关联
}
```

默认配置为：
- 5条日记
- 2个账本
- 每个账本10条账单
- 包含图片和视频关联

## 5. 数据生成规则

### 5.1 日记生成规则

- 标题从预定义的12个标题中随机选择
- 内容从预定义的12个内容段落中随机选择
- 点赞数在0-9之间随机
- 创建时间在过去30天内随机
- 每条日记随机关联1-3个日记类型标签
- 每条日记随机分配一个权限类型
- 50%概率添加1-3张图片
- 25%概率添加一个视频

### 5.2 账本生成规则

- 名称从预定义的8个名称中随机选择
- 创建时间在过去60天内随机

### 5.3 账单生成规则

- 75%概率为支出，25%概率为收入
- 支出金额范围：10-1010元
- 收入金额范围：1000-6000元
- 备注从预定义的12个备注中随机选择
- 10%概率包含图片
- 根据账单类型（收入/支出）关联相应的标签
- 每条账单随机关联1-2个标签

## 6. 示例输出

执行数据生成后，您可以看到类似以下的输出：

```
开始生成通用数据...
开始生成标签数据...
创建标签: 工资 (bill)
创建标签: 奖金 (bill)
创建标签: 兼职 (bill)
...
标签数据生成完成
开始生成权限类型数据...
创建权限类型: private
创建权限类型: public
创建权限类型: shared_read
创建权限类型: shared_edit
权限类型数据生成完成
公共数据生成完成
开始为用户 02b67436-1ec9-4635-94cf-50f61eaba009 生成数据...
开始为用户 02b67436-1ec9-4635-94cf-50f61eaba009 生成 5 条日记...
创建日记: 美好的一天
创建日记: 工作笔记
创建日记: 旅行记录
创建日记: 读书感悟
创建日记: 健身日志
用户 02b67436-1ec9-4635-94cf-50f61eaba009 的日记数据生成完成
开始为用户 02b67436-1ec9-4635-94cf-50f61eaba009 生成 2 个账本...
创建账本: 个人账本 包含 10 条账单
创建账本: 家庭开支 包含 10 条账单
用户 02b67436-1ec9-4635-94cf-50f61eaba009 的账本和账单数据生成完成
用户 02b67436-1ec9-4635-94cf-50f61eaba009 的数据生成完成
所有数据生成完成
```

## 7. 扩展和定制

### 7.1 添加新的随机内容

如果您想添加更多随机内容，只需更新 `user_data_generator.go` 文件中的相应数组：

```go
diaryTitles = []string{
    "美好的一天", "工作笔记", "旅行记录", 
    // 添加更多标题...
}

diaryContents = []string{
    "今天天气很好，心情也很愉快。早上起床后，我决定去公园散步...",
    // 添加更多内容...
}
```

### 7.2 修改数据生成逻辑

数据生成的核心逻辑位于以下函数中：

- `generateDiaries`: 生成日记相关数据
- `generateAccountBooksAndBills`: 生成账本和账单数据

您可以根据项目需求修改这些函数，调整数据生成的规则和逻辑。

## 8. 注意事项

1. **性能考虑**：生成大量数据可能会导致性能问题，特别是在数据库插入操作中。建议在测试环境中使用。

2. **数据库连接**：确保在运行数据生成器之前正确配置数据库连接，可以通过环境变量设置：
   ```
   export MYSQLUSER=root
   export MYSQLPASSWORD=password
   export MYSQLHOST=localhost
   export MYSQLPORT=3306
   export MYSQLDBNAME=diarygo
   ```

3. **预先迁移**：在使用数据生成器之前，请确保已经执行了数据库迁移，创建了所有必要的表结构。

4. **数据重复**：重复执行数据生成可能会创建重复的数据，但对于标签和权限等基础数据，脚本会检查是否已存在，避免重复创建。

5. **随机性**：生成的数据具有随机性，每次运行会产生不同的结果。如果需要可重复的测试数据，需要修改随机数生成逻辑。 