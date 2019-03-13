package main

import (
	"ican/types"
	"log"
	"vgof/core"

	"go.uber.org/zap"
)

type icanModule struct {
}

// ModuleEntry 模块入口点
var ModuleEntry = icanModule{}

var _ core.Module = (*icanModule)(nil)

func (t *icanModule) CheckDepend(s core.Service) bool {
	// 该模块需要zaplog模块,config模块和sendmail模块
	return s.CheckServices(types.SrvZapLogUUID, types.SrvConfigUUID, types.SrvSendMailUUID)
}

// Start 实现module的初始化，如果返回true，则表示安装了新的service
func (t *icanModule) Start(s core.Service) bool {
	appObj := &ican{}
	zaplogSrv, err := s.LocateService(types.SrvZapLogUUID)
	if err != nil {
		panic(err)
	}
	// 安装日志插件
	appObj.logger = ((zaplogSrv.(types.ZapLogService)).GetZapLogger(s)).(*zap.Logger)
	// 获得邮件发送服务
	sSrv, err := s.LocateService(types.SrvSendMailUUID)
	if err != nil {
		panic(err)
	}
	appObj.sendMail = sSrv.(types.SendMailService)
	// 获得配置接口
	cSrv, err := s.LocateService(types.SrvConfigUUID)
	if err != nil {
		panic(err)
	}
	configSrv := cSrv.(types.ConfigService)
	config := configSrv.GetConfig("System")
	sysCfg := config.(map[interface{}]interface{})
	appObj.taskPath = sysCfg["TaskPath"].(string)
	appObj.subject = sysCfg["Subjects"].([]interface{})
	appObj.body = sysCfg["Bodys"].([]interface{})
	// 安装Application服务
	err = s.InstallService(core.SrvApplicationUUID, appObj)
	if err != nil { // 如果出错
		panic(err)
	}
	return true
}

// Stop 关闭插件
func (t *icanModule) Stop(s core.Service) {
	log.Print("Ican! module stoped!")
}

func main() {
	log.Fatal("This is a vgof module, please build this package with \"-buildmode=plugin\".")
}
