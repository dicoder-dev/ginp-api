# Ginp 设计文档

## 概述

Ginp 是基于 Gin 框架的自定义增强包（GIN Plus），位于 `pkg/ginp` 目录。它对 Gin 进行了扩展封装，提供了更便捷的 API 响应处理、统一的路由管理、参数绑定自动化等功能，让开发者可以更高效地编写 RESTful API。

## 目录结构

```
pkg/ginp/
├── bind_handler.go      # 参数绑定处理器
├── code.go              # 响应码和消息配置
├── ctx.go               # ContextPlus 核心扩展
├── func.go              # 工具函数
├── middleware.go        # 内置中间件
├── operation_type.go    # 操作类型定义
├── request_binding.go   # 请求参数绑定
├── response.go          # API 响应结构体
├── router_manager.go    # 路由管理器
└── main_test.go         # 测试示例
```

## 核心组件

### 1. ContextPlus 上下文扩展

`ContextPlus` 是对 `gin.Context` 的增强封装，提供了便捷的响应方法：

```go
type ContextPlus struct {
    *gin.Context
}
```

#### 响应方法

| 方法 | 说明 | 返回示例 |
|------|------|----------|
| `Success(messages...string)` | 返回成功响应 | `{"code": 1, "msg": "success"}` |
| `SuccessData(data any, messages...string)` | 返回带数据的成功响应 | `{"code": 1, "msg": "success", "data": {...}}` |
| `Fail(strs...string)` | 返回失败响应 | `{"code": 0, "msg": "fail"}` |
| `FailData(data any, messages...string)` | 返回带数据的失败响应 | `{"code": 0, "msg": "...", "data": {...}}` |
| `R(code int, obj any)` | 返回自定义 JSON 响应 | 自定义 |
| `Log(data any)` | 记录请求日志 | - |
| `SuccessHtml(path string)` | 返回 HTML 页面 | - |

#### 用户认证方法

```go
// 获取当前登录用户 ID
func (c *ContextPlus) GetUserID() uint

// 获取 JWT Claims
func (c *ContextPlus) getJWTClaims() map[string]interface{}
```

#### 获取路由列表

```go
// 获取所有已注册的路由
func (c *ContextPlus) GetApiList() []RouterItem
```

---

### 2. RouterItem 路由项

`RouterItem` 是路由的基本单元，包含了路由的完整配置信息：

```go
type RouterItem struct {
    Path           string              // API 路径
    AliasePaths    []string            // 路径别名列表，用于兼容旧接口
    Handler        gin.HandlerFunc     // 处理函数
    Middlewares    []gin.HandlerFunc   // 中间件链
    HttpType       string              // HTTP 方法 (GET/POST/PUT/PATCH/ANY)
    NeedLogin      bool                // 是否需要登录
    NeedPermission bool                // 是否需要权限验证
    PermissionName string              // 权限名称
    OperationType  OperationType      // 操作类型
    Swagger        *SwaggerInfo       // Swagger 文档信息
    ParamTypes     []interface{}       // 参数类型元数据
}
```

#### HTTP 方法常量

```go
const (
    HttpPost  = "POST"
    HttpGet   = "GET"
    HttpPut   = "PUT"
    HttpPatch = "PATCH"
    HttpAny   = "ANY"
)
```

---

### 3. SwaggerInfo Swagger 信息

用于生成 API 文档的元信息：

```go
type SwaggerInfo struct {
    Title        string   // 接口标题
    Description  string   // 接口描述
    RequestParams any     // 请求参数结构体（不要传入指针）
    Consumes    []string  // 请求 Content-Type，默认 ["application/json"]
    Produces    []string  // 响应 Content-Type，默认 ["application/json"]
    IsIgnore    bool      // 是否忽略该接口
}
```

---

### 4. OperationType 操作类型

定义了常见的操作类型，用于权限控制和审计：

