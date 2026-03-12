# Where 查询操作文档

本文档详细介绍 `pkg/where` 和 `pkg/dbops` 包中的查询相关功能。

---

## 一、where 包概述

`pkg/where` 包负责构建查询条件，支持链式调用，语法简洁。

### 1.1 文件结构

| 文件 | 说明 |
|------|------|
| `condition.go` | Condition 结构体定义，操作符函数 |
| `manager.go` | whereManager 链式构建器 |
| `opts.go` | 操作符常量定义 |
| `where.go` | 条件格式化辅助函数 |
| `extra.go` | 分页、排序等额外参数 |
| `check.go` | 条件校验函数 |
| `func.go` | 通用工具函数 |

---

## 二、条件构建

### 2.1 基础方式

直接创建 Condition 切片：

```go
// 单条件
wheres := where.FormatOne("id", "=", 1)

// 多条件（AND 连接）
wheres := where.Format(
    where.FormatOne("status", "=", 1),
    where.FormatOne("username", "LIKE", "%admin%"),
)

// OR 条件
wheres := where.Format(
    where.FormatOneOr("status", "=", 1),
    where.FormatOneOr("status", "=", 2),
)
```

### 2.2 链式方式（推荐）

使用 whereManager 进行链式构建：

```go
// 创建单条件
wheres := where.New("id", "=", 1).Conditions()

// 链式添加 AND 条件
wheres := where.New("status", "=", 1).
    And("username", "=", "admin").
    And("email", "LIKE", "%@example.com").
    Conditions()

// 链式添加 OR 条件
wheres := where.New("status", "=", 1).
    Or("status", "=", 2).
    Or("status", "=", 3).
    Conditions()

// 混合 AND/OR（需要手动设置连接符）
wheres := where.New("status", "=", 1).Conditions()
```

---

## 三、操作符详解

### 3.1 支持的操作符常量

```go
// pkg/where/opts.go
const OptLike = "LIKE"      // 模糊匹配
const OptIn = "IN"          // IN 查询
const OptBetween = "BETWEEN" // 范围查询
const Greater = ">"         // 大于
const GreaterEqual = ">="   // 大于等于
const Less = "<"            // 小于
const LessEqual = "<="      // 小于等于
const Equal = "="           // 等于
```

### 3.2 操作符使用示例

| 操作符 | 说明 | 示例 | 示例值 |
|--------|------|------|--------|
| `=` | 等于 | `where.New("id", "=", 1)` | `1` |
| `>` | 大于 | `where.New("age", ">", 18)` | `18` |
| `<` | 小于 | `where.New("age", "<", 65)` | `65` |
| `>=` | 大于等于 | `where.New("score", ">=", 60)` | `60` |
| `<=` | 小于等于 | `where.New("stock", "<=", 0)` | `0` |
| `LIKE` | 模糊匹配 | `where.New("username", "LIKE", "%admin%")` | `"%admin%"` |
| `IN` | IN 查询 | `where.New("status", "IN", []int{1,2,3})` | `[]int{1, 2, 3}` |
| `BETWEEN` | 范围查询 | `where.New("age", "BETWEEN", []int{18, 30})` | `[]int{18, 30}` |

### 3.3 便捷函数

```go
// 使用 OptEqual 创建等值条件
wheres := where.OptEqual("id", 1).Conditions()

// 使用 Opt 创建任意操作符条件
wheres := where.Opt("username", "LIKE", "%test%").Conditions()
```

---

## 四、Extra 额外参数

### 4.1 Extra 结构体

```go
type Extra struct {
    OrderByColumn string // 排序字段
    OrderByDesc   bool   // 是否倒序
    PageSize      int    // 每页数量
    PageNum       int    // 页码
}
```

### 4.2 创建方式

```go
// 方式1：直接创建
extra := &where.Extra{
    OrderByColumn: "created_at",
    OrderByDesc:   true,
    PageSize:      10,
    PageNum:       1,
}

// 方式2：使用构造函数
extra := where.NewExtra()

// 方式3：使用便捷构造函数（默认按 created_at 排序）
extra := where.NewExtraParam(10, true) // 10条，倒序

// 方式4：链式调用
extra := where.NewExtra().
    PSize(10).      // 设置 PageSize
    PNum(1).       // 设置 PageNum
    OrderBy("id", false) // 设置排序
```

### 4.3 参数说明

| 方法 | 说明 |
|------|------|
| `PSize(size)` | 设置每页数量 |
| `PNum(num)` | 设置页码 |
| `OrderBy(column, desc)` | 设置排序字段和是否倒序 |

---

## 五、dbops 查询操作

