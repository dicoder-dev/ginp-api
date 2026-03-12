# Ginp-API

A powerful RESTful API project based on Go + Gin framework with built-in code generation tools.

## Features

### 🚀 Ginp Framework Extension
Enhanced Gin framework with convenient APIs:

**ContextPlus - Handy Response Methods:**
```go
// Success response
c.Success()

// Success with data
c.SuccessData(&User{ID: 1, Name: "John"})

// Fail response
c.Fail("Invalid parameter")

// Fail with data
c.FailData(err.Error(), map[string]any{"field": "username"})
```

**Auto Parameter Binding:**
```go
// Define request struct
type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Email    string `json:"email"`
}

// Controller - params auto-bound!
func CreateUser(c *ginp.ContextPlus, params *CreateUserRequest) error {
    user, err := service.CreateUser(params.Username, params.Password)
    if err != nil {
        return err  // Auto returns fail response
    }
    return nil  // Auto returns success response
}

// Register route
ginp.RouterAppend(ginp.RouterItem{
    Path:    "/api/user/create",
    Handler: ginp.BindParamsHandler(CreateUser, CreateUserRequest{}),
    HttpType: ginp.HttpPost,
})
```

**Route Management with Aliases:**
```go
ginp.RouterAppend(ginp.RouterItem{
    Path:        "/api/user/info",
    AliasePaths: []string{"/api/user/profile", "/api/u/info"},
    Handler:     ginp.BindParamsHandler(GetUserInfo, Request{}),
    HttpType:    ginp.HttpGet,
})
```

### 🔧 Code Generation Tool
Automated CRUD code generation with interactive menu:

```bash
go run cmd/gencode/main.go

=== GAPI Code Generator ===
1. Generate Entity CRUD Code      # Auto generate Controller/Service/Model
2. Add New API Controller         # Add custom API to existing module
3. Generate Field Constants      # Scan entities, generate FieldXXX constants
4. Delete Entity CRUD Code        # Remove all CRUD files for entity
```

### 📦 Where Query Builder
Elegant chain-style query condition builder:

```go
// Simple condition
wheres := where.New("id", "=", 1).Conditions()

// Chain with AND
wheres := where.New("status", "=", 1).
    And("username", "=", "admin").
    And("email", "LIKE", "%@example.com").
    Conditions()

// Chain with OR
wheres := where.New("status", "=", 1).
    Or("status", "=", 2).
    Conditions()

// With operators
wheres := where.New("age", ">", 18).
    And("score", ">=", 60).
    And("status", "IN", []int{1,2,3}).
    Conditions()
```

### 🗄️ Multi-Database Support
```go
// MySQL
db.InitMySQL()

// PostgreSQL
db.InitPgSQL()

// SQLite
db.InitSqlite()
```

### 📁 Layered Architecture
```
Controller → Service → Model → Entity
```

## Quick Start

### 1. Create Entity
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

### 2. Generate Code
```bash
go run cmd/gencode/main.go
# Select option 1, enter entity name "user"
```

### 3. Implement Business Logic
Edit generated controller files in `internal/gapi/controller/`

### 4. Run Server
```bash
go run cmd/gapi/main.go
```

## Project Structure

```
ginp-api/
├── cmd/
│   ├── gapi/           # Main application
│   └── gencode/        # Code generation tool
├── internal/gapi/
│   ├── controller/     # HTTP handlers (.a.go)
│   ├── service/        # Business logic (.s.go)
│   ├── model/          # Database operations (.m.go)
│   ├── entity/         # Data models (.e.go)
│   └── router/        # Route config
├── pkg/
│   ├── ginp/           # Framework extension
│   ├── where/          # Query builder
│   ├── dbops/          # DB operations
│   └── utils/          # Utilities
└── configs/            # Configuration files
```

## Configuration

Edit `configs/config.yaml`:
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

## API Response Format

```json
{
  "code": 1,
  "msg": "success",
  "data": {...}
}
```

## Tech Stack

- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL / PostgreSQL / SQLite

## License

MIT
