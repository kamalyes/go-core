# Database 包使用指南

## 📋 目录

- [概述](#概述)
- [核心组件](#核心组件)
- [快速开始](#快速开始)
- [基础查询](#基础查询)
- [查询构建器](#查询构建器)
- [高级查询](#高级查询)
- [事务操作](#事务操作)
- [分页查询](#分页查询)
- [数据库兼容性](#数据库兼容性)
- [性能优化](#性能优化)
- [错误处理](#错误处理)
- [最佳实践](#最佳实践)
- [API 参考](#api-参考)
- [配置示例](#配置示例)
- [常见问题](#常见问题)

## 🌟 概述

`database` 包是一个基于 GORM 构建的高级数据库操作层，提供了：

- 🚀 **高性能**: 优化的查询构建和执行
- 🔧 **灵活性**: 支持复杂查询条件和动态查询构建
- 🛡️ **安全性**: 内置 SQL 注入防护和事务安全
- 🔄 **兼容性**: 支持 MySQL、PostgreSQL、SQLite 等多种数据库
- 📊 **可观测性**: 完整的日志记录和性能监控
- 🎯 **易用性**: 链式 API 和直观的接口设计

## 🏗️ 核心组件

### 1. Handler 接口

`Handler` 是数据库操作的核心接口：

```go
type Handler interface {
    // 基础操作
    DB() *gorm.DB                                    // 获取原始 GORM 实例
    Query(param QueryParam) *gorm.DB                 // 执行查询
    Close() error                                    // 关闭数据库连接
    AutoMigrate(dst ...any) error                    // 自动迁移表结构
    
    // 事务操作
    Begin(opts ...*sql.TxOptions) Handler            // 开始事务
    Commit() error                                   // 提交事务
    Rollback() error                                 // 回滚事务
}
```

### 2. QueryParam 接口

定义查询参数的统一接口：

```go
type QueryParam interface {
    Where(db *gorm.DB) *gorm.DB  // 应用查询条件到 GORM 实例
}
```

### 3. 查询参数类型

#### SimpleQueryParam - 简单查询
```go
type SimpleQueryParam struct {
    Query string        // SQL 查询条件
    Args  []interface{} // 查询参数
}
```

#### PageQueryParam - 分页查询
```go 
type PageQueryParam struct {
    Query   string        // SQL 查询条件
    Args    []interface{} // 查询参数
    Limit   int          // 每页数量
    Offset  int          // 偏移量  
    OrderBy string       // 排序字段
}
```

#### AdvancedQueryParam - 高级查询
支持复杂的动态查询构建，包括多种过滤器、时间范围、FIND_IN_SET 等。

## 🚀 快速开始

### 1. 初始化数据库处理器

```go
package main

import (
    "log"
    "github.com/kamalyes/go-core/pkg/database"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    // 方式1: 使用默认配置（推荐）
    handler := database.GetDefaultHandler()
    if handler == nil {
        log.Fatal("Failed to get default database handler")
    }
    defer handler.Close()
    
    // 方式2: 自定义数据库连接
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    handler = database.NewHandler(db)
    defer handler.Close()
}
```

### 2. 定义数据模型

```go
// 用户模型
type User struct {
    ID         uint      `gorm:"primarykey" json:"id"`
    Username   string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
    Email      string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
    Age        int       `json:"age"`
    BusinessID int64     `gorm:"index" json:"business_id"`
    ShopID     int64     `gorm:"index" json:"shop_id"`
    Status     int       `gorm:"default:1;index" json:"status"`
    Tags       string    `gorm:"type:text" json:"tags"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

// 产品模型  
type Product struct {
    ID         uint      `gorm:"primarykey" json:"id"`
    Name       string    `gorm:"size:100;not null" json:"name"`
    Price      float64   `gorm:"type:decimal(10,2)" json:"price"`
    BusinessID int64     `gorm:"index" json:"business_id"`
    ShopID     int64     `gorm:"index" json:"shop_id"`
    CategoryID int       `gorm:"index" json:"category_id"`
    Status     int       `gorm:"default:1;index" json:"status"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

### 3. 自动迁移表结构

```go
// 自动迁移表结构
err := handler.AutoMigrate(&User{}, &Product{})
if err != nil {
    log.Fatal("Failed to migrate database:", err)
}
```

## 🎯 基础查询

### 1. 简单查询

```go
// 使用 SimpleQueryParam 进行基础查询
param := database.NewSimpleQueryParam("age > ? AND status = ?", 18, 1)

var users []User
result := handler.Query(param).Find(&users)
if result.Error != nil {
    log.Printf("Query error: %v", result.Error)
    return
}

fmt.Printf("Found %d users\n", len(users))
```

### 2. 单条记录查询

```go
// 查询单个用户
param := database.NewSimpleQueryParam("username = ?", "john_doe")

var user User
result := handler.Query(param).First(&user)
if result.Error != nil {
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        fmt.Println("User not found")
    } else {
        log.Printf("Query error: %v", result.Error)
    }
    return
}
```

### 3. 计数查询

```go
// 统计符合条件的记录数
param := database.NewSimpleQueryParam("business_id = ? AND status = ?", 1, 1)

var count int64
result := handler.Query(param).Model(&User{}).Count(&count)
if result.Error != nil {
    log.Printf("Count error: %v", result.Error)
    return
}

fmt.Printf("Total active users: %d\n", count)
```

## 🏗️ 查询构建器

### 1. 基础用法

查询构建器提供了链式 API 来构建复杂查询：

```go
// 构建复杂查询
param := database.NewQueryBuilder().
    WithBusinessId(1).                                    // 业务ID过滤
    WithShopId(101).                                     // 店铺ID过滤  
    WhereIn("status", []interface{}{1, 2}).              // IN 查询
    WithOrder("created_at", "DESC").                     // 排序
    WithPagination(10, 0).                               // 分页
    Build()

var users []User
result := handler.Query(param).Find(&users)
if result.Error != nil {
    log.Printf("Query error: %v", result.Error)
}
```

### 2. 支持的查询条件

#### IN 查询

```go
// 状态 IN 查询
param := database.NewQueryBuilder().
    WhereIn("status", []interface{}{1, 2, 3}).
    Build()

// 分类 IN 查询  
param = database.NewQueryBuilder().
    WhereIn("category_id", []interface{}{10, 20, 30}).
    Build()
```

#### 模糊查询

```go
// 用户名模糊查询（左匹配，性能更好）
param := database.NewQueryBuilder().
    WhereLike("username", []interface{}{"john", "jane"}, false). // false: 左匹配 user%
    Build()

// 邮箱全匹配查询  
param = database.NewQueryBuilder().
    WhereLike("email", []interface{}{"@gmail.com"}, true). // true: 全匹配 %@gmail.com%
    Build()
```

#### 时间范围查询

```go
// 创建时间范围查询
param := database.NewQueryBuilder().
    WhereTimeRange("created_at", "2023-01-01 00:00:00", "2023-12-31 23:59:59").
    Build()

// 更新时间范围查询
param = database.NewQueryBuilder().
    WhereTimeRange("updated_at", "2023-06-01", "2023-06-30").
    Build()
```

#### FIND_IN_SET 查询

适用于存储逗号分隔值的字段：

```go
// 标签查询（适用于 tags 字段存储 "vip,active,premium" 格式）
param := database.NewQueryBuilder().
    WhereFindInSet("tags", []string{"vip", "active"}).
    Build()

// 权限查询
param = database.NewQueryBuilder().
    WhereFindInSet("permissions", []string{"read", "write"}).
    Build()
```

### 3. 排序和分页

```go
// 多字段排序
param := database.NewQueryBuilder().
    WithBusinessId(1).
    WithOrder("status", "ASC").                          // 先按状态升序
    WithOrder("created_at", "DESC").                     // 再按创建时间降序
    WithPagination(20, 40).                              // 每页20条，跳过前40条
    Build()

// 分组查询
param = database.NewQueryBuilder().
    WithBusinessId(1).
    WithGroupBy("category_id").                          // 按分类分组
    Build()
```

## 🔧 高级查询

### 1. 使用 AdvancedQueryParam

`AdvancedQueryParam` 提供了最灵活的查询构建方式：

```go
// 创建高级查询选项
option := &database.FindOptionCommon{
    BusinessId:               1,                    // 业务ID
    ShopId:                  101,                   // 店铺ID
    ExcludeBusiness:         false,                 // 是否排除业务ID条件
    ExcludeShop:            false,                 // 是否排除店铺ID条件
    ExcludeBusinessAndShop: false,                 // 是否排除业务和店铺条件
    IncludeBusinessIdZero:  false,                 // 是否包含业务ID为0的情况
    Limit:                  10,                    // 限制数量
    Offset:                 0,                     // 偏移量
    By:                     "created_at",          // 排序字段
    Order:                  "DESC",                // 排序方向
    GroupBy:                "category_id",         // 分组字段
    TablePrefix:            "",                    // 表前缀（单表查询留空）
    DisableOrderBy:         false,                 // 是否禁用排序
}

// 创建高级查询参数
param := database.NewAdvancedQueryParam(option)

// 添加精确匹配过滤器
param.AddFilter(&database.BaseInfoFilter{
    DBField:    "age",
    Values:     []interface{}{25, 30, 35},
    ExactMatch: true,                              // 精确匹配 IN 查询
})

// 添加模糊匹配过滤器
param.AddFilter(&database.BaseInfoFilter{
    DBField:    "username", 
    Values:     []interface{}{"admin", "user"},
    ExactMatch: false,                             // 模糊匹配
    AllRegex:   false,                             // false: 左匹配，true: 全匹配
})

// 添加时间范围查询
param.AddTimeRange("created_at", "2023-01-01 00:00:00", "2023-12-31 23:59:59")

// 添加 FIND_IN_SET 查询
param.AddFindInSet("tags", []string{"vip", "active", "premium"})

// 执行查询
var users []User
result := handler.Query(param).Find(&users)
if result.Error != nil {
    log.Printf("Advanced query error: %v", result.Error)
}
```

### 2. 自定义过滤器

提供便捷的过滤器创建函数：

```go
// 精确匹配过滤器（IN 查询）
inFilter := database.NewInFilter("status", []interface{}{1, 2, 3})
param.AddFilter(inFilter)

// 模糊匹配过滤器
likeFilter := database.NewLikeFilter("username", []interface{}{"admin", "user"}, false)
param.AddFilter(likeFilter)

// 组合多个过滤器
param := database.NewAdvancedQueryParam(option)
param.AddFilter(database.NewInFilter("status", []interface{}{1}))
param.AddFilter(database.NewLikeFilter("email", []interface{}{"@company.com"}, true))
```

### 3. 联表查询示例

当需要进行 JOIN 查询时，正确使用 TablePrefix：

```go
// ✅ 正确的联表查询用法
func GetUsersWithOrders(handler database.Handler) ([]UserWithOrder, error) {
    // 定义结果结构
    type UserWithOrder struct {
        UserID      uint   `json:"user_id"`
        Username    string `json:"username"`
        OrderID     uint   `json:"order_id"`
        OrderAmount float64 `json:"order_amount"`
    }
    
    // 使用表前缀进行联表查询
    option := &database.FindOptionCommon{
        BusinessId:  1,
        TablePrefix: "u.",                             // 为主表设置别名前缀
        Limit:       50,
        By:          "created_at", 
        Order:       "DESC",
    }
    
    param := database.NewAdvancedQueryParam(option)
    
    var results []UserWithOrder
    err := handler.Query(param).
        Table("users u").                               // 主表使用别名 u
        Select("u.id as user_id, u.username, o.id as order_id, o.amount as order_amount").
        Joins("LEFT JOIN orders o ON u.id = o.user_id").
        Find(&results).Error
        
    return results, err
}

// ❌ 单表查询错误用法
option := &database.FindOptionCommon{
    BusinessId:  1,
    TablePrefix: "u.",                                 // 单表查询不应使用前缀
}

// ✅ 单表查询正确用法  
option := &database.FindOptionCommon{
    BusinessId:  1,
    TablePrefix: "",                                   // 单表查询留空
}
```

### 4. 动态查询构建

根据条件动态构建查询：

```go
func GetUsersByFilter(handler database.Handler, filter UserFilter) ([]User, error) {
    option := &database.FindOptionCommon{
        BusinessId: filter.BusinessID,
        Limit:      filter.Limit,
        Offset:     filter.Offset,
    }
    
    param := database.NewAdvancedQueryParam(option)
    
    // 根据条件动态添加过滤器
    if len(filter.Statuses) > 0 {
        param.AddFilter(database.NewInFilter("status", 
            interfaceSlice(filter.Statuses)))
    }
    
    if filter.UsernameKeyword != "" {
        param.AddFilter(database.NewLikeFilter("username", 
            []interface{}{filter.UsernameKeyword}, false))
    }
    
    if filter.StartDate != "" && filter.EndDate != "" {
        param.AddTimeRange("created_at", filter.StartDate, filter.EndDate)
    }
    
    if len(filter.Tags) > 0 {
        param.AddFindInSet("tags", filter.Tags)
    }
    
    var users []User
    err := handler.Query(param).Find(&users).Error
    return users, err
}

// 辅助函数：转换 slice 类型
func interfaceSlice(slice []int) []interface{} {
    result := make([]interface{}, len(slice))
    for i, v := range slice {
        result[i] = v
    }
    return result
}
```

## 💾 事务操作

### 1. 基础事务模式

```go
// 标准事务操作模板
func CreateUserWithProfile(handler database.Handler, user User, profile UserProfile) error {
    // 开始事务
    tx := handler.Begin()
    
    // 确保异常情况下的回滚
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r) // 重新抛出 panic
        }
    }()
    
    // 创建用户
    if err := tx.DB().Create(&user).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    // 设置关联ID并创建用户资料
    profile.UserID = user.ID
    if err := tx.DB().Create(&profile).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create profile: %w", err)
    }
    
    // 提交事务
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}
```

### 2. 事务中的查询操作

```go
// 在事务中进行查询和更新
func TransferBalance(handler database.Handler, fromUserID, toUserID uint, amount float64) error {
    tx := handler.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 查询转出用户
    fromParam := database.NewSimpleQueryParam("id = ?", fromUserID)
    var fromUser User
    if err := tx.Query(fromParam).First(&fromUser).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("from user not found: %w", err)
    }
    
    // 检查余额
    if fromUser.Balance < amount {
        tx.Rollback()
        return fmt.Errorf("insufficient balance")
    }
    
    // 查询转入用户
    toParam := database.NewSimpleQueryParam("id = ?", toUserID)
    var toUser User
    if err := tx.Query(toParam).First(&toUser).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("to user not found: %w", err)
    }
    
    // 更新余额
    fromUser.Balance -= amount
    toUser.Balance += amount
    
    if err := tx.DB().Save(&fromUser).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update from user: %w", err)
    }
    
    if err := tx.DB().Save(&toUser).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update to user: %w", err)
    }
    
    // 记录转账日志
    log := TransferLog{
        FromUserID: fromUserID,
        ToUserID:   toUserID,
        Amount:     amount,
    }
    if err := tx.DB().Create(&log).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create transfer log: %w", err)
    }
    
    return tx.Commit()
}
```

### 3. 批量操作事务

```go
// 批量创建用户（事务保护）
func BatchCreateUsers(handler database.Handler, users []User) error {
    tx := handler.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 批量插入（分批处理，避免单次插入过多数据）
    batchSize := 100
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        if err := tx.DB().CreateInBatches(batch, batchSize).Error; err != nil {
            tx.Rollback()
            return fmt.Errorf("batch create failed at index %d: %w", i, err)
        }
    }
    
    return tx.Commit()
}
```

### 4. 嵌套事务和保存点

```go
// 使用保存点处理复杂业务逻辑
func ComplexBusinessOperation(handler database.Handler, data BusinessData) error {
    tx := handler.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 第一阶段：创建主要数据
    mainRecord := MainRecord{Name: data.Name}
    if err := tx.DB().Create(&mainRecord).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // 创建保存点
    savePoint := tx.DB().SavePoint("sp1")
    if savePoint.Error != nil {
        tx.Rollback()
        return savePoint.Error
    }
    
    // 第二阶段：创建相关数据（可能失败但不影响主要数据）
    for _, item := range data.Items {
        item.MainRecordID = mainRecord.ID
        if err := tx.DB().Create(&item).Error; err != nil {
            // 回滚到保存点，保留主要数据
            tx.DB().RollbackTo("sp1")
            log.Printf("Failed to create item %v: %v", item, err)
            break // 跳出循环，继续后续处理
        }
    }
    
    // 第三阶段：更新统计信息
    if err := updateStatistics(tx, mainRecord.ID); err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit()
}
```

## 📄 分页查询

### 1. 使用 PageQueryParam

适用于简单的分页查询场景：

```go
// 基础分页查询
func GetUsersByPage(handler database.Handler, businessId int64, page, size int) ([]User, int64, error) {
    // 计算偏移量
    offset := (page - 1) * size
    
    // 创建分页查询参数
    param := database.NewPageQueryParam(
        "business_id = ? AND status = ?",           // 查询条件
        []interface{}{businessId, 1},               // 查询参数
        size,                                       // 每页数量
        offset,                                     // 偏移量
        "created_at DESC",                          // 排序
    )
    
    // 执行查询
    var users []User
    result := handler.Query(param).Find(&users)
    if result.Error != nil {
        return nil, 0, fmt.Errorf("failed to query users: %w", result.Error)
    }
    
    // 获取总数
    var total int64
    countResult := handler.DB().Model(&User{}).
        Where("business_id = ? AND status = ?", businessId, 1).
        Count(&total)
    if countResult.Error != nil {
        return nil, 0, fmt.Errorf("failed to count users: %w", countResult.Error)
    }
    
    return users, total, nil
}
```

### 2. 使用查询构建器分页

更强大的分页查询，支持复杂条件：

```go
// 高级分页查询
func GetUsersAdvancedPage(handler database.Handler, filter UserPageFilter) (*PageResult, error) {
    // 构建查询条件
    param := database.NewQueryBuilder().
        WithBusinessId(filter.BusinessID).
        WithShopId(filter.ShopID).
        WhereIn("status", []interface{}{1, 2}).
        WithPagination(filter.Size, (filter.Page-1)*filter.Size).
        WithOrder(filter.SortBy, filter.SortOrder).
        Build()
    
    // 添加可选的模糊查询
    if filter.Keyword != "" {
        likeParam := database.NewQueryBuilder().
            WhereLike("username", []interface{}{filter.Keyword}, false).
            Build()
        // 注意：这里需要组合查询条件，实际使用中可能需要更复杂的逻辑
    }
    
    // 执行查询
    var users []User
    result := handler.Query(param).Find(&users)
    if result.Error != nil {
        return nil, fmt.Errorf("failed to query users: %w", result.Error)
    }
    
    // 获取总数（使用相同的查询条件，但不包括分页）
    countParam := database.NewQueryBuilder().
        WithBusinessId(filter.BusinessID).
        WithShopId(filter.ShopID).
        WhereIn("status", []interface{}{1, 2}).
        Build()
    
    var total int64
    countResult := handler.Query(countParam).Model(&User{}).Count(&total)
    if countResult.Error != nil {
        return nil, fmt.Errorf("failed to count users: %w", countResult.Error)
    }
    
    // 计算分页信息
    totalPages := (total + int64(filter.Size) - 1) / int64(filter.Size)
    
    return &PageResult{
        Data:        users,
        Total:       total,
        Page:        filter.Page,
        Size:        filter.Size,
        TotalPages:  int(totalPages),
        HasNext:     filter.Page < int(totalPages),
        HasPrevious: filter.Page > 1,
    }, nil
}

// 分页结果结构
type PageResult struct {
    Data        interface{} `json:"data"`
    Total       int64       `json:"total"`
    Page        int         `json:"page"`
    Size        int         `json:"size"`
    TotalPages  int         `json:"total_pages"`
    HasNext     bool        `json:"has_next"`
    HasPrevious bool        `json:"has_previous"`
}

// 用户分页过滤器
type UserPageFilter struct {
    BusinessID int64  `json:"business_id"`
    ShopID     int64  `json:"shop_id,omitempty"`
    Keyword    string `json:"keyword,omitempty"`
    Status     []int  `json:"status,omitempty"`
    Page       int    `json:"page"`
    Size       int    `json:"size"`
    SortBy     string `json:"sort_by"`
    SortOrder  string `json:"sort_order"`
}
```

### 3. 游标分页（适用于大数据集）

对于大数据集，传统的 offset 分页性能较差，可以使用游标分页：

```go
// 游标分页查询
func GetUsersByCursor(handler database.Handler, businessId int64, cursor uint, size int) ([]User, uint, error) {
    query := "business_id = ? AND status = ?"
    args := []interface{}{businessId, 1}
    
    // 如果有游标，添加游标条件
    if cursor > 0 {
        query += " AND id > ?"
        args = append(args, cursor)
    }
    
    param := database.NewPageQueryParam(
        query,
        args,
        size+1, // 多查一条，用于判断是否还有下一页
        0,      // 游标分页不使用 offset
        "id ASC", // 必须按ID排序
    )
    
    var users []User
    result := handler.Query(param).Find(&users)
    if result.Error != nil {
        return nil, 0, fmt.Errorf("failed to query users: %w", result.Error)
    }
    
    // 处理结果
    var nextCursor uint
    hasMore := len(users) > size
    
    if hasMore {
        // 移除多查的一条记录
        users = users[:size]
        nextCursor = users[len(users)-1].ID
    }
    
    return users, nextCursor, nil
}
```

### 4. 分页查询最佳实践

```go
// 通用分页查询封装
type PaginationService struct {
    handler database.Handler
}

func NewPaginationService(handler database.Handler) *PaginationService {
    return &PaginationService{handler: handler}
}

// 通用分页方法
func (s *PaginationService) Paginate(
    model interface{},
    param database.QueryParam,
    page, size int,
) (*PageResult, error) {
    if page < 1 {
        page = 1
    }
    if size < 1 || size > 100 { // 限制最大页面大小
        size = 20
    }
    
    offset := (page - 1) * size
    
    // 查询数据
    var items []interface{}
    result := s.handler.Query(param).
        Limit(size).
        Offset(offset).
        Find(&items)
    if result.Error != nil {
        return nil, fmt.Errorf("failed to query data: %w", result.Error)
    }
    
    // 查询总数
    var total int64
    countResult := s.handler.Query(param).Model(model).Count(&total)
    if countResult.Error != nil {
        return nil, fmt.Errorf("failed to count data: %w", countResult.Error)
    }
    
    totalPages := (total + int64(size) - 1) / int64(size)
    
    return &PageResult{
        Data:        items,
        Total:       total,
        Page:        page,
        Size:        size,
        TotalPages:  int(totalPages),
        HasNext:     page < int(totalPages),
        HasPrevious: page > 1,
    }, nil
}
```

## 🗄️ 数据库兼容性

### 1. 支持的数据库类型

该 database 包设计为兼容多种数据库：

| 数据库类型 | 支持状态 | 特殊说明 |
|------------|----------|----------|
| MySQL      | ✅ 完全支持 | 推荐用于生产环境 |
| PostgreSQL | ✅ 完全支持 | 支持高级特性 |
| SQLite     | ✅ 完全支持 | 适用于测试和小型应用 |
| SQL Server | ⚠️ 基础支持 | 部分高级特性可能不可用 |
| Oracle     | ⚠️ 基础支持 | 需要额外配置 |

### 2. 数据库特定功能

#### MySQL 特性
```go
// MySQL 的 FIND_IN_SET 函数
param := database.NewQueryBuilder().
    WhereFindInSet("tags", []string{"vip", "premium"}).
    Build()

// MySQL 全文搜索（需要配置全文索引）
result := handler.DB().
    Where("MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE)", "搜索关键词").
    Find(&articles)
```

#### PostgreSQL 特性
```go
// PostgreSQL 的 JSONB 查询
result := handler.DB().
    Where("metadata->>'key' = ?", "value").
    Find(&records)

// PostgreSQL 数组查询
result := handler.DB().
    Where("tags && ?", pq.Array([]string{"tag1", "tag2"})).
    Find(&posts)
```

#### SQLite 兼容性
```go
// SQLite 不支持某些 MySQL 函数，包已自动处理兼容性
// FIND_IN_SET 会自动转换为 LIKE 查询
param := database.NewQueryBuilder().
    WhereFindInSet("tags", []string{"vip"}). // 自动转换为 LIKE 查询
    Build()

// REGEXP 会自动转换为 LIKE 查询
param = database.NewQueryBuilder().
    WhereLike("username", []interface{}{"admin"}, false). // 使用 LIKE 替代 REGEXP
    Build()
```

### 3. 跨数据库查询注意事项

```go
// ✅ 推荐：使用兼容性好的查询方式
param := database.NewQueryBuilder().
    WhereIn("status", []interface{}{1, 2}).              // 所有数据库都支持
    WhereLike("username", []interface{}{"admin"}, false). // 使用 LIKE 而不是正则
    WithOrder("created_at", "DESC").                     // 标准 SQL 排序
    Build()

// ❌ 避免：数据库特定的 SQL 函数
// 不要直接使用 MySQL 特定函数，除非确保只在 MySQL 上运行
```

## ⚡ 性能优化

### 1. 查询优化策略

#### 索引优化
```go
// 为常用查询字段创建索引
type User struct {
    ID         uint   `gorm:"primarykey"`
    Username   string `gorm:"size:50;not null;uniqueIndex"`      // 唯一索引
    Email      string `gorm:"size:100;not null;uniqueIndex"`     // 唯一索引
    BusinessID int64  `gorm:"index"`                             // 普通索引
    ShopID     int64  `gorm:"index"`                             // 普通索引
    Status     int    `gorm:"default:1;index"`                   // 状态索引
    CreatedAt  time.Time `gorm:"index"`                          // 时间索引
    
    // 复合索引
    _ struct{} `gorm:"index:idx_business_shop,unique,sort:desc"`
}

// 复合索引定义
func (User) TableName() string { return "users" }

// 为模型添加复合索引
func CreateUserIndexes(db *gorm.DB) error {
    // 业务+店铺复合索引
    err := db.Exec("CREATE INDEX idx_user_business_shop ON users(business_id, shop_id)").Error
    if err != nil {
        return err
    }
    
    // 状态+创建时间复合索引
    err = db.Exec("CREATE INDEX idx_user_status_created ON users(status, created_at DESC)").Error
    return err
}
```

#### 查询性能优化
```go
// ✅ 高效查询方式
func GetActiveUsersByBusiness(handler database.Handler, businessId int64, limit int) ([]User, error) {
    // 1. 使用索引字段查询
    param := database.NewQueryBuilder().
        WithBusinessId(businessId).           // business_id 有索引
        WhereIn("status", []interface{}{1}).  // status 有索引
        WithOrder("created_at", "DESC").      // created_at 有索引
        WithPagination(limit, 0).
        Build()
    
    var users []User
    // 2. 只选择需要的字段，减少网络传输
    result := handler.Query(param).
        Select("id, username, email, status, created_at").
        Find(&users)
    
    return users, result.Error
}

// ❌ 低效查询方式
func GetUsersSlowly(handler database.Handler) ([]User, error) {
    // 1. 避免全表扫描
    param := database.NewSimpleQueryParam("age > ?", 18) // age 字段没有索引
    
    var users []User
    // 2. 避免 SELECT *
    result := handler.Query(param).Find(&users) // 查询了所有字段
    
    return users, result.Error
}
```

### 2. 连接池优化

```go
// 配置连接池参数
func ConfigureConnectionPool(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    
    // 设置空闲连接池中连接的最大数量
    sqlDB.SetMaxIdleConns(10)
    
    // 设置打开数据库连接的最大数量
    sqlDB.SetMaxOpenConns(100)
    
    // 设置连接可复用的最大时间
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    // 设置连接空闲的最大时间
    sqlDB.SetConnMaxIdleTime(time.Minute * 10)
    
    return nil
}
```

### 3. 批量操作优化

```go
// 高效批量插入
func BatchInsertUsers(handler database.Handler, users []User) error {
    if len(users) == 0 {
        return nil
    }
    
    // 使用事务进行批量操作
    tx := handler.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 分批插入，避免单次操作过大
    batchSize := 500
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        if err := tx.DB().CreateInBatches(batch, batchSize).Error; err != nil {
            tx.Rollback()
            return fmt.Errorf("batch insert failed: %w", err)
        }
    }
    
    return tx.Commit()
}

// 批量更新优化
func BatchUpdateUserStatus(handler database.Handler, userIds []uint, status int) error {
    if len(userIds) == 0 {
        return nil
    }
    
    // 使用 IN 查询进行批量更新
    return handler.DB().
        Model(&User{}).
        Where("id IN ?", userIds).
        Update("status", status).Error
}
```

### 4. 查询缓存策略

```go
// 缓存查询结果（伪代码，需要配合缓存系统）
type CachedUserService struct {
    handler database.Handler
    cache   CacheInterface
}

func (s *CachedUserService) GetUsersByBusinessId(businessId int64) ([]User, error) {
    // 检查缓存
    cacheKey := fmt.Sprintf("users:business:%d", businessId)
    if cached, err := s.cache.Get(cacheKey); err == nil {
        var users []User
        json.Unmarshal(cached, &users)
        return users, nil
    }
    
    // 缓存未命中，查询数据库
    param := database.NewQueryBuilder().
        WithBusinessId(businessId).
        WhereIn("status", []interface{}{1}).
        Build()
    
    var users []User
    result := s.handler.Query(param).Find(&users)
    if result.Error != nil {
        return nil, result.Error
    }
    
    // 写入缓存
    if data, err := json.Marshal(users); err == nil {
        s.cache.Set(cacheKey, data, time.Minute*10)
    }
    
    return users, nil
}
```

## 🚨 错误处理

### 1. 错误类型识别

```go
import (
    "errors"
    "gorm.io/gorm"
)

func HandleDatabaseErrors(err error) {
    switch {
    case errors.Is(err, gorm.ErrRecordNotFound):
        log.Println("Record not found")
        
    case errors.Is(err, gorm.ErrInvalidTransaction):
        log.Println("Invalid transaction")
        
    case errors.Is(err, gorm.ErrDuplicatedKey):
        log.Println("Duplicated key constraint violation")
        
    case errors.Is(err, gorm.ErrCheckConstraint):
        log.Println("Check constraint violation")
        
    case errors.Is(err, gorm.ErrForeignKeyViolated):
        log.Println("Foreign key constraint violation")
        
    default:
        log.Printf("Unknown database error: %v", err)
    }
}
```

### 2. 优雅的错误处理

```go
// 统一错误处理封装
type DatabaseService struct {
    handler database.Handler
}

func (s *DatabaseService) GetUserByID(id uint) (*User, error) {
    param := database.NewSimpleQueryParam("id = ?", id)
    
    var user User
    result := s.handler.Query(param).First(&user)
    
    switch {
    case result.Error == nil:
        return &user, nil
        
    case errors.Is(result.Error, gorm.ErrRecordNotFound):
        return nil, NewNotFoundError("user", id)
        
    default:
        return nil, NewDatabaseError("failed to get user", result.Error)
    }
}

// 自定义错误类型
type NotFoundError struct {
    Resource string
    ID       interface{}
}

func (e NotFoundError) Error() string {
    return fmt.Sprintf("%s with id %v not found", e.Resource, e.ID)
}

func NewNotFoundError(resource string, id interface{}) *NotFoundError {
    return &NotFoundError{Resource: resource, ID: id}
}

type DatabaseError struct {
    Message string
    Err     error
}

func (e DatabaseError) Error() string {
    return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func NewDatabaseError(message string, err error) *DatabaseError {
    return &DatabaseError{Message: message, Err: err}
}
```

### 3. 重试机制

```go
// 数据库操作重试机制
func WithRetry(operation func() error, maxRetries int, delay time.Duration) error {
    var err error
    
    for i := 0; i <= maxRetries; i++ {
        err = operation()
        if err == nil {
            return nil
        }
        
        // 检查是否是可重试的错误
        if !isRetryableError(err) {
            return err
        }
        
        if i < maxRetries {
            log.Printf("Operation failed (attempt %d/%d): %v. Retrying in %v...", 
                i+1, maxRetries+1, err, delay)
            time.Sleep(delay)
            delay *= 2 // 指数退避
        }
    }
    
    return fmt.Errorf("operation failed after %d retries: %w", maxRetries+1, err)
}

func isRetryableError(err error) bool {
    // 判断是否是可重试的错误（如网络错误、超时等）
    errStr := err.Error()
    retryableErrors := []string{
        "connection refused",
        "timeout",
        "connection lost",
        "database is locked",
    }
    
    for _, retryableErr := range retryableErrors {
        if strings.Contains(errStr, retryableErr) {
            return true
        }
    }
    
    return false
}

// 使用重试机制
func CreateUserWithRetry(handler database.Handler, user User) error {
    return WithRetry(func() error {
        return handler.DB().Create(&user).Error
    }, 3, time.Second)
}
```

## 💡 最佳实践

### 1. 资源管理

```go
// ✅ 正确的资源管理方式
func UserService() error {
    // 获取数据库处理器
    handler := database.GetDefaultHandler()
    if handler == nil {
        return errors.New("failed to get database handler")
    }
    defer handler.Close() // 确保资源释放
    
    // 业务逻辑处理
    return processUsers(handler)
}

// ❌ 错误的资源管理方式
func BadUserService() error {
    handler := database.GetDefaultHandler()
    // 忘记释放资源，可能导致连接泄漏
    return processUsers(handler)
}
```

### 2. 查询优化最佳实践

```go
// ✅ 高效查询模式
func GetUsersByBusinessEfficient(handler database.Handler, businessId int64, page, size int) ([]User, int64, error) {
    // 1. 使用索引字段
    param := database.NewQueryBuilder().
        WithBusinessId(businessId).                    // business_id 有索引
        WhereIn("status", []interface{}{1}).           // status 有索引
        WithOrder("created_at", "DESC").               // created_at 有索引
        WithPagination(size, (page-1)*size).
        Build()
    
    var users []User
    // 2. 只查询必要字段
    result := handler.Query(param).
        Select("id, username, email, created_at").     // 避免 SELECT *
        Find(&users)
    if result.Error != nil {
        return nil, 0, result.Error
    }
    
    // 3. 高效计数查询
    countParam := database.NewQueryBuilder().
        WithBusinessId(businessId).
        WhereIn("status", []interface{}{1}).
        Build()
    
    var total int64
    countResult := handler.Query(countParam).Model(&User{}).Count(&total)
    
    return users, total, countResult.Error
}
```

### 3. 事务处理最佳实践

```go
// ✅ 完整的事务处理模板
func CompleteTransactionExample(handler database.Handler, data ProcessData) error {
    tx := handler.Begin()
    
    // 使用 defer 确保异常情况下的清理
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            log.Printf("Transaction panicked and rolled back: %v", r)
            panic(r) // 重新抛出 panic
        }
    }()
    
    // 业务逻辑 1
    if err := step1(tx, data); err != nil {
        tx.Rollback()
        return fmt.Errorf("step1 failed: %w", err)
    }
    
    // 业务逻辑 2
    if err := step2(tx, data); err != nil {
        tx.Rollback()
        return fmt.Errorf("step2 failed: %w", err)
    }
    
    // 提交事务
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}
```

### 4. 模型设计最佳实践

```go
// ✅ 良好的模型设计
type User struct {
    // 主键
    ID uint `gorm:"primarykey" json:"id"`
    
    // 业务字段（带索引）
    Username   string `gorm:"size:50;not null;uniqueIndex" json:"username"`
    Email      string `gorm:"size:100;not null;uniqueIndex" json:"email"`
    BusinessID int64  `gorm:"index" json:"business_id"`
    ShopID     int64  `gorm:"index" json:"shop_id"`
    
    // 状态字段
    Status int `gorm:"default:1;index;comment:用户状态:1=正常,2=禁用" json:"status"`
    
    // 可选字段
    Age     *int    `gorm:"comment:年龄" json:"age,omitempty"`
    Phone   *string `gorm:"size:20;comment:手机号" json:"phone,omitempty"`
    Avatar  *string `gorm:"size:255;comment:头像URL" json:"avatar,omitempty"`
    
    // 时间字段（带索引）
    CreatedAt time.Time  `gorm:"index" json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `gorm:"index" json:"-"` // 软删除
}

// 表名
func (User) TableName() string {
    return "users"
}

// 模型钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 创建前的验证和处理
    if u.Username == "" {
        return errors.New("username is required")
    }
    return nil
}
```

### 5. 并发安全

```go
// ✅ 并发安全的服务设计
type UserService struct {
    handler database.Handler
    mu      sync.RWMutex
    cache   map[string]interface{}
}

func NewUserService(handler database.Handler) *UserService {
    return &UserService{
        handler: handler,
        cache:   make(map[string]interface{}),
    }
}

// Handler 是线程安全的，可以在多个 goroutine 中使用
func (s *UserService) GetUser(id uint) (*User, error) {
    // 数据库操作是并发安全的
    param := database.NewSimpleQueryParam("id = ?", id)
    
    var user User
    result := s.handler.Query(param).First(&user)
    return &user, result.Error
}

// 本地缓存需要考虑并发安全
func (s *UserService) GetCachedUser(id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    
    // 读锁
    s.mu.RLock()
    if cached, exists := s.cache[cacheKey]; exists {
        s.mu.RUnlock()
        return cached.(*User), nil
    }
    s.mu.RUnlock()
    
    // 查询数据库
    user, err := s.GetUser(id)
    if err != nil {
        return nil, err
    }
    
    // 写锁
    s.mu.Lock()
    s.cache[cacheKey] = user
    s.mu.Unlock()
    
    return user, nil
}

## 📚 API 参考

### 1. Handler 接口方法

#### 核心方法
```go
// DB() 获取原始 GORM 实例
func (h *Handler) DB() *gorm.DB

// Query() 执行查询参数
func (h *Handler) Query(param QueryParam) *gorm.DB

// Close() 关闭数据库连接
func (h *Handler) Close() error

// AutoMigrate() 自动迁移表结构
func (h *Handler) AutoMigrate(dst ...any) error
```

#### 事务方法
```go
// Begin() 开始事务
func (h *Handler) Begin(opts ...*sql.TxOptions) Handler

// Commit() 提交事务
func (h *Handler) Commit() error

// Rollback() 回滚事务
func (h *Handler) Rollback() error
```

### 2. 查询构建器方法

#### 基础条件
```go
// WithBusinessId() 设置业务ID
func (qb *QueryBuilder) WithBusinessId(businessId int64) *QueryBuilder

// WithShopId() 设置店铺ID
func (qb *QueryBuilder) WithShopId(shopId int64) *QueryBuilder

// WhereIn() 添加IN条件
func (qb *QueryBuilder) WhereIn(field string, values []interface{}) *QueryBuilder

// WhereLike() 添加LIKE条件
func (qb *QueryBuilder) WhereLike(field string, values []interface{}, allRegex bool) *QueryBuilder
```

#### 时间和特殊查询
```go
// WhereTimeRange() 添加时间范围条件
func (qb *QueryBuilder) WhereTimeRange(field, startTime, endTime string) *QueryBuilder

// WhereFindInSet() 添加FIND_IN_SET条件
func (qb *QueryBuilder) WhereFindInSet(field string, values []string) *QueryBuilder
```

#### 分页和排序
```go
// WithPagination() 设置分页
func (qb *QueryBuilder) WithPagination(limit, offset int) *QueryBuilder

// WithOrder() 设置排序
func (qb *QueryBuilder) WithOrder(field, order string) *QueryBuilder

// WithGroupBy() 设置分组
func (qb *QueryBuilder) WithGroupBy(groupBy string) *QueryBuilder
```

### 3. 工厂方法

#### 查询参数创建
```go
// NewSimpleQueryParam() 创建简单查询参数
func NewSimpleQueryParam(query string, args ...interface{}) QueryParam

// NewPageQueryParam() 创建分页查询参数
func NewPageQueryParam(query string, args []interface{}, limit, offset int, orderBy string) QueryParam

// NewAdvancedQueryParam() 创建高级查询参数
func NewAdvancedQueryParam(option *FindOptionCommon) *AdvancedQueryParam

// NewQueryBuilder() 创建查询构建器
func NewQueryBuilder() *QueryBuilder
```

#### 过滤器创建
```go
// NewInFilter() 创建IN过滤器
func NewInFilter(field string, values []interface{}) *BaseInfoFilter

// NewLikeFilter() 创建LIKE过滤器
func NewLikeFilter(field string, values []interface{}, allRegex bool) *BaseInfoFilter
```

## ⚙️ 配置示例

### 1. MySQL 配置

```go
package main

import (
    "fmt"
    "time"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "github.com/kamalyes/go-core/pkg/database"
)

func setupMySQLDatabase() database.Handler {
    // 数据库连接配置
    config := MySQLConfig{
        Host:     "127.0.0.1",
        Port:     "3306",
        Username: "root",
        Password: "password",
        Database: "myapp",
        Charset:  "utf8mb4",
        ParseTime: true,
        Loc:      "Local",
    }
    
    // 构建 DSN
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
        config.Username,
        config.Password,
        config.Host,
        config.Port,
        config.Database,
        config.Charset,
        config.ParseTime,
        config.Loc,
    )
    
    // GORM 配置
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info), // 日志级别
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true, // 使用单数表名
        },
    }
    
    // 创建数据库连接
    db, err := gorm.Open(mysql.Open(dsn), gormConfig)
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to MySQL: %v", err))
    }
    
    // 连接池配置
    sqlDB, err := db.DB()
    if err != nil {
        panic(fmt.Sprintf("Failed to get underlying sql.DB: %v", err))
    }
    
    // 连接池参数
    sqlDB.SetMaxIdleConns(10)                    // 空闲连接数
    sqlDB.SetMaxOpenConns(100)                   // 最大连接数
    sqlDB.SetConnMaxLifetime(time.Hour)          // 连接最大生存时间
    sqlDB.SetConnMaxIdleTime(time.Minute * 10)   // 连接最大空闲时间
    
    return database.NewHandler(db)
}

