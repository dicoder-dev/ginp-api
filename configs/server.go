package configs

import "ginp-api/pkg/cfg"

// ServerConfig 服务配置
type ServerConfig struct {
	Port string `default:"8082"`
}

// Server 全局配置变量
var Server = new(ServerConfig)

func init() {
	cfg.ParseConfigStruct(Server)
}
