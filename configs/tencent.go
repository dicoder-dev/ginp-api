package configs

import "ginp-api/pkg/cfg"

// TencentCos 全局配置变量
var TencentCos = new(TencentCosConfig)

// TencentCosConfig 腾讯云COS配置
type TencentCosConfig struct {
	SecretID    string `default:""`
	SecretKey   string `default:""`
	BucketName  string `default:""`
	BucketAppID string `default:""`
	Region      string `default:""`
	Duration    int    `default:"0"`
	AllowPrefix string `default:""`
}

func init() {
	cfg.ParseConfigStruct(TencentCos)
}
