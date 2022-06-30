package global

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/HaHadaxigua/conversion"
	"github.com/HaHadaxigua/surtr/util"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var Conf Config

func init() {
	confFilepath := "conf.yaml"
	confFile, err := os.ReadFile(confFilepath)
	if err != nil {
		logrus.Fatalf("failed to load config file: %s, %v", confFilepath, err)
	}

	if err = yaml.Unmarshal(confFile, &Conf); err != nil {
		logrus.Fatalf("failed to parse config file: %s, %v", confFilepath, err)
	}

	if Conf.Port == nil {
		Conf.Port = conversion.Intptr(HttpPort)
	}
	Conf.Domain = fmt.Sprintf(":%d", Conf.Port)
	if Conf.HomeDir == nil {
		current, err := user.Current()
		if err != nil {
			logrus.Fatalf("failed to load home dir")
		}
		Conf.HomeDir = &current.HomeDir
	}
	Conf.Storage = filepath.Join(*Conf.HomeDir, StoragePath)

	logrus.Infof("start init storage path: %s", Conf.Storage)
	if err = util.CreateDirIfNotExist(Conf.Storage); err != nil {
		logrus.Fatalf("failed to init storage path")
	}
}

type Config struct {
	httpConfig `yaml:"http_config"`
	HomeDir    *string `json:"homeDir" yaml:"home_dir"`
	Storage    string
}

type httpConfig struct {
	Domain string
	Port   *int `json:"port" yaml:"port"`
}
