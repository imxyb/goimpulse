package conf

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"runtime"
)

type Config struct {
	App struct{
		Host string `json:"host"`
	} `json:"app"`

	Etcd struct {
		Host []string `json:"host"`
	} `json:"etcd"`

	Type []string `json:"type"`

	Batch int `json:"batch"`

	Auth struct{
		User string `json:"user"`
		Pass string `json:"pass"`
		Enable bool `json:"enable"`
	} `json:"auth"`

	NodeManager struct{
		Host string `json:"host"`
	} `json:"node_manager"`
}

var Cfg *Config

func init() {
	Cfg = &Config{}

	_, filename, _, _ := runtime.Caller(0)
	curPath := path.Dir(filename)
	data, err := ioutil.ReadFile(path.Join(curPath, "config.json"))
	if err != nil {
		panic("配置文件不存在")
	}

	json.Unmarshal(data, Cfg)
}
