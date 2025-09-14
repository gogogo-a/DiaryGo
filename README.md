# DiaryGo

DiaryGo是一个基于Go语言和Gin框架开发的日记应用后端服务。

## 项目结构

```
DiaryGo/
  - api/             # API定义和版本控制
    - v1/            # API v1版本
    - v2/            # API v2版本
  - config/          # 配置文件
  - internal/        # 内部包
    - handler/       # HTTP处理器
    - middleware/    # 中间件
    - models/        # 数据模型
    - repository/    # 数据仓库
  - main.go          # 程序入口
  - migrations/      # 数据库迁移文件
  - pkg/             # 可重用的公共包
    - database/      # 数据库连接
    - logger/        # 日志工具
    - utils/         # 通用工具函数
  - scripts/         # 脚本文件
```

## 环境要求

- Go 1.18或更高版本
- 支持的操作系统: Linux, macOS, Windows

## 安装与启动

### 1. 克隆仓库

```bash
git clone https://github.com/haogeng/DiaryGo.git
cd DiaryGo
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 启动服务

```bash
go run main.go
```

服务将在 http://localhost:8080 启动。

### 4. 构建二进制文件

```bash
go build -o bin/diary-go
```

构建完成后，可以通过以下命令运行:

```bash
./bin/diary-go
```

## API测试

启动服务后，可以通过以下命令测试API是否正常工作:

```bash
curl http://localhost:8080/ping
```

预期返回:

```json
{"message":"pong"}
```

## 配置

可以通过环境变量或配置文件修改服务配置（待实现）。

## 开发

### 添加新路由

在`main.go`中或适当的处理器文件中添加新的路由:

```go
r.GET("/api/v1/diaries", handler.GetDiaries)
```

### 构建并运行测试

```bash
go test ./...
``` 