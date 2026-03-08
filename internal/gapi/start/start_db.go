package start

import (
	"fmt"
	"ginp-api/internal/db/dbs"
	"ginp-api/pkg/cfg"
)

func startDB() {
	// 根据配置文件中的 system.db.type 决定使用哪种数据库
	dbType := cfg.GetString("system.db.type")

	switch dbType {
	case "pgsql", "postgresql":
		dbs.InitDb(dbs.DbTypePgsql)
	case "mysql":
		dbs.InitDb(dbs.DbTypeMysql)
	case "sqlite":
		dbs.InitDb(dbs.DbTypeSqlite)
	default:
		// 默认使用 MySQL
		dbs.InitDb(dbs.DbTypeMysql)
	}

	//迁移表
	if dbs.GetWriteDb() != nil {
		//自动迁移表结构
		err := dbs.GetWriteDb().AutoMigrate(EntityAutoMigrateList...)
		if err != nil {
			fmt.Println("迁移表结构失败" + err.Error())
			panic(err)
		}
	}
}
