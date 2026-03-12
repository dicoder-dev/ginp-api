# Ginp-API

基于 Go + Gin 框架开发的 RESTful API 项目，内置强大的代码生成工具，开箱即用。

## 快速开始 (5 分钟)

### 1. 定义实体
```go
// internal/gapi/entity/user.e.go
package entity

type User struct {
    ID       uint
    Username string
    Password string
    Email    string
    Status   int
}

func (User) TableName() string { return "sys_user" }
func (User) GenConfig() *gen.EntityConfig {
    return &gen.EntityConfig{TableName: "sys_user"}
}
```

### 2. 生成CRUD代码
```bash
go run cmd/gencode/main.go
# 选择: 1 → 输入: user
```

### 3. 启动
```bash
go run cmd/gapi/main.go
```

---

## 完整接口示例文件

使用代码生成工具生成的文件：
```
internal/gapi/controller/user/cuser/sys_user_create.a.go
```

```go
package cuser

import (
    "ginp-api/internal/gapi/service/user/suser"
    "ginp-api/pkg/ginp"
)

// 请求结构体 - 自动从 JSON 绑定
type RequestSysUserCreate struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Email    string `json:"email"`
    Status   int    `json:"status"`
}

// 响应结构体
type RespondSysUserCreate struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
}

// 控制器处理函数 - 参数自动绑定！
func SysUserCreate(c *ginp.ContextPlus, params *RequestSysUserCreate) error {
    id, err := suser.Create(params.Username, params.Password, params.Email, params.Status)
    if err != nil {
        return err  // 自动返回: {"code": 0, "msg": "错误信息"}
    }
    // 自动返回: {"code": 1, "msg": "success", "data": {...}}
    return c.SuccessData(RespondSysUserCreate{
        ID:       id,
        Username: params.Username,
    })
}

// 路由注册 (通过 init 自动调用)
func init() {
    ginp.RouterAppend(ginp.RouterItem{
        Path:     "/api/user/create",
        Handler:  ginp.BindParamsHandler(SysUserCreate, RequestSysUserCreate{}),
        HttpType: ginp.HttpPost,
        NeedLogin: false,
        Swagger: &ginp.SwaggerInfo{
            Title:         "创建用户",
            Description:   "创建新用户账号",
            RequestParams: RequestSysUserCreate{},
        },
    })
}
```

---

## 核心特性

### 🚀 Ginp 框架扩展

**自动参数绑定：**
```go
// 定义结构体，添加 binding 标签
type Request struct {
    Name string `json:"name" binding:"required"`
    Age  int    `json:"age"`
}

// 处理器直接接收绑定好的参数，无需手动解析！
func Handler(c *ginp.ContextPlus, params *Request) error {
    // params 已经自动从请求体绑定完成
    return c.SuccessData(params)
}

// 一行代码注册路由
ginp.RouterAppend(ginp.RouterItem{
    Path:    "/api/test",
    Handler: ginp.BindParamsHandler(Handler, Request{}),
    HttpType: ginp.HttpPost,
})
```

**便捷的响应方法：**
```go
c.Success()                                    // {"code": 1, "msg": "success"}
c.SuccessData(user)                           // {"code": 1, "msg": "success", "data": {...}}
c.Fail("错误信息")                              // {"code": 0, "msg": "错误信息"}
c.FailData("错误信息", map[string]any{"key": 1}) // {"code": 0, "msg": "错误信息", "data": {...}}
```

**路由别名支持：**
```go
ginp.RouterAppend(ginp.RouterItem{
    Path:        "/api/user/info",
    AliasePaths: []string{"/api/user/profile", "/api/u/info"},
    // ...
})
```

### 🔧 代码生成工具
```bash
go run cmd/gencode/main.go

=== GAPI 代码生成工具 ===
1. 生成实体 CRUD 代码          # 自动创建所有层级代码
2. 新增 API 接口控制器        # 添加自定义 API
3. 生成实体字段常量           # 扫描实体生成 FieldXXX
4. 删除实体 CRUD 代码         # 删除实体的所有文件
```

### 📦 Where 查询构建器
```go
// 链式查询
wheres := where.New("status", "=", 1).
    And("username", "LIKE", "%admin%").
    Or("email", "=", "admin@test.com").
    Conditions()

// 在 Model 中使用
users, err := model.Where(wheres).Limit(10).Find()
```

### ⚙️ 零配置 (开箱即用)

只需声明结构体，无需配置文件：
```go
// internal/gapi/start/setting.go
func LoadConfig() {
    // 默认值已配置好！
    // 服务端口: :8080
    // 数据库: sqlite ./data.db
    // 如需修改，只需覆盖：
    // ginp.SetSuccessCode(1)
    // ginp.SetFailCode(0)
}
```

代码中覆盖配置：
```go
func init() {
    ginp.SetSuccessCode(200)      // 自定义成功码
    ginp.SetFailCode(400)         // 自定义失败码
    ginp.SetSuccessMsg("OK")      // 自定义成功消息
    ginp.SetNoLoginCode(401)      // 自定义未登录码
}
```

或使用配置文件：
```yaml
# configs/config.yaml
server:
  port: 8080
database:
  type: mysql
  mysql:
    host: localhost
    port: 3306
    user: root
    password: ""
    dbname: ginp_api
```

### 🗄️ 多数据库支持
```go
db.InitMySQL()    // MySQL
db.InitPgSQL()    // PostgreSQL
db.InitSqlite()   // SQLite (默认)
```

---

## 项目结构

```
ginp-api/
├── cmd/
│   ├── gapi/            # 主程序: go run cmd/gapi/main.go
│   └── gencode/         # 代码生成: go run cmd/gencode/main.go
├── internal/gapi/
│   ├── controller/     # .a.go 文件 - HTTP 处理器
│   ├── service/        # .s.go 文件 - 业务逻辑
│   ├── model/          # .m.go 文件 - 数据库操作
│   ├── entity/         # .e.go 文件 - 数据模型
│   └── router/         # 路由注册
├── pkg/
│   ├── ginp/           # 框架扩展
│   ├── where/          # 查询构建器
│   └── utils/          # 工具函数
└── configs/            # 可选配置文件
```

## API 响应格式

| Code | Msg | 说明 |
|------|-----|------|
| 1 | success | 成功 |
| 0 | 错误信息 | 业务错误 |
| 401 | unauthorized | 未登录 |

```json
{"code": 1, "msg": "success", "data": {...}}
{"code": 0, "msg": "用户名已存在", "data": null}
```

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL / PostgreSQL / SQLite
