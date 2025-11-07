# go-core

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/kamalyes/go-core)
[![license](https://img.shields.io/github/license/kamalyes/go-core)]()
[![download](https://img.shields.io/github/downloads/kamalyes/go-core/total)]()
[![release](https://img.shields.io/github/v/release/kamalyes/go-core)]()
[![commit](https://img.shields.io/github/last-commit/kamalyes/go-core)]()
[![issues](https://img.shields.io/github/issues/kamalyes/go-core)]()
[![pull](https://img.shields.io/github/issues-pr/kamalyes/go-core)]()
[![fork](https://img.shields.io/github/forks/kamalyes/go-core)]()
[![star](https://img.shields.io/github/stars/kamalyes/go-core)]()
[![go](https://img.shields.io/github/go-mod/go-version/kamalyes/go-core)]()
[![size](https://img.shields.io/github/repo-size/kamalyes/go-core)]()
[![contributors](https://img.shields.io/github/contributors/kamalyes/go-core)]()
[![codecov](https://codecov.io/gh/kamalyes/go-core/branch/master/graph/badge.svg)](https://codecov.io/gh/kamalyes/go-core)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamalyes/go-core)](https://goreportcard.com/report/github.com/kamalyes/go-core)
[![Go Reference](https://pkg.go.dev/badge/github.com/kamalyes/go-core?status.svg)](https://pkg.go.dev/github.com/kamalyes/go-core?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/kamalyes/go-core/-/badge.svg)](https://sourcegraph.com/github.com/kamalyes/go-core?badge)

### 介绍

go-core 是 go web 应用开发脚手架，从全局配置文件读取，zap日志组件始化，gorm数据库连接初始化，redis客户端初始化，http server启动等。最终实现简化流程、提高效率、统一规范。

## 📦 核心组件

| 组件 | 功能 | 特性 |
|------|------|------|
| 🗄️ **Database** | 数据库操作层 | 支持MySQL/PostgreSQL/SQLite，查询构建器，事务管理，87.8%测试覆盖率 |
| 🔐 **JWT** | JWT令牌管理 | 令牌生成/验证，黑名单机制，Redis存储支持 |
| 🛡️ **Casbin** | 权限控制 | RBAC权限模型，动态权限管理，多域支持 |
| 📊 **Response** | 统一响应格式 | 多框架适配，标准化错误码，场景化响应 |
| 🚀 **SRun** | 服务启动器 | 优雅关闭，健康检查，多平台支持 |
| 📝 **Zap** | 日志管理 | 结构化日志，性能优化，多输出目标 |
| 💾 **Redis** | 缓存客户端 | 连接池管理，序列化支持，分布式锁 |
| 🔧 **Global** | 全局管理 | 配置管理，依赖注入，组件协调 |
| 🎨 **Captcha** | 验证码服务 | 图形验证码，Redis存储，防暴力破解 |
| 🌐 **MQTT** | 消息队列 | 发布订阅，连接管理，异步处理 |
| ☁️ **OSS** | 对象存储 | MinIO集成，文件上传，存储管理 |
| 🔗 **Consul** | 服务发现 | 微服务注册，健康检查，配置中心 |

## ⚡ 快速开始

### 安装依赖

```bash
go get -u github.com/kamalyes/go-core
```

### 项目结构

```
your-project/
├── resources/                 # 配置文件目录（必需）
│   ├── dev_config.yaml       # 开发环境配置
│   ├── fat_config.yaml       # 测试环境配置  
│   ├── pro_config.yaml       # 生产环境配置
│   └── uat_config.yaml       # 预生产环境配置
├── main.go                   # 主程序入口
└── go.mod                    # 依赖管理
```

> 📝 **注意**: 配置文件格式请参考 [go-config](https://github.com/kamalyes/go-config) 项目

### 基础用法

```go
package main

import (
    "github.com/gin-gonic/gin"
    goconfig "github.com/kamalyes/go-config"
    "github.com/kamalyes/go-core/pkg/database"
    "github.com/kamalyes/go-core/pkg/global" 
    "github.com/kamalyes/go-core/pkg/redis"
    "github.com/kamalyes/go-core/pkg/srun"
    "github.com/kamalyes/go-core/pkg/zap"
)

func main() {
    // 1. 加载配置文件（根据环境变量自动选择）
    global.CONFIG = goconfig.GlobalConfig()

    // 2. 初始化日志系统
    global.LOG = zap.Zap()

    // 3. 初始化数据库连接
    global.DB = database.Gorm()

    // 4. 初始化Redis客户端
    global.REDIS = redis.Redis()

    // 5. 启动Web服务
    r := gin.Default()
    
    // 健康检查端点
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // 启动HTTP服务器（支持优雅关闭）
    srun.RunHttpServer(r)
}
```

## 🔧 高级用法

### 数据库操作示例

```go
package main

import (
    "github.com/kamalyes/go-core/pkg/database"
    "github.com/kamalyes/go-core/pkg/global"
)

type User struct {
    ID       uint   `gorm:"primarykey"`
    Username string `gorm:"uniqueIndex"`
    Email    string
    Age      int
}

func DatabaseExample() {
    // 使用查询构建器进行复杂查询
    builder := database.NewQueryBuilder().
        Where("age > ?", 18).
        Where("status = ?", "active").
        OrderBy("created_at DESC").
        Limit(10)
    
    var users []User
    if err := builder.Find(global.DB, &users); err != nil {
        global.LOG.Error("查询失败", zap.Error(err))
        return
    }
    
    // 使用事务处理
    tx := database.NewTransactionHandler(global.DB)
    if err := tx.WithTransaction(func() error {
        // 在事务中执行多个操作
        user := &User{Username: "john", Email: "john@example.com"}
        return tx.Create(user)
    }); err != nil {
        global.LOG.Error("事务执行失败", zap.Error(err))
    }
}
```

### JWT认证示例

```go
package main

import (
    "github.com/kamalyes/go-core/pkg/jwt"
    "github.com/kamalyes/go-core/pkg/global"
)

func AuthExample() {
    j := jwt.NewJWT()
    
    // 生成Token
    claims := jwt.BaseClaims{
        UUID:     "user-uuid-123",
        ID:       1,
        Username: "john",
        NickName: "John Doe",
    }
    
    token, err := j.CreateToken(claims)
    if err != nil {
        global.LOG.Error("Token生成失败", zap.Error(err))
        return
    }
    
    // 验证Token
    parsedClaims, err := j.ParseToken(token)
    if err != nil {
        global.LOG.Error("Token验证失败", zap.Error(err))
        return
    }
    
    global.LOG.Info("用户认证成功", zap.String("username", parsedClaims.Username))
}
```

### 统一响应处理

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/kamalyes/go-core/pkg/response"
)

func APIExample() {
    r := gin.Default()
    
    // 成功响应
    r.GET("/users", func(c *gin.Context) {
        users := []User{{ID: 1, Username: "john"}}
        response.OkWithData(users, c)
    })
    
    // 错误响应
    r.POST("/users", func(c *gin.Context) {
        // 参数验证失败
        response.FailWithMessage("用户名已存在", c)
    })
    
    // 自定义响应
    r.GET("/custom", func(c *gin.Context) {
        response.Result(response.SUCCESS, "操作成功", gin.H{
            "timestamp": time.Now().Unix(),
        }, c)
    })
}
```

## 🏗️ 架构设计

### 分层架构

```text
┌─────────────────────────────────────────┐
│              应用层 (Your App)            │  
├─────────────────────────────────────────┤
│              go-core 核心层               │
│  ┌─────────┬─────────┬─────────┬──────┐ │
│  │ Response│   JWT   │ Database│ Zap  │ │
│  └─────────┴─────────┴─────────┴──────┘ │
│  ┌─────────┬─────────┬─────────┬──────┐ │
│  │ Casbin  │  Redis  │  MQTT   │ OSS  │ │
│  └─────────┴─────────┴─────────┴──────┘ │
├─────────────────────────────────────────┤
│           基础设施层 (Infrastructure)     │
│     Database │ Cache │ MessageQueue     │
└─────────────────────────────────────────┘
```

### 核心原则

- **🔌 插拔式设计**: 每个组件都可以独立使用
- **📏 标准化**: 统一的接口设计和错误处理
- **⚡ 高性能**: 优化的数据库查询和缓存策略  
- **🛡️ 安全优先**: 内置安全机制和最佳实践
- **📊 可观测**: 完整的日志、监控和调试支持

## 📋 详细文档

### 核心组件文档

| 组件 | 文档链接 | 说明 |
|------|----------|------|
| 🗄️ Database | [database使用指南](pkg/database/database_usage.md) | 数据库操作详细指南，包含查询构建器、事务管理等 |
| 🔐 JWT | [jwt组件](pkg/jwt/) | JWT令牌管理和认证机制 |
| 📊 Response | [response组件](pkg/response/) | 统一响应格式和错误码管理 |
| 🛡️ Casbin | [casbin组件](pkg/casbin/) | RBAC权限控制系统 |

### 配置管理

项目使用 [go-config](https://github.com/kamalyes/go-config) 进行配置管理，支持：

- 🌍 **多环境配置**: dev/fat/uat/pro 环境自动切换
- 🔄 **热更新**: 配置文件变更自动重载
- 🗂️ **格式支持**: YAML/JSON/TOML 等多种格式
- 🔒 **敏感信息**: 环境变量和加密支持

## 🚀 性能特性

### 数据库性能

- ✅ **连接池优化**: 智能连接池管理，最大化资源利用
- ✅ **查询优化**: 预编译语句和查询缓存  
- ✅ **批量操作**: 高效的批量插入和更新
- ✅ **读写分离**: 支持主从数据库配置

### 缓存策略

- ⚡ **多级缓存**: 内存缓存 + Redis 分布式缓存
- 🔄 **智能失效**: TTL 和事件驱动的缓存失效
- 🎯 **缓存穿透防护**: 布隆过滤器和空值缓存

### 并发处理

- 🔒 **分布式锁**: Redis 实现的分布式锁机制
- ⚖️ **负载均衡**: 支持多种负载均衡算法
- 📊 **限流控制**: 令牌桶和滑动窗口限流

## 🛡️ 安全特性

### 认证与授权

- 🎫 **JWT认证**: 无状态令牌认证
- 👥 **RBAC权限**: 基于角色的访问控制
- 🔐 **多因子认证**: 支持验证码等多因子认证

### 数据安全

- 💉 **SQL注入防护**: 参数化查询和输入验证
- 🔒 **数据加密**: 敏感数据自动加密存储
- 📝 **审计日志**: 完整的操作审计追踪

## 🧪 测试覆盖

项目具有完整的测试体系：

- 📊 **覆盖率**: Database模块测试覆盖率达到 87.8%
- 🧪 **测试类型**: 单元测试、集成测试、基准测试
- 🔄 **持续集成**: GitHub Actions 自动化测试

## 🤝 贡献指南

我们欢迎社区贡献！请参考以下步骤：

1. 🍴 Fork 项目仓库
2. 🌿 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 💾 提交更改 (`git commit -m 'Add amazing feature'`)
4. 📤 推送分支 (`git push origin feature/amazing-feature`)
5. 🔄 提交 Pull Request

### 开发规范

- 📝 **代码规范**: 遵循 Go 官方编码规范
- 🧪 **测试要求**: 新功能需要包含相应测试
- 📚 **文档更新**: 功能变更需要更新对应文档
- 🎯 **性能考虑**: 关注性能影响和资源使用

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

感谢以下优秀的开源项目：

- [GORM](https://gorm.io/) - Go ORM 框架
- [Gin](https://gin-gonic.com/) - HTTP Web 框架  
- [Zap](https://github.com/uber-go/zap) - 高性能日志库
- [Redis](https://redis.io/) - 内存数据结构存储
- [Casbin](https://casbin.org/) - 访问控制库

## 📞 联系方式

- 📧 **Email**: [501893067@qq.com](mailto:501893067@qq.com)
- 🐛 **Issues**: [GitHub Issues](https://github.com/kamalyes/go-core/issues)
- 📖 **Documentation**: [项目文档](https://pkg.go.dev/github.com/kamalyes/go-core)

---

### ⭐ 如果这个项目对你有帮助，请给我们一个 Star！⭐