### 5.1 FindOneConfig - 查询单条

```go
type FindOneConfig struct {
    Fields         []string              // 要查询的字段列表
    Wheres         []*where.Condition    // 查询条件
    NewEntity      any                   // 结果载体（指针）
    Db             *gorm.DB              // 数据库实例
    getSoftDelData bool                  // 是否包含软删除数据
    RelationList   []*RelationItem       // 关联查询配置
}
```

### 5.2 FindListConfig - 查询列表

```go
type FindListConfig struct {
    Conditions     []*where.Condition    // 查询条件
    Extra          *where.Extra          // 分页、排序参数
    NewEntityList  any                   // 结果载体（指针）
    GetSoftDelData bool                  // 是否包含软删除数据
    Db             *gorm.DB              // 数据库实例
    RelationList   []*RelationItem       // 关联查询配置
    Fields         []string              // 要查询的字段列表
}
```

---

## 六、查询示例

### 6.1 查询单条数据

```go
// 简单条件查询
entityInfo := new(Entity)
err := dbops.FindOne(&dbops.FindOneConfig{
    Wheres:    where.New("id", "=", 1).Conditions(),
    Db:        dbRead,
    NewEntity: entityInfo,
})

// 指定查询字段
entityInfo := new(Entity)
err := dbops.FindOne(&dbops.FindOneConfig{
    Wheres:    where.New("id", "=", 1).Conditions(),
    Db:        dbRead,
    NewEntity: entityInfo,
    Fields:    []string{"id", "username", "email"},
})

// 包含软删除的数据
entityInfo := new(Entity)
err := dbops.FindOne(&dbops.FindOneConfig{
    Wheres:          where.New("id", "=", 1).Conditions(),
    Db:              dbRead,
    NewEntity:       entityInfo,
    getSoftDelData:  true,
})
```

### 6.2 查询列表数据

```go
// 简单列表查询
var entityList []*Entity
err := dbops.FindList(&dbops.FindListConfig{
    Conditions:    where.New("status", "=", 1).Conditions(),
    Db:           dbRead,
    NewEntityList: &entityList,
})

// 带分页和排序
var entityList []*Entity
err := dbops.FindList(&dbops.FindListConfig{
    Conditions:    where.New("status", "=", 1).Conditions(),
    Db:            dbRead,
    Extra:         where.NewExtra().PSize(10).PNum(1).OrderBy("created_at", true),
    NewEntityList: &entityList,
})

// 指定查询字段
var entityList []*Entity
err := dbops.FindList(&dbops.FindListConfig{
    Conditions:    where.New("status", "=", 1).Conditions(),
    Db:            dbRead,
    NewEntityList: &entityList,
    Fields:        []string{"id", "username", "email"},
})

// 包含软删除的数据
var entityList []*Entity
err := dbops.FindList(&dbops.FindListConfig{
    Conditions:      where.New("status", "=", 1).Conditions(),
    Db:              dbRead,
    NewEntityList:   &entityList,
    GetSoftDelData:  true,
})
```

### 6.3 统计总数

```go
// 获取符合条件的数据总数
total, err := dbops.GetTotal(
    where.New("status", "=", 1).Conditions(), // 条件
    new(Entity),                               // 实体类型
    dbRead,                                   // 数据库实例
)
```

### 6.4 关联查询

#### RelationItem 结构体

```go
type RelationItem struct {
    RelationName string              // 关联名称（对应 GORM 的 Associations 标签）
    Wheres       []*where.Condition  // 关联查询的筛选条件
}
```

#### 关联查询示例

```go
// 查询用户及其角色
var entityList []*Entity
err := dbops.FindList(&dbops.FindListConfig{
    Conditions:    where.New("status", "=", 1).Conditions(),
    Db:            dbRead,
    NewEntityList: &entityList,
    RelationList: []*dbops.RelationItem{
        {
            RelationName: "Roles", // 对应实体中的 Roles 字段
            Wheres:       where.New("status", "=", 1).Conditions(), // 可选：关联查询筛选条件
        },
    },
})

// 单条关联查询
entityInfo := new(Entity)
err := dbops.FindOne(&dbops.FindOneConfig{
    Wheres:    where.New("id", "=", 1).Conditions(),
    Db:        dbRead,
    NewEntity: entityInfo,
    RelationList: []*dbops.RelationItem{
        {
            RelationName: "Roles",
        },
    },
})
```

---

## 七、复杂查询场景

### 7.1 多条件组合

