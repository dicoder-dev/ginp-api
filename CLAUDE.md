# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

GINP-API 是一个基于 Gin 框架的 API 开发工具，提供代码生成、实体管理等功能。项目包含两个主要二进制文件：
- `gapi`: 主应用程序（API 服务）
- `gencode`: 代码生成工具

## 常用命令

```bash
# 构建 gapi 主程序
./scripts/build_gapi.sh

# 生成 swagger 文档
cd ./cmd/gencode && go run main.go swagger

# 查看版本
gapi -v

# 创建实体并生成 CRUD 代码
gapi gen entity -c UserGroup -p user

# 新增 API 接口
gapi gen api -a GetUserInfo -d user/cuser

# 生成实体字段常量
gapi gen const

# 为现有实体生成 CRUD
gapi gen crud -e tableNameDemoTable1 -p demo

# 删除实体的 CRUD 文件
gapi gen rm crud -e tableNameDemoTable1 -p demo

# 运行测试
go test ./...
go test -v ./pkg/xxx # 运行指定包测试
```

## 架构概览

```
cmd/
├── gapi/           # 主应用程序入口 (main.go)
└── gencode/        # 代码生成工具入口
    ├── cmd/        # 子命令定义 (gen, rm, swagger)
    ├── desc/       # 代码生成描述逻辑
    ├── swagen/     # Swagger 生成
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
    └── dbs/        # 数据库实例管理

pkg/                 # 可复用公共库
├── ginp/            # Gin 框架扩展 (ContextPlus)
├── dbops/           # 数据库操作
├── where/           # 查询条件构建
├── utils/           # 工具函数
├── cos/             # 腾讯云 COS
├── email/           # 邮件发送
├── filehelper/      # 文件处理
├── logger/          # 日志
└── httpclient/      # HTTP 客户端
```

## 代码规范

### 命名约定
- 控制器文件: `xxx.a.go` (如 `user_search.a.go`)
- 服务文件: `xxx.s.go` (如 `user.s.go`)
- 模型文件: `xxx.m.go` (如 `user.m.go`)
- 实体文件: `xxx.e.go` (如 `user.e.go`)
- 路由文件: `xxx.r.go` (如 `index.r.go`)

### 实体字段命名
- 使用 `RequestParams` 替代 `RequestDto`
- 使用 `BindParamsHandler` 进行参数绑定

### 框架特性
- 使用 `pkg/ginp` 的 `ContextPlus` 处理响应
- 响应方法: `Success()`, `Fail()`, `SuccessData()`, `FailData()`
- 使用 GORM 作为 ORM，支持 MySQL、PostgreSQL、SQLite