type MySQLConfig struct {
    Host      string
    Port      string
    Username  string
    Password  string
    Database  string
    Charset   string
    ParseTime bool
    Loc       string
}
```

### 2. PostgreSQL 配置

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func setupPostgreSQLDatabase() database.Handler {
    // PostgreSQL 连接配置
    dsn := "host=localhost user=postgres password=password dbname=myapp port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    
    // GORM 配置
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    }
    
    db, err := gorm.Open(postgres.Open(dsn), gormConfig)
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to PostgreSQL: %v", err))
    }
    
    // 连接池配置
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return database.NewHandler(db)
}
```

### 3. SQLite 配置

```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupSQLiteDatabase() database.Handler {
    // SQLite 数据库文件路径
    dbPath := "./app.db"
    
    // GORM 配置
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent), // SQLite 通常使用 Silent 模式
    }
    
    db, err := gorm.Open(sqlite.Open(dbPath), gormConfig)
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to SQLite: %v", err))
    }
    
    // SQLite 特殊配置
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(1)  // SQLite 推荐使用单连接
    sqlDB.SetMaxOpenConns(1)  // SQLite 不支持真正的并发写入
    
    return database.NewHandler(db)
}
```

### 4. 生产环境配置

```go
// 生产环境数据库配置
func setupProductionDatabase() database.Handler {
    // 从环境变量获取配置
    config := &DatabaseConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnv("DB_PORT", "3306"),
        Username: getEnv("DB_USERNAME", "root"),
        Password: getEnv("DB_PASSWORD", ""),
        Database: getEnv("DB_NAME", "myapp"),
        
        // 连接池配置
        MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
        MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
        ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", time.Hour),
    }
    
    // 构建 DSN
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
        config.Username,
        config.Password,
        config.Host,
        config.Port,
        config.Database,
    )
    
    // 生产环境 GORM 配置
    gormConfig := &gorm.Config{
        Logger:                 logger.Default.LogMode(logger.Error), // 只记录错误日志
        DisableForeignKeyConstraintWhenMigrating: true,               // 禁用外键约束（提升性能）
        PrepareStmt:           true,                                  // 启用预处理语句缓存
        SkipDefaultTransaction: true,                                 // 跳过默认事务（提升性能）
    }
    
    db, err := gorm.Open(mysql.Open(dsn), gormConfig)
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to database: %v", err))
    }
    
    // 连接池配置
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(config.MaxIdleConns)
    sqlDB.SetMaxOpenConns(config.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
    
    return database.NewHandler(db)
}

// 辅助函数
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}

type DatabaseConfig struct {
    Host            string
    Port            string
    Username        string
    Password        string
    Database        string
    MaxIdleConns    int
    MaxOpenConns    int
    ConnMaxLifetime time.Duration
}
```

