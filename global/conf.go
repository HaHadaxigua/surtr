package global

import (
	"os/user"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Conf Config

func init() {
	current, err := user.Current()
	if err != nil {
		logrus.Errorf("failed to load home dir")
	}

	Conf = Config{
		httpConfig: httpConfig{
			Domain: HttpDomain,
		},
		HomeDir: current.HomeDir,
		Storage: filepath.Join(current.HomeDir, StoragePath),
	}
}

type Config struct {
	httpConfig
	HomeDir string
	Storage string
}

type httpConfig struct {
	Domain string `json:"domain"`
}
