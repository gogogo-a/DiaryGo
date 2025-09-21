# DiaryGo 开发规范指南

## 目录

1. [项目架构概述](#项目架构概述)
2. [API接口规范](#API接口规范)
3. [模型(Models)规范](#模型规范)
4. [仓库(Repository)规范](#仓库规范)
5. [数据库迁移规范](#数据库迁移规范)

## 项目架构概述

DiaryGo 采用分层架构设计，主要分为以下几层：

- **API层**：处理HTTP请求和响应，定义在 `/api` 目录下
- **仓库层**：处理数据存取逻辑，定义在 `/internal/repository` 目录下
- **模型层**：定义数据结构和关系，定义在 `/internal/models` 目录下
- **中间件层**：处理认证、日志等横切关注点，定义在 `/internal/middleware` 目录下
- **工具层**：提供通用工具功能，定义在 `/pkg` 目录下
- **迁移工具**：管理数据库结构变更，定义在 `/migrations` 目录下

## API接口规范

### 文件结构

API接口文件应放置在 `/api/v{version}/` 目录下，按功能模块分文件，例如：

```
/api/v1/
  ├── user.go         # 用户相关API
  ├── diary.go        # 日记相关API
  ├── account_book.go # 账本相关API
  ├── bill.go         # 账单相关API
  ├── tag.go          # 标签相关API
  └── routes.go       # 路由注册
```

### 接口定义规范
全部采用restful风格，使用gin框架，使用swagger注释，使用gin-swagger生成swagger文档。

1. **处理器结构体**：每个功能模块定义一个Handler结构体，包含对应的仓库依赖

```go
// TagHandler 标签处理器
type TagHandler struct {
    repo repository.TagRepository
}

// NewTagHandler 创建标签处理器
func NewTagHandler() *TagHandler {
    return &TagHandler{
        repo: repository.NewTagRepository(),
    }
}
```

2. **路由注册**：每个Handler提供一个RegisterRoutes方法注册路由

```go
// RegisterRoutes 注册标签相关路由
func (h *TagHandler) RegisterRoutes(router *gin.RouterGroup) {
    tags := router.Group("/tags")
    {
        tags.POST("", h.CreateTag)          // 创建标签
        tags.GET("", h.GetTags)             // 获取标签列表
        tags.GET("/:id", h.GetTag)          // 获取标签详情
        tags.PUT("/:id", h.UpdateTag)       // 更新标签
        tags.DELETE("/:id", h.DeleteTag)    // 删除标签
    }
}
```

3. **请求/响应结构体**：为复杂请求定义专用的结构体

```go
// TagRequest 标签请求参数
type TagRequest struct {
    TagName  string `json:"tag_name" binding:"required"`
    Type     string `json:"type" binding:"required"`
    Category string `json:"category" binding:"required"`
}
```

4. **统一响应格式**：使用 `pkg/response` 包中定义的响应函数

```go
// 成功响应
response.Success(c, data)
response.SuccessWithMessage(c, "操作成功", data)

// 错误响应
response.ParamError(c, err.Error())
response.NotFound(c, "资源不存在")
response.ServerError(c, "服务器错误")
response.Unauthorized(c, "未授权")
response.Forbidden(c, "禁止访问")
```

5. **API文档注释**：使用Swagger注释格式

```go
// CreateTag 创建标签
// @Summary 创建标签
// @Description 创建新的标签
// @Tags 标签
// @Accept json
// @Produce json
// @Param tag body TagRequest true "标签信息"
// @Success 201 {object} response.Response{data=models.Tag}
// @Failure 400 {object} response.Response
// @Router /api/v1/tags [post]
```

### API响应数据处理规范

1. **关联对象处理**：避免在API响应中包含未加载的空关联对象

```go
// 错误示例：直接返回带有未加载关联的模型
func (h *DiaryImageHandler) AddImage(c *gin.Context) {
    // ... 创建图片逻辑
    image, err := h.repo.AddImage(diaryId, imageUrl)
    // 直接返回会包含空的关联对象
    response.Success(c, image) // 错误：返回了空的diary对象
}

// 正确示例：返回自定义响应对象
func (h *DiaryImageHandler) AddImage(c *gin.Context) {
    // ... 创建图片逻辑
    image, err := h.repo.AddImage(diaryId, imageUrl)
    
    // 创建精简响应对象
    result := map[string]interface{}{
        "id":        image.Id,
        "diary_id":  image.DiaryId,
        "image_url": image.ImageUrl,
    }
    
    response.Success(c, result)
}
```

2. **集合数据处理**：处理关联集合时，创建精简响应

```go
// 获取集合数据时的处理
func (h *DiaryImageHandler) GetImages(c *gin.Context) {
    // ... 获取图片列表逻辑
    images, err := h.repo.GetDiaryImages(diaryId)
    
    // 创建精简响应列表
    var result []map[string]interface{}
    for _, image := range images {
        result = append(result, map[string]interface{}{
            "id":        image.Id,
            "diary_id":  image.DiaryId,
            "image_url": image.ImageUrl,
        })
    }
    
    response.Success(c, result)
}
```

3. **预加载关联**：当确实需要返回关联对象时，确保完整预加载

```go
// 需要返回完整关联对象时
func (r *diaryRepository) GetDiaryWithDetails(diaryId uuid.UUID) (*models.Diary, error) {
    var diary models.Diary
    err := r.db.Preload("Tags").
        Preload("Images").
        Preload("Videos").
        Preload("Permission").
        Where("id = ?", diaryId).
        First(&diary).Error
    
    return &diary, err
}
```

4. **自定义DTO**：对于复杂对象，定义专门的DTO（数据传输对象）

```go
// DiaryDetailDTO 日记详情数据传输对象
type DiaryDetailDTO struct {
    Id          string             `json:"id"`
    Title       string             `json:"title"`
    Content     string             `json:"content"`
    Address     string             `json:"address"`
    Like        int                `json:"like"`
    CreatedAt   time.Time          `json:"created_at"`
    Tags        []TagDTO           `json:"tags"`
    Images      []ImageDTO         `json:"images"`
    Videos      []VideoDTO         `json:"videos"`
    Permission  PermissionDTO      `json:"permission"`
}

// 转换函数
func toDiaryDetailDTO(diary *models.Diary, tags []models.Tag, ...) DiaryDetailDTO {
    // 转换逻辑...
}
```

5. **空值处理**：对可能为空的字段进行合理处理

```go
// 处理可能为nil的指针
func getUserName(user *models.User) string {
    if user == nil {
        return ""
    }
    return user.UserName
}
```

## 模型规范

### 文件结构

模型文件应放置在 `/internal/models/` 目录下，每个模型一个文件：

```
/internal/models/
  ├── user.go
  ├── diary.go
  ├── tag.go
  ├── diary_tag.go     # 多对多关系表
  ├── diary_image.go   # 一对多关系表
  └── ...
```

### 模型定义规范

1. **基本结构**：每个模型都应该定义为一个结构体，包含数据库字段和JSON标签

```go
// Tag 标签模型
type Tag struct {
    Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
    TagName   string    `json:"tag_name" gorm:"type:varchar(255);not null"`
    Type      string    `json:"type" gorm:"type:varchar(255);not null"`
    Category  string    `json:"category" gorm:"type:varchar(255);not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

2. **表名定义**：使用TableName方法指定表名

```go
// TableName 指定表名
func (Tag) TableName() string {
    return "tags"
}
```

3. **UUID主键**：使用BeforeCreate钩子生成UUID主键

```go
// BeforeCreate 创建前钩子
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
    t.Id = uuid.New()
    t.CreatedAt = time.Now()
    t.UpdatedAt = time.Now()
    return nil
}
```

4. **时间戳更新**：使用BeforeUpdate钩子更新UpdatedAt字段

```go
// BeforeUpdate 更新前钩子
func (t *Tag) BeforeUpdate(tx *gorm.DB) error {
    t.UpdatedAt = time.Now()
    return nil
}
```

5. **关联关系定义**：使用GORM标签定义关联关系

```go
// 一对多关系
type Diary struct {
    Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
    // 其他字段
}

type DiaryImage struct {
    Id        uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
    DiaryId   uuid.UUID `json:"diary_id" gorm:"type:char(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    Diary     Diary     `json:"diary" gorm:"foreignKey:DiaryId;references:Id"` // 定义关联
    // 其他字段
}

// 多对多关系
type DiaryTag struct {
    DiaryId   uuid.UUID `json:"diary_id" gorm:"type:char(36);primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    TagId     uuid.UUID `json:"tag_id" gorm:"type:char(36);primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    Diary     Diary     `json:"diary" gorm:"foreignKey:DiaryId;references:Id"`
    Tag       Tag       `json:"tag" gorm:"foreignKey:TagId;references:Id"`
}
```

6. **MySQL UUID兼容性**：UUID字段使用char(36)类型

```go
Id uuid.UUID `json:"id" gorm:"primaryKey;type:char(36)"`
```

## 仓库规范

### 文件结构

仓库文件应放置在 `/internal/repository/` 目录下，每个功能模块一个文件：

```
/internal/repository/
  ├── user_repository.go
  ├── diary_repository.go
  ├── tag_repository.go
  └── ...
```

### 仓库定义规范

1. **接口定义**：先定义接口，再实现接口

```go
// TagRepository 标签仓库接口
type TagRepository interface {
    Create(tag *models.Tag) error
    GetByID(id uuid.UUID) (*models.Tag, error)
    GetAll(category string) ([]models.Tag, error)
    Update(tag *models.Tag) error
    Delete(id uuid.UUID) error
}
```

2. **仓库实现**：定义结构体，实现上述接口

```go
// tagRepository 标签仓库实现
type tagRepository struct {
    db *gorm.DB
}

// NewTagRepository 创建标签仓库
func NewTagRepository() TagRepository {
    return &tagRepository{
        db: database.GetDB(),
    }
}
```

3. **事务处理**：复杂操作使用事务

```go
func (r *tagRepository) BatchCreate(tags []*models.Tag) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        for _, tag := range tags {
            if err := tx.Create(tag).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

4. **错误处理**：使用标准errors包和gorm.ErrRecordNotFound

```go
func (r *tagRepository) GetByID(id uuid.UUID) (*models.Tag, error) {
    var tag models.Tag
    err := r.db.Where("id = ?", id).First(&tag).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("标签不存在")
        }
        return nil, err
    }
    return &tag, nil
}
```

5. **查询构建**：使用链式调用构建复杂查询

```go
func (r *diaryRepository) GetDiaries(query DiaryQuery) ([]models.Diary, int64, error) {
    var diaries []models.Diary
    var total int64
    
    db := r.db.Model(&models.Diary{})
    
    if query.Keyword != "" {
        db = db.Where("title LIKE ? OR content LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
    }
    
    if query.TagID != uuid.Nil {
        db = db.Joins("JOIN diary_tags ON diaries.id = diary_tags.diary_id").
             Where("diary_tags.tag_id = ?", query.TagID)
    }
    
    // 先获取总数
    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // 再获取分页数据
    if err := db.Limit(query.PageSize).Offset((query.Page - 1) * query.PageSize).
        Order("created_at DESC").Find(&diaries).Error; err != nil {
        return nil, 0, err
    }
    
    return diaries, total, nil
}
```

6. **预加载关联**：使用Preload加载关联数据

```go
func (r *diaryRepository) GetDiaryById(id uuid.UUID) (*models.Diary, error) {
    var diary models.Diary
    err := r.db.Preload("Tags").
        Preload("Images").
        Preload("Videos").
        Where("id = ?", id).First(&diary).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("日记不存在")
        }
        return nil, err
    }
    return &diary, nil
}
```

## 数据库迁移规范

### 文件结构

迁移文件应放置在 `/migrations/` 目录下：

```
/migrations/
  └── migrate.go  # 主迁移文件
```

### 迁移规范

1. **模型分组**：按照功能和依赖关系对模型进行分组

```go
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
        // 其他分组...
    }

    // 迁移每个分组的模型
    for _, group := range modelGroups {
        logger.Info("正在迁移%s相关表...", group.name)
        if err := db.AutoMigrate(group.models...); err != nil {
            return fmt.Errorf("迁移%s表失败: %w", group.name, err)
        }
    }
    
    return nil
}
```

2. **迁移顺序**：按照依赖关系顺序迁移

- 先迁移基本模型（如User, Tag）
- 再迁移依赖基本模型的模型（如Diary）
- 最后迁移关系模型（如DiaryTag, DiaryImage）

3. **初始数据**：在迁移后可以插入初始数据

```go
// 可以添加一个函数来插入初始数据
func SeedInitialData(db *gorm.DB) error {
    // 创建默认权限
    permissions := []models.DPermission{
        {Name: "private", Description: "私密，仅创建者可见"},
        {Name: "public", Description: "公开，所有人可见"},
        {Name: "shared", Description: "共享，特定用户可见"},
    }
    
    for _, p := range permissions {
        if err := db.Where("name = ?", p.Name).FirstOrCreate(&p).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

4. **日志记录**：记录迁移过程

```go
logger.Info("开始数据库迁移...")
// 执行迁移...
logger.Info("数据库迁移完成")
```

5. **错误处理**：详细记录迁移错误

```go
if err := db.AutoMigrate(group.models...); err != nil {
    logger.Error("迁移%s表失败: %v", group.name, err)
    return fmt.Errorf("迁移%s表失败: %w", group.name, err)
}
```
