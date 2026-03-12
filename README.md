# Ginp-API

A powerful RESTful API project based by Go + Gin framework with built-in code generation tools.

## Quick Start (5 minutes)

### 1. Define Entity
```go
// internal/gapi/entity/user.e.go
package entity

type User struct {
    ID        uint
    Username  string
    Password  string
    Email     string
    Status    int
}

func (User) TableName() string { return "sys_user" }
func (User) GenConfig() *gen.EntityConfig {
    return &gen.EntityConfig{TableName: "sys_user"}
}
```

### 2. Generate CRUD Code
```bash
go run cmd/gencode/main.go
# Select: 1 → Enter: user
```

### 3. Run
```bash
go run cmd/gapi/main.go
```

---

## Complete API Example File

Generate this file with code tool:
```
internal/gapi/controller/user/cuser/sys_user_create.a.go
```

```go
package cuser

import (
    "ginp-api/internal/gapi/service/user/suser"
    "ginp-api/pkg/ginp"
)

// Request structure - auto bound from JSON
type RequestSysUserCreate struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Email    string `json:"email"`
    Status   int    `json:"status"`
}

// Response structure
type RespondSysUserCreate struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
}

// Controller handler - params auto-bound!
func SysUserCreate(c *ginp.ContextPlus, params *RequestSysUserCreate) error {
    id, err := suser.Create(params.Username, params.Password, params.Email, params.Status)
    if err != nil {
        return err  // Auto returns: {"code": 0, "msg": "error message"}
    }
    // Auto returns: {"code": 1, "msg": "success", "data": {...}}
    return c.SuccessData(RespondSysUserCreate{
        ID:       id,
        Username: params.Username,
    })
}

// Route registration (auto called via init)
func init() {
    ginp.RouterAppend(ginp.RouterItem{
        Path:        "/api/user/create",
        Handler:     ginp.BindParamsHandler(SysUserCreate, RequestSysUserCreate{}),
        HttpType:    ginp.HttpPost,
        NeedLogin:   false,
        Swagger: &ginp.SwaggerInfo{
            Title:         "Create User",
            Description:   "Create a new user account",
            RequestParams: RequestSysUserCreate{},
        },
    })
}
```

---

## Features

### 🚀 Ginp Framework Extension

**Auto Parameter Binding:**
```go
// Define struct with binding tags
type Request struct {
    Name string `json:"name" binding:"required"`
    Age  int    `json:"age"`
}

// Handler receives typed params - no manual binding needed!
func Handler(c *ginp.ContextPlus, params *Request) error {
    // params already bound from request body
    return c.SuccessData(params)
}

// Register with one line
ginp.RouterAppend(ginp.RouterItem{
    Path:    "/api/test",
    Handler: ginp.BindParamsHandler(Handler, Request{}),
    HttpType: ginp.HttpPost,
})
```

**Handy Response Methods:**
```go
c.Success()                                    // {"code": 1, "msg": "success"}
c.SuccessData(user)                           // {"code": 1, "msg": "success", "data": {...}}
c.Fail("error")                                // {"code": 0, "msg": "error"}
c.FailData("error", map[string]any{"key": 1}) // {"code": 0, "msg": "error", "data": {...}}
```

**Route with Aliases:**
```go
ginp.RouterAppend(ginp.RouterItem{
    Path:        "/api/user/info",
    AliasePaths: []string{"/api/user/profile", "/api/u/info"},
    // ...
})
```

### 🔧 Code Generation Tool
```bash
go run cmd/gencode/main.go

=== GAPI Code Generator ===
1. Generate Entity CRUD Code      # Auto create all layers
2. Add New API Controller        # Add custom API
3. Generate Field Constants     # Scan & generate FieldXXX
4. Delete Entity CRUD Code      # Remove entity files
```

### 📦 Where Query Builder
```go
// Chain query
wheres := where.New("status", "=", 1).
    And("username", "LIKE", "%admin%").
    Or("email", "=", "admin@test.com").
    Conditions()

// Use in model
users, err := model.Where(wheres).Limit(10).Find()
```

### ⚙️ Zero Config (Defaults Work Out of Box)

Just declare struct, no config files needed:
```go
// internal/gapi/start/setting.go
func LoadConfig() {
    // Default values already set!
    // Server: :8080
    // Database: sqlite ./data.db
    // Just override if needed:
    // ginp.SetSuccessCode(1)
    // ginp.SetFailCode(0)
}
```

Want to change? Override in code:
```go
func init() {
    ginp.SetSuccessCode(200)      // Custom success code
    ginp.SetFailCode(400)         // Custom fail code
    ginp.SetSuccessMsg("OK")      // Custom success message
    ginp.SetNoLoginCode(401)      // Custom no-login code
}
```

Or use config file:
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

### 🗄️ Multi-Database
```go
db.InitMySQL()    // MySQL
db.InitPgSQL()    // PostgreSQL
db.InitSqlite()   // SQLite (default)
```

---

## Project Structure

```
ginp-api/
├── cmd/
│   ├── gapi/            # Main: go run cmd/gapi/main.go
│   └── gencode/        # Generator: go run cmd/gencode/main.go
├── internal/gapi/
│   ├── controller/     # .a.go files - HTTP handlers
│   ├── service/        # .s.go files - business logic
│   ├── model/          # .m.go files - DB operations
│   ├── entity/         # .e.go files - data models
│   └── router/         # route registration
├── pkg/
│   ├── ginp/           # framework extension
│   ├── where/          # query builder
│   └── utils/          # utilities
└── configs/            # optional config files
```

## API Response Format

| Code | Msg | Description |
|------|-----|-------------|
| 1 | success | Success |
| 0 | fail message | Business error |
| 401 | unauthorized | Not logged in |

```json
{"code": 1, "msg": "success", "data": {...}}
{"code": 0, "msg": "username already exists", "data": null}
```

## Tech Stack

- **Web**: Gin
- **ORM**: GORM
- **DB**: MySQL / PostgreSQL / SQLite
