package global

import (
	"os"
	"path/filepath"

	"gin-boot/utils"
)

var (
	AppConfig *Config
)

func Bootstrap() {
	AppConfig = utils.YamlConfig(filepath.Dir(os.Args[0])+"/conf/application.yml")
}
