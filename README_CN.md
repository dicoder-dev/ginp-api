# Ginp-API

基于 Go + Gin 框架开发的 RESTful API 项目，内置强大的代码生成工具。

## 核心特性

### 🚀 Ginp 框架扩展
对 Gin 框架进行增强封装，提供便捷的 API：

**ContextPlus - 便捷的响应方法：**
```go
// 成功响应
c.Success()

// 带数据的成功响应
c.SuccessData(&User{ID: 1, Name: "张三"})

// 失败响应
c.Fail("参数错误")

// 带数据的失败响应
c.FailData(err.Error(), map[string]any{"field": "username"})
```

**自动参数绑定：**
```go
// 定义请求结构体
type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Email    string `json:"email"`
}

// 控制器 - 参数自动绑定！
func CreateUser(c *ginp.ContextPlus, params *CreateUserRequest) error {
    user, err := service.CreateUser(params.Username, params.Password)
    if err != nil {
        return err  // 自动返回失败响应
    }
    return nil  // 自动返回成功响应
}

// 注册路由
ginp.RouterAppend(ginp.RouterItem{
    Path:     "/api/user/create",
    Handler:  ginp.BindParamsHandler(CreateUser, CreateUserRequest{}),
    HttpType: ginp.HttpPost,
})
```

**路由管理支持别名：**
```go
ginp.RouterAppend(ginp.RouterItem{
    Path:        "/api/user/info",
    AliasePaths: []string{"/api/user/profile", "/api/u/info"},
    Handler:     ginp.BindParamsHandler(GetUserInfo, Request{}),
    HttpType:    ginp.HttpGet,
})
```

### 🔧 代码生成工具
交互式菜单，自动化生成 CRUD 代码：

```bash
go run cmd/gencode/main.go

=== GAPI 代码生成工具 ===
1. 生成实体 CRUD 代码          # 自动生成 Controller/Service/Model
2. 新增 API 接口控制器         # 在已有模块中添加自定义 API
3. 生成实体字段常量            # 扫描实体，生成 FieldXXX 常量
4. 删除实体 CRUD 代码          # 删除实体的所有 CRUD 文件
```

### 📦 Where 查询构建器
优雅的链式调用查询条件构建器：

```go
// 简单条件
wheres := where.New("id", "=", 1).Conditions()

// 链式 AND
wheres := where.New("status", "=", 1).
    And("username", "=", "admin").
    And("email", "LIKE", "%@example.com").
    Conditions()

// 链式 OR
wheres := where.New("status", "=", 1).
    Or("status", "=", 2).
    Conditions()

// 支持多种操作符
wheres := where.New("age", ">", 18).
    And("score", ">=", 60).
    And("status", "IN", []int{1,2,3}).
    Conditions()
```

### 🗄️ 多数据库支持
```go
// MySQL
db.InitMySQL()

// PostgreSQL
db.InitPgSQL()

// SQLite
db.InitSqlite()
```

### 📁 分层架构
```
Controller → Service → Model → Entity
```

## 快速开始

### 1. 创建实体
```go
// internal/gapi/entity/user.e.go
package entity

type User struct {
    ID        uint
    Username  string
    Password  string
    Email     string
}

func (User) TableName() string {
    return "sys_user"
}

func (User) GenConfig() *gen.EntityConfig {
    return &gen.EntityConfig{
        TableName: "sys_user",
    }
}
```

### 2. 生成代码
```bash
go run cmd/gencode/main.go
# 选择选项 1，输入实体名称 "user"
```

### 3. 实现业务逻辑
在 `internal/gapi/controller/` 目录下修改生成的控制器文件

### 4. 启动服务
```bash
go run cmd/gapi/main.go
```

## 项目结构

```
ginp-api/
├── cmd/
│   ├── gapi/           # 主应用程序
│   └── gencode/        # 代码生成工具
├── internal/gapi/
│   ├── controller/    # HTTP 处理器 (.a.go)
│   ├── service/       # 业务逻辑 (.s.go)
│   ├── model/         # 数据库操作 (.m.go)
│   ├── entity/        # 数据模型 (.e.go)
│   └── router/        # 路由配置
├── pkg/
│   ├── ginp/          # 框架扩展
│   ├── where/         # 查询构建器
│   ├── dbops/         # 数据库操作
│   └── utils/         # 工具函数
└── configs/           # 配置文件
```

## 配置

编辑 `configs/config.yaml`:
```yaml
server:
  port: 8080

database:
  mysql:
    host: localhost
    port: 3306
    user: root
    password: ""
    dbname: ginp_api
```

## API 响应格式

```json
{
  "code": 1,
  "msg": "success",
  "data": {...}
}
```

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL / PostgreSQL / SQLite

## 许可证

MIT
