package global

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {

}

func YamlConfig(filePath string) *Config{
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