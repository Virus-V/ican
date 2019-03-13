package main

import (
	"ican/types"
	"log"
	"vgof/core"

	"go.uber.org/zap"
)

type sendmailMod struct {
	obj sendmail
}

// ModuleEntry 模块入口点
var ModuleEntry = sendmailMod{}

var _ core.Module = (*sendmailMod)(nil)

func (t *sendmailMod) CheckDepend(s core.Service) bool {
	// 该库依赖配置服务,zap日志服务
	return s.CheckServices(types.SrvConfigUUID, types.SrvZapLogUUID)
}

// Start 实现module的初始化，如果返回true，则表示安装了新的service
func (t *sendmailMod) Start(s core.Service) bool {
	t.obj = sendmail{}
	cSrv, err := s.LocateService(types.SrvConfigUUID)
	if err != nil {
		panic(err)
	}
	configSrv := cSrv.(types.ConfigService)
	zaplogSrv, err := s.LocateService(types.SrvZapLogUUID)
	if err != nil {
		panic(err)
	}
	// 安装日志插件
	t.obj.logger = ((zaplogSrv.(types.ZapLogService)).GetZapLogger(s)).(*zap.Logger)
	// 获得log的配置信息
	config := configSrv.GetConfig("Mail")
	mailCfg := config.(map[interface{}]interface{})
	t.obj.From = mailCfg["From"].(string)                 // 发件人
	t.obj.FromName = mailCfg["FromName"].(string)         // 发件人姓名
	t.obj.SMTPAddr = mailCfg["SMTPAddr"].(string)         // smtp服务器地址
	t.obj.SMTPPort = mailCfg["SMTPPort"].(int)            // smtp服务器端口
	t.obj.SMTPUsername = mailCfg["SMTPUsername"].(string) // smtp服务器用户名
	t.obj.SMTPPassword = mailCfg["SMTPPassword"].(string) // smtp服务器密码
	err = s.InstallService(types.SrvSendMailUUID, &t.obj)
	if err != nil { // 如果出错
		panic(err)
	}
	return true
}

// Stop 关闭插件
func (t *sendmailMod) Stop(s core.Service) {
	log.Print("Sendmail module stoped!")
}

func main() {
	log.Fatal("This is a vgof module, please build this package with \"-buildmode=plugin\".")
}