## ❓ 常见问题

### 1. 连接相关问题

**Q: 数据库连接池满了怎么办？**

A: 检查连接池配置和连接泄漏：

```go
// 正确管理连接
func correctUsage() error {
    handler := database.GetDefaultHandler()
    defer handler.Close() // 确保释放连接
    
    // 业务逻辑...
    return nil
}

// 调整连接池大小
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(200)        // 增加最大连接数
sqlDB.SetMaxIdleConns(20)         // 增加空闲连接数
sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接生命周期
```

**Q: 如何处理数据库连接断开？**

A: 实现重连机制：

```go
func withReconnect(handler database.Handler, operation func() error) error {
    err := operation()
    if isConnectionError(err) {
        // 尝试重新连接
        newHandler := database.GetDefaultHandler()
        defer newHandler.Close()
        
        return operation()
    }
    return err
}

func isConnectionError(err error) bool {
    if err == nil {
        return false
    }
    errStr := strings.ToLower(err.Error())
    return strings.Contains(errStr, "connection") || 
           strings.Contains(errStr, "broken pipe") ||
           strings.Contains(errStr, "reset by peer")
}
```

### 2. 性能相关问题

**Q: 查询速度慢怎么优化？**

A: 从多个角度优化：

```go
// 1. 添加索引
type User struct {
    BusinessID int64 `gorm:"index"`           // 单字段索引
    ShopID     int64 `gorm:"index"`           // 单字段索引
    Status     int   `gorm:"index"`           // 状态索引
    CreatedAt  time.Time `gorm:"index"`       // 时间索引
}

// 2. 使用复合索引
// CREATE INDEX idx_business_shop_status ON users(business_id, shop_id, status);

// 3. 只查询需要的字段
param := database.NewQueryBuilder().
    WithBusinessId(1).
    Build()

var users []User
handler.Query(param).
    Select("id, username, email"). // 避免 SELECT *
    Find(&users)

// 4. 使用 LIMIT 限制返回数量
param = database.NewQueryBuilder().
    WithBusinessId(1).
    WithPagination(50, 0). // 限制返回数量
    Build()
```

