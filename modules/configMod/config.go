package main

import (
	"ican/types"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// 配置文件
type config map[string]interface{}

var _ types.ConfigService = (config)(nil)

// GetConfig 获得配置内容信息
func (c config) GetConfig(item string) interface{} {
	// TODO 通过反射
	return c[item]
}

// ReadConfig 读取指定的配置项
func (c config) ReadConfig(name string, data interface{}) {
	configContent, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configContent, data)
	if err != nil {
		panic(err)
	}
}