```go
// AND 多条件
wheres := where.New("status", "=", 1).
    And("age", ">", 18).
    And("city", "=", "Beijing").
    Conditions()

// OR 多条件
wheres := where.Format(
    where.New("status", "=", 1).Conditions()[0],
    where.FormatOneOr("status", "=", 2)[0],
)
// 或使用 New + Or（注意 Or 方法的 Connect 默认是 AND，需要调整）
wheres := where.New("status", "=", 1).
    Or("status", "=", 2). // 注意：Or 方法内部设置为 AND，需改进
    Conditions()
```

### 7.2 IN 查询

```go
// 查询状态为 1, 2, 3 的记录
wheres := where.New("status", "IN", []int{1, 2, 3}).Conditions()

// 查询用户名在列表中的记录
wheres := where.New("username", "IN", []string{"admin", "user", "guest"}).Conditions()
```

### 7.3 BETWEEN 范围查询

```go
// 查询年龄在 18-30 之间的记录
wheres := where.New("age", "BETWEEN", []int{18, 30}).Conditions()

// 查询创建时间范围内的记录
wheres := where.New("created_at", "BETWEEN", []string{"2024-01-01", "2024-12-31"}).Conditions()
```

### 7.4 LIKE 模糊查询

```go
// 前后模糊匹配
wheres := where.New("username", "LIKE", "%admin%").Conditions()

// 前缀匹配
wheres := where.New("username", "LIKE", "admin%").Conditions()

// 后缀匹配
wheres := where.New("username", "LIKE", "%admin").Conditions()
```

### 7.5 分页查询完整示例

```go
// Service 层分页查询
func GetUserList(page, pageSize int, status int) ([]*entity.User, uint, error) {
    // 构建查询条件
    wheres := where.New("status", "=", status).Conditions()

    // 构建分页参数
    extra := where.NewExtra().
        PSize(pageSize).
        PNum(page).
        OrderBy("created_at", true) // 按创建时间倒序

    // 执行查询
    var list []*entity.User
    err := dbops.FindList(&dbops.FindListConfig{
        Conditions:    wheres,
        Db:            dbRead,
        Extra:         extra,
        NewEntityList: &list,
    })
    if err != nil {
        return nil, 0, err
    }

    // 获取总数
    total, err := dbops.GetTotal(wheres, new(entity.User), dbRead)
    if err != nil {
        return []*entity.User{}, 0, err
    }

    return list, uint(total), nil
}
```

---

## 八、条件校验

### 8.1 Check 函数

在执行查询前校验条件是否合法：

```go
wheres := where.New("id", "=", 1).
    And("status", "=", 1).
    Conditions()

// 校验条件
err := where.Check(wheres)
if err != nil {
    // 处理错误
}
```

### 8.2 校验内容

- 检查条件是否为 nil
- 检查操作符是否合法（=, >, <, >=, <=, LIKE, IN, BETWEEN）
- 检查值是否包含危险字符（SQL 注入预防）

---

## 九、Model 层集成

在实际项目中，通常会将 dbops 封装到 Model 层：

```go
// Model 层定义
type Model struct {
    dbWrite *gorm.DB
    dbRead  *gorm.DB
}

// 查询单条
func (m *Model) FindOne(wheres []*where.Condition) (*Entity, error) {
    entityInfo := new(Entity)
    err := dbops.FindOne(&dbops.FindOneConfig{
        Wheres:    wheres,
        Db:        m.dbRead,
        NewEntity: entityInfo,
    })
    return entityInfo, err
}

// 分页列表查询
func (m *Model) FindList(wheres []*where.Condition, extra *where.Extra) ([]*Entity, uint, error) {
    var list []*Entity
    err := dbops.FindList(&dbops.FindListConfig{
        Conditions:    wheres,
        Db:            m.dbRead,
        Extra:         extra,
        NewEntityList: &list,
    })
    if err != nil {
        return nil, 0, err
    }
    total, _ := dbops.GetTotal(wheres, new(Entity), m.dbRead)
    return list, uint(total), nil
}
```

---

## 十、总结

| 功能 | 函数/结构体 | 说明 |
|------|-------------|------|
| 创建条件 | `where.New()` | 链式构建查询条件 |
| 添加 AND | `.And()` | 添加 AND 条件 |
| 添加 OR | `.Or()` | 添加 OR 条件 |
| 分页排序 | `where.Extra` | 分页、排序参数 |
| 查询单条 | `dbops.FindOne()` | 查询单条数据 |
| 查询列表 | `dbops.FindList()` | 查询列表数据 |
| 统计总数 | `dbops.GetTotal()` | 统计符合条件的数据总数 |
| 关联查询 | `RelationItem` | 预加载关联数据 |
| 条件校验 | `where.Check()` | 校验条件合法性 |