**Q: 大批量数据操作如何优化？**

A: 使用批量操作：

```go
// 批量插入
func BatchInsert(handler database.Handler, users []User) error {
    batchSize := 1000
    
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        if err := handler.DB().CreateInBatches(batch, batchSize).Error; err != nil {
            return err
        }
    }
    
    return nil
}

// 批量更新
func BatchUpdate(handler database.Handler, userIds []uint, status int) error {
    return handler.DB().
        Model(&User{}).
        Where("id IN ?", userIds).
        Update("status", status).Error
}
```

### 3. 事务相关问题

**Q: 事务死锁怎么处理？**

A: 避免死锁的最佳实践：

```go
// 1. 按相同顺序获取锁
func transferMoney(handler database.Handler, fromID, toID uint, amount float64) error {
    // 确保按ID顺序获取锁，避免死锁
    if fromID > toID {
        fromID, toID = toID, fromID
    }
    
    tx := handler.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 按顺序锁定记录
    var fromUser, toUser User
    tx.DB().Clauses(clause.Locking{Strength: "UPDATE"}).First(&fromUser, fromID)
    tx.DB().Clauses(clause.Locking{Strength: "UPDATE"}).First(&toUser, toID)
    
    // 执行转账逻辑...
    
    return tx.Commit()
}

// 2. 使用超时和重试
func withDeadlockRetry(operation func() error) error {
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        if strings.Contains(strings.ToLower(err.Error()), "deadlock") {
            time.Sleep(time.Duration(i*100) * time.Millisecond) // 递增延迟
            continue
        }
        
        return err
    }
    
    return errors.New("max deadlock retries exceeded")
}
```