```go
const (
    // CRUD 基础操作
    OpCreate   = "CREATE"   // 新增/创建
    OpRead     = "READ"     // 查询/读取
    OpUpdate   = "UPDATE"   // 修改/更新
    OpDelete   = "DELETE"   // 删除
    OpSearch   = "SEARCH"   // 搜索/列表查询

    // 其他常见操作
    OpImport   = "IMPORT"   // 导入
    OpExport   = "EXPORT"   // 导出
    OpDownload = "DOWNLOAD" // 下载
    OpUpload   = "UPLOAD"   // 上传
    OpSync     = "SYNC"     // 同步
    OpAudit    = "AUDIT"    // 审核
    OpApprove  = "APPROVE"  // 批准
    OpReject   = "REJECT"   // 拒绝
    OpCancel   = "CANCEL"   // 取消

    OpOther       = "OTHER"        // 其他
    OpUserCustom  = "USER_CUSTOM"  // 用户自定义
)
```

---

## 核心功能

### 1. 响应码配置

Ginp 支持自定义响应码和消息：

```go
// 初始化时设置默认值
func init() {
    SetFailCode(0)       // 失败码默认为 0
    SetSuccessCode(1)    // 成功码默认为 1
    SetNoLoginCode(401)  // 未登录码默认为 401
    codeHttpSuccess = http.StatusOK
    codeHttpFail = http.StatusOK
}

// 自定义配置
SetSuccessCode(100)        // 设置成功码
SetFailCode(0)             // 设置失败码
SetNoLoginCode(401)        // 设置未登录码
SetSuccessMsg("操作成功")   // 设置成功消息
SetFailMsg("操作失败")      // 设置失败消息
SetShowLog(true)           // 开启/关闭日志
SetLogLevel(LogLevelDebug) // 设置日志级别
```

---

### 2. 参数绑定处理器

Ginp 提供了强大的参数绑定功能，可以自动将请求参数绑定到结构体：

#### BindHandler 基础绑定

```go
// 将 gin.HandlerFunc 转换为自定义处理函数
func BindHandler(handler func(c *ContextPlus)) func(c *gin.Context)
```

#### BindParamsHandler 自动参数绑定

这是最强大的功能，可以自动从请求中绑定参数到结构体：

```go
func BindParamsHandler(handler interface{}, paramTypes ...interface{}) gin.HandlerFunc
```

**使用示例**：

```go
// 定义请求参数结构体
type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Email    string `json:"email"`
}

// 控制器函数签名：最后一个参数可以是 error 类型
func CreateUser(c *ginp.ContextPlus, params *CreateUserRequest) error {
    // params 已经自动从请求体绑定
    user, err := service.CreateUser(params.Username, params.Password, params.Email)
    if err != nil {
        return err  // 自动返回失败响应
    }
    return nil  // 自动返回成功响应
}

// 注册路由时使用
ginp.RouterAppend(ginp.RouterItem{
    Path:     "/api/user/create",
    Handler:  ginp.BindParamsHandler(CreateUser, CreateUserRequest{}),
    HttpType: ginp.HttpPost,
    // ...
})
```

**参数绑定规则**：

| 请求方法 | 绑定方式 |
|----------|----------|
| POST/PUT/PATCH | `ShouldBindJSON` |
| GET | `ShouldBindQuery` |

**返回值处理**：

- 如果 handler 返回 error，自动调用 `ctx.Fail(err.Error())`
- 如果 handler 中已经调用了 `ctx.Success/SuccessData/Fail` 等方法，不会重复响应
- 如果没有显式响应，默认调用 `ctx.Success()`

---

### 3. 请求参数绑定

除了自动绑定，还提供了手动绑定函数：

