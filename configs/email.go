package configs

import "ginp-api/pkg/cfg"

// ClientConfig 客户端配置（嵌套结构体示例）
type ClientConfig struct {
	Account string `default:"dicoder@126.com"`
	Pwd     string `default:"12345"`
	Port    int    `default:"465"`
	Host    string `default:"smtp.126.com"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	Client ClientConfig
}

// Email 全局配置变量
var Email = new(EmailConfig)

func init() {
	cfg.ParseConfigStruct(Email)
}