### 4. 数据一致性问题

**Q: 如何确保数据一致性？**

A: 使用事务和约束：

```go
// 1. 使用事务确保操作原子性
func createOrderWithItems(handler database.Handler, order Order, items []OrderItem) error {
    tx := handler.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 创建订单
    if err := tx.DB().Create(&order).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // 创建订单项
    for i := range items {
        items[i].OrderID = order.ID
    }
    
    if err := tx.DB().CreateInBatches(items, 100).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit()
}

// 2. 使用数据库约束
type Order struct {
    ID         uint      `gorm:"primarykey"`
    UserID     uint      `gorm:"not null;index"` // 外键约束
    TotalAmount float64  `gorm:"check:total_amount >= 0"` // 检查约束
    Status     int       `gorm:"default:1;check:status IN (1,2,3,4)"`
    CreatedAt  time.Time
}

type OrderItem struct {
    ID       uint    `gorm:"primarykey"`
    OrderID  uint    `gorm:"not null;index"` // 外键
    Quantity int     `gorm:"check:quantity > 0"` // 数量必须大于0
    Price    float64 `gorm:"check:price >= 0"`   // 价格不能为负
}
```

### 5. 调试和监控

**Q: 如何调试 SQL 查询？**

A: 启用 SQL 日志：

```go
import "gorm.io/gorm/logger"

// 开发环境：启用详细日志
gormConfig := &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // 记录所有 SQL
}

// 自定义日志格式
newLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
    logger.Config{
        SlowThreshold:             time.Second,  // 慢查询阈值
        LogLevel:                  logger.Info,  // 日志级别
        IgnoreRecordNotFoundError: true,         // 忽略 RecordNotFound 错误
        Colorful:                  true,         // 彩色输出
    },
)

gormConfig := &gorm.Config{Logger: newLogger}
```

