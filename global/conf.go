package global

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/HaHadaxigua/surtr/util"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Ip      string `yaml:"ip"`
	Port    int    `yaml:"port"`
	HomeDir string `yaml:"home_dir"`
	Addr    string
	Storage string
}

func NewConfig(path string) (*Config, error) {
	if util.IsFileNotExist(path) {
		logrus.Info("not found config file, will use default config")
		return defaultConfig(), nil
	}

	buf, err := os.ReadFile(path)
	if err != nil {
		logrus.Errorf("failed to load config file: %s", path)
		return nil, err
	}

	var conf Config
	if err = yaml.Unmarshal(buf, &conf); err != nil {
		logrus.Errorf("failed to parse config file: %s", path)
		return nil, err
	}

	for _, fn := range conf.setters() {
		fn()
	}

	return &conf, nil
}

func defaultConfig() *Config {
	conf := &Config{
		Port: HttpPort,
	}
	for _, fn := range conf.setters() {
		fn()
	}
	return conf
}

func (c *Config) setters() []func() {
	return []func(){
		c.setPort,
		c.setDomain,
		c.setHomeDir,
		c.setStorage,
		c.setIp,
	}
}

func (c *Config) setIp() {
	if c.Ip == "" {
		c.Ip = "localhost"
	}
	setApiAddr(fmt.Sprintf("http://%s:%d/api/", c.Ip, c.Port))
}

func (c *Config) setPort() {
	if c.Port != 0 {
		return
	}
	c.Port = HttpPort
}

func (c *Config) setHomeDir() {
	if c.HomeDir != "" {
		return
	}
	current, err := user.Current()
	if err != nil {
		logrus.Fatalf("failed to load home dir %v", err)
	}
	c.HomeDir = current.HomeDir
}

func (c *Config) setDomain() {
	if c.Port == 0 {
		logrus.Error("invalid port")
		return
	}
	c.Addr = fmt.Sprintf(":%d", c.Port)
}

func (c *Config) setStorage() {
	if c.HomeDir == "" {
		logrus.Error("invalid config, homedir have not been set")
		return
	}
	c.Storage = filepath.Join(c.HomeDir, StoragePath)
	if err := util.CreateDirIfNotExist(c.Storage); err != nil {
		logrus.Errorf("failed to create setup home on %s, %v", c.Storage, err)
		return
	}
	logrus.Infof("start init storage path: %s", c.Storage)
	setStoragePath(c.Storage)
}
