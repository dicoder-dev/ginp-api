# GINP-API 命令行工具

## 简介

GINP-API 是一个基于 Gin 框架的 API 开发工具，提供了代码生成、实体管理等功能。

## 安装

1. 克隆仓库
```bash
git clone https://github.com/dicoder-cn/ginp-api.git
cd ginp-api
```

2. 编译命令行工具
```bash
./scripts/build_gapi.sh
```

3. 添加到 PATH（可选）
```bash
export PATH="/path/to/ginp-api/build:$PATH"
```

## 使用方法

### 交互式命令行

运行以下命令进入交互式菜单：
```bash
gapi gen
```

#### 菜单选项说明

1. **生成实体 CRUD 代码** - 创建新实体并生成 CRUD 代码，或为已存在实体生成 CRUD
2. **新增 API 接口控制器** - 添加新的 API 接口文件
3. **生成实体字段常量** - 生成实体字段常量定义
4. **删除实体 CRUD 代码** - 删除指定实体的 CRUD 代码
5. **退出** - 退出程序

#### 菜单操作流程

**生成实体 CRUD 代码**
1. 选择操作：创建新实体 或 为已存在实体生成 CRUD
2. 输入实体名称（大驼峰命名，如 UserGroup）
3. 选择父级目录（system, user 等）
4. 系统自动生成 Entity、Service、Model、Fields、Controller 文件
5. 生成的实体文件会自动记录 `FatherFolderName` 配置

**新增 API 接口控制器**
1. 选择目标目录（system, user 等）或创建新目录
2. 如果目录下有子目录（以 c 开头），继续选择子目录
3. 如果没有子目录，选择在父级目录下新建（新建的文件夹必须以 c 开头）
4. 输入 API 名称（大驼峰命名，如 GetUserInfo）
5. 系统在指定目录生成 API 文件

**删除实体 CRUD 代码**
1. 选择要删除的实体
2. 系统自动从实体配置中读取 `FatherFolderName` 定位目录
3. 确认删除操作

### 命令行参数

#### 生成 swagger 文档
```bash
cd ./cmd/gencode && go run main.go swagger
```

#### 查看版本
```bash
gapi -v
```

#### 创建实体并生成 CRUD 代码
```bash
# 交互式方式
gapi gen entity

# 直接指定实体名称和父级目录
gapi gen entity -c UserGroup -p user
```

#### 生成实体字段常量
```bash
gapi gen const
```

#### 新增 API 接口
```bash
# 交互式方式
gapi gen api

# 直接指定 API 名称和目录
gapi gen api -a GetUserInfo -d user/cuser
```

## 命令说明

### 根命令
- `gapi`: 显示帮助信息
- `gapi -v`: 显示版本信息

### 创建 UserGroup 实体
```bash
gapi gen entity -c UserGroup -p user
```

### 新增一个接口
```bash
# -a API名称,大驼峰命名法 -d指定api所在文件夹,
# 存放于controller/user/cuser文件夹，命名为get_user_info.go 采用一个api接口一个文件的方式
gapi gen api -d user/cuser -a GetUserInfo
```

### 生成实体字段常量
```bash
gapi gen const
```

### 现有实体生成 CRUD
```bash
# 指定多个实体 -p指定父级目录 -e 指定实体名称列表
gapi gen crud -e tableNameDemoTable1,tableNameDemoTable2 -p demo
# 指定一个实体 -p指定父级目录 -e 指定实体名称
gapi gen crud -e tableNameDemoTable1 -p demo
```

### 删除现有实体的 CRUD 文件
```bash
# 删除多个实体的CRUD文件 -p指定父级目录 -e 指定实体名称列表
gapi gen rm crud -e tableNameDemoTable1,tableNameDemoTable2 -p demo
# 删除一个实体的CRUD文件 -p指定父级目录 -e 指定实体名称
gapi gen rm crud -e tableNameDemoTable1 -p demo
# 交互式删除（不指定参数时进入交互模式）
gapi gen rm crud
```

## 项目结构

```
cmd/
├── gapi/           # 主应用程序入口 (main.go)
└── gencode/        # 代码生成工具入口
    ├── cmd/        # 子命令定义 (gen, rm, swagger)
    ├── desc/       # 代码生成描述逻辑
    └── templates/  # 代码生成模板 (.tmpl 文件)

internal/
├── gapi/           # 内部应用代码
│   ├── controller/ # 控制器层 (cxxx 目录)
│   ├── service/    # 业务逻辑层 (sxxx 目录)
│   ├── model/      # 数据模型层 (mxxx 目录)
│   ├── entity/     # 实体定义 (.e.go)
│   ├── dto/        # 数据传输对象
│   ├── router/     # 路由配置
│   └── start/      # 启动配置
└── db/             # 数据库初始化 (mysql/pgsql/sqlite)

pkg/                 # 可复用公共库
├── ginp/            # Gin 框架扩展 (ContextPlus)
├── dbops/           # 数据库操作
├── where/           # 查询条件构建
└── utils/           # 工具函数
```

## 代码规范

### 命名约定
- 控制器文件: `xxx.a.go` (如 `user_search.a.go`)
- 服务文件: `xxx.s.go` (如 `user.s.go`)
- 模型文件: `xxx.m.go` (如 `user.m.go`)
- 实体文件: `xxx.e.go` (如 `user.e.go`)
- 路由文件: `xxx.r.go` (如 `index.r.go`)

### API 接口目录规范
- 存放 API 接口的文件夹必须以 `c` 开头
- 例如：`cuser`、`csystem`、`cuser/test`

### 实体配置
- 实体文件中的 `GenConfig()` 方法可配置 `FatherFolderName` 字段
- 该字段用于指定实体所属的父级目录

## 贡献

欢迎提交 Pull Request 或提出 Issue。

## 许可证

[MIT](LICENSE)
