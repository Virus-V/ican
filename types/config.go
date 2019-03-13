package types

import (
	"github.com/google/uuid"
)

// SrvConfigUUID 配置文件服务uuid
var SrvConfigUUID = uuid.UUID{0xa8, 0xec, 0x91, 0x39, 0x89, 0xc7, 0x42, 0x83, 0xb2, 0x3b, 0x62, 0x79, 0xc1, 0xec, 0x1e, 0x94}

// ConfigService 配置读取服务
type ConfigService interface {
	GetConfig(string) interface{}   // 获取指定的全局配置项
	ReadConfig(string, interface{}) // 读取配置文件
}
