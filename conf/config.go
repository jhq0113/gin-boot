package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	Conf *Config
)

func init() {
	Conf = YamlConfig(filepath.Dir(os.Args[0]) + "/conf/application.yml")
}

type Config struct {
	App    App                      `yaml:"app" json:"app"`
	Redis  map[string][]RedisOption `yaml:"redis" json:"redis"`
	Params map[string]interface{}   `yaml:"params" json:"params"`
}

type App struct {
	Addr string `yaml:"addr" json:"addr"`
	//可以是debug、test、release
	Mode string `yaml:"mode" json:"mode"`
}

type RedisOption struct {
	Host string `yaml:"host" json:"host"`
	Port string `yaml:"port" json:"port"`
	Auth string `yaml:"auth" json:"auth"`
	Db   uint8  `yaml:"db" json:"db"`
	//单位s
	MaxConnLifetime time.Duration `yaml:"db" json:"db"`
	MaxIdle         int           `yaml:"maxIdle" json:"maxIdle"`
	MaxActive       int           `yaml:"maxActive" json:"maxActive"`
	Wait            bool          `yaml:"wait" json:"wait"`
	//单位ms
	ConnectTimeout int `yaml:"connectTimeout" json:"connectTimeout"`
	//单位ms
	ReadTimeout int `yaml:"readTimeout" json:"readTimeout"`
	//单位ms
	WriteTimeout int `yaml:"writeTimeout" json:"writeTimeout"`
}

func YamlConfig(filePath string) *Config {
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}

	config := &Config{}
	err = yaml.Unmarshal(conf, config)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func JsonConfig(filePath string) *Config {
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}

	config := &Config{}
	err = json.Unmarshal(conf, config)
	if err != nil {
		panic(err.Error())
	}

	return config
}