**Q: 如何监控数据库性能？**

A: 添加监控指标：

```go
// 监控连接池状态
func monitorConnectionPool(db *gorm.DB) {
    sqlDB, _ := db.DB()
    
    stats := sqlDB.Stats()
    fmt.Printf("Open connections: %d\n", stats.OpenConnections)
    fmt.Printf("In use: %d\n", stats.InUse)
    fmt.Printf("Idle: %d\n", stats.Idle)
    fmt.Printf("Wait count: %d\n", stats.WaitCount)
    fmt.Printf("Wait duration: %v\n", stats.WaitDuration)
}

// 监控查询执行时间
func monitorQuery(handler database.Handler, param database.QueryParam) time.Duration {
    start := time.Now()
    
    var users []User
    handler.Query(param).Find(&users)
    
    duration := time.Since(start)
    if duration > time.Second {
        log.Printf("Slow query detected: %v", duration)
    }
    
    return duration
}
```

---

## 📄 总结

这个 database 包为 Go 应用提供了一个强大、灵活且易用的数据库操作层。通过合理使用其提供的功能，可以构建高性能、可维护的数据库应用。

### 关键要点回顾

1. **正确使用资源管理**：始终记得关闭数据库连接
2. **合理设计查询**：使用索引字段，避免全表扫描  
3. **恰当处理事务**：确保异常情况下的回滚机制
4. **优化性能**：批量操作、连接池配置、查询优化
5. **处理错误**：完善的错误处理和重试机制
6. **监控和调试**：适当的日志记录和性能监控

遵循这些最佳实践，可以充分发挥 database 包的能力，构建稳定可靠的数据库应用。