```go
// 从 JSON 绑定
result := ginp.BindJSON(ctx, &params)
if !result.Success {
    ctx.Fail(result.Message)
    return
}

// 从 Query 参数绑定
result := ginp.BindQuery(ctx, &params)

// 从 URL 路径参数绑定
result := ginp.BindURI(ctx, &params)

// 强制绑定（失败自动返回错误）
ginp.MustBindJSON(ctx, &params)
ginp.MustBindQuery(ctx, &params)
ginp.MustBindURI(ctx, &params)
```

**BindResult 结构**：

```go
type BindResult struct {
    Success bool        // 是否成功
    Data    interface{} // 绑定后的数据
    Error   error       // 原始错误
    Message string      // 用户友好的错误消息
}
```

---

### 4. 内置中间件

Ginp 提供了多个内置中间件：

#### LoggingMiddleware 请求日志

```go
// 记录请求方法、路径、状态码和耗时
func LoggingMiddleware() gin.HandlerFunc
```

#### CORSMiddleware 跨域

```go
// 配置跨域允许的方法、头部和来源
func CORSMiddleware() gin.HandlerFunc
```

#### RecoveryMiddleware 异常恢复

```go
// 捕获 panic 异常，返回 500 错误
func RecoveryMiddleware() gin.HandlerFunc
```

#### RequestIDMiddleware 请求 ID

```go
// 为每个请求生成唯一 ID，便于日志追踪
func RequestIDMiddleware() gin.HandlerFunc
```

---

### 5. API 响应结构体

标准化的 API 响应体：

```go
type ApiResponse struct {
    Code    interface{} `json:"code"`
    Msg     string      `json:"msg"`
    Data    interface{} `json:"data,omitempty"`
    Payload interface{} `json:"payload,omitempty"`
}

// 创建响应
ginp.NewSuccessResponse(data)
ginp.NewFailResponse("错误信息")
ginp.NewFailResponseWithData(data, "错误信息")
```

---

### 6. 路由注册与管理

#### 注册路由

```go
// 向路由池添加路由
ginp.RouterAppend(ginp.RouterItem{
    Path:        "/api/user/login",
    Handler:     ginp.BindParamsHandler(Login, LoginRequest{}),
    HttpType:    ginp.HttpPost,
    NeedLogin:   false,
    PermissionName: "user.login",
    Swagger: &ginp.SwaggerInfo{
        Title:         "用户登录",
        Description:   "通过用户名密码登录",
        RequestParams: LoginRequest{},
    },
})

// 注册到 Gin Engine
ginp.RegisterRouter(r *gin.Engine)

// 获取所有路由
routers := ginp.GetAllRouter()
```

#### 路径别名

支持为同一接口注册多个路径：

```go
ginp.RouterAppend(ginp.RouterItem{
    Path:        "/api/user/info",
    AliasePaths: []string{"/api/user/profile", "/api/u/info"}, // 别名
    Handler:     ginp.BindParamsHandler(GetUserInfo, UserInfoRequest{}),
    HttpType:    ginp.HttpGet,
    // ...
})
```

#### 路由分组

```go
type RouterGroup struct {
    Prefix      string
    Items       []RouterItem
    Middlewares []gin.HandlerFunc
}
```

---

## 项目集成

### 1. 路由初始化流程

```
┌─────────────────────┐
│   main.go           │
│   startGinServer()   │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│   gin.Default()      │
│   创建 Gin 实例       │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│   router.Register() │
│   注册中间件和路由    │
└─────────┬───────────┘
          │
    ┌─────┴─────┐
    │           │
    ▼           ▼
┌────────┐  ┌────────────┐
│ CORS   │  │ ginp.RegisterRouter() │
│中间件  │  │ 注册所有路由 │
└────────┘  └────────────┘
```

### 2. 控制器注册

在 `internal/gapi/controller` 目录下，每个控制器文件使用 `init()` 函数注册路由：

