package utils

import (
	"encoding/json"
	"io/ioutil"

	"gin-boot/global"

	"gopkg.in/yaml.v2"
)

func YamlConfig(filePath string) *global.Config{
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}

	config := &global.Config{}
	err = yaml.Unmarshal(conf, config)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func JsonConfig(filePath string) *global.Config {
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}

	config := &global.Config{}
	err = json.Unmarshal(conf, config)
	if err != nil {
		panic(err.Error())
	}

	return config
}
