package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type File struct {
	MdPath   string `json:"md_path"`
	JsonPath string `json:"json_path"`
	WebPort  string `json:"web_port"`
}

type RedisProps struct {
	Addr    string `json:"addr"`
	Passwd  string `json:"passwd"`
	Timeout int64  `json:"timeout"`
}
type RateLimit struct {
	Retry int
	Limit int
	Burst int
}
type Config struct {
	Files     map[string]File
	Cache     map[string]RedisProps
	RateLimit RateLimit
}

func New(configFile string) *Config {

	if configFile == "" {
		panic("配置文件不能为空")
	}

	exists, err := FileExists(configFile)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic(fmt.Errorf("配置文件 : %s 不存在", configFile))
	}

	// 实例化配置文件加载器
	configLoader := viper.New()
	// 设置配置文件路径
	configLoader.SetConfigFile(configFile)
	// 设置配置文件格式
	configLoader.SetConfigType(filepath.Ext(configFile)[1:])
	// 读取配置文件
	err = configLoader.ReadInConfig()
	if err != nil {
		panic(err)
	}
	conf := &Config{}
	// 解析配置文件信息到 config 结构体
	if err := configLoader.Unmarshal(&conf); err != nil {
		panic(err)
	}
	return conf
}

func FileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
