package main

import (
	"ican/types"
	"log"
	"vgof/core"
)

type zaplogMod struct {
	obj *zaplog
}

// ModuleEntry 模块入口点
var ModuleEntry = zaplogMod{}

var _ core.Module = (*zaplogMod)(nil)

func (t *zaplogMod) CheckDepend(s core.Service) bool {
	// 该库依赖配置服务
	return s.CheckServices(types.SrvConfigUUID)
}

// Start 实现module的初始化，如果返回true，则表示安装了新的service
func (t *zaplogMod) Start(s core.Service) bool {
	cSrv, err := s.LocateService(types.SrvConfigUUID)
	if err != nil {
		panic(err)
	}
	configSrv := cSrv.(types.ConfigService)
	t.obj = &zaplog{}
	// 获得log的配置信息
	t.obj.config = configSrv.GetConfig("Log")
	err = s.InstallService(types.SrvZapLogUUID, t.obj)
	if err != nil { // 如果出错
		panic(err)
	}
	return true
}

// Stop 关闭插件
func (t *zaplogMod) Stop(s core.Service) {
	log.Print("Zap module stoped!")
	// 同步日志缓冲区
	t.obj.logger.Sync()
}

func main() {
	log.Fatal("This is a vgof module, please build this package with \"-buildmode=plugin\".")
}
