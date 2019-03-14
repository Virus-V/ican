package main

import (
	"flag"
	"ican/types"
	"io/ioutil"
	"log"
	"vgof/core"

	yaml "gopkg.in/yaml.v2"
)

type configModule struct {
	// 配置文件路径
	configFile string
}

// ModuleEntry 模块入口点
var ModuleEntry = configModule{}

var _ core.Module = (*configModule)(nil)

func (t *configModule) CheckDepend(s core.Service) bool {
	// 该模块无依赖
	return true
}

// Start 实现module的初始化，如果返回true，则表示安装了新的service
func (t *configModule) Start(s core.Service) bool {
	flag.StringVar(&t.configFile, "config", "conf/config.yaml", "Config file")
	flag.Parse() // 解析参数
	obj := make(config)
	configContent, err := ioutil.ReadFile(t.configFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configContent, obj)
	if err != nil {
		panic(err)
	}
	// 安装Application服务
	err = s.InstallService(types.SrvConfigUUID, obj)
	if err != nil { // 如果出错
		panic(err)
	}
	return true
}

// Stop 关闭插件
func (t *configModule) Stop(s core.Service) {
	log.Print("Config module stoped!")
}

// String ...
func (t *configModule) String() string {
	return "Config Module"
}

func main() {
	log.Fatal("This is a vgof module, please build this package with \"-buildmode=plugin\".")
}