```go
// 文件：internal/gapi/controller/user/cuser/login_by_username.a.go

package cuser

import (
    "ginp-api/internal/gapi/service/user/suser"
    "ginp-api/pkg/ginp"
)

func init() {
    ginp.RouterAppend(ginp.RouterItem{
        Path:           "/api/sys_user/login_by_username",
        Handler:        ginp.BindParamsHandler(LoginByUsername, RequestLoginByUsername{}),
        HttpType:       ginp.HttpPost,
        NeedLogin:      false,
        NeedPermission: false,
        PermissionName: "User.login_by_username",
        Swagger: &ginp.SwaggerInfo{
            Title:         "login_by_username",
            Description:   "用户名登录",
            RequestParams: RequestLoginByUsername{},
        },
    })
}

func LoginByUsername(c *ginp.ContextPlus, requestParams *RequestLoginByUsername) {
    userInfo, token, err := suser.LoginByUsername(requestParams.Username, requestParams.Password)
    if err != nil {
        c.FailData(err.Error())
        return
    }
    c.SuccessData(&RespondLogin{
        Token:    token,
        UserInfo: userInfo,
    })
}

type RequestLoginByUsername struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RespondLogin struct {
    Token    string      `json:"token"`
    UserInfo interface{} `json:"user_info"`
}
```

### 3. 路由导入聚合

通过 `routers_import.go` 聚合所有控制器的导入：

```go
// 文件：internal/gapi/router/routers_import.go

package router

// 显式导入确保 init 函数被调用
import (
    _ "ginp-api/internal/gapi/controller/system/cindex"
    _ "ginp-api/internal/gapi/controller/user/cuser"
    // {{placeholder_router_import}}//
)
```

### 4. 路由注册主函数

```go
// 文件：internal/gapi/router/router.go

package router

import (
    "ginp-api/pkg/ginp"
    "github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
    // 1. 中间件配置
    r.Use(CORSMiddleware())  // 跨域

    // 2. 路由注册
    ginp.RegisterRouter(r)
}
```

---

## 最佳实践

### 1. 控制器函数签名

推荐使用以下签名：

```go
// 简单场景：不处理错误
func Handler(c *ginp.ContextPlus, params *RequestStruct) {
    // 处理业务逻辑
    c.SuccessData(result)
}

// 错误处理场景：返回 error
func Handler(c *ginp.ContextPlus, params *RequestStruct) error {
    err := service.DoSomething(params)
    if err != nil {
        return err  // 自动返回失败响应
    }
    return nil  // 自动返回成功响应
}
```

### 2. 命名规范

| 类型 | 命名规则 | 示例 |
|------|----------|------|
| 请求结构体 | `Request{API名称}` | `RequestLogin` |
| 响应结构体 | `Respond{API名称}` | `RespondLogin` |
| 路由路径 | `/api/{资源}/{操作}` | `/api/user/login` |
| 权限名称 | `{模块}.{资源}.{操作}` | `user.login` |

### 3. Swagger 文档

每个接口都应配置 Swagger 信息：

```go
Swagger: &ginp.SwaggerInfo{
    Title:         "用户登录",
    Description:   "通过用户名密码登录系统",
    RequestParams: RequestLoginByUsername{},  // 注意：不使用指针
    Consumes:      []string{"application/json"},
    Produces:      []string{"application/json"},
}
```

---

## 与 Gin 原生对比

| 特性 | Gin 原生 | Ginp |
|------|----------|------|
| 响应 | `c.JSON(200, gin.H{...})` | `c.SuccessData(data)` |
| 参数绑定 | 手动 `c.ShouldBindJSON(&params)` | 自动 `BindParamsHandler` |
| 路由管理 | 分散注册 | 集中管理、支持并发安全 |
| 统一响应格式 | 手动封装 | 内置支持 |
| 日志 | 手动配置 | 内置中间件 |

---

## 扩展阅读

- [代码生成逻辑](./代码生成逻辑.md) - 了解如何使用代码生成工具生成基于 ginp 的 CRUD 代码
- [项目结构](./项目结构.md) - 了解项目整体架构
- [配置文件](./配置文件.md) - 了解项目配置
