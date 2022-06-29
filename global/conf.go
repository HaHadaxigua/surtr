package global

var Conf Config

func init() {
	Conf = Config{
		httpConfig{
			Domain: HttpDomain,
		},
	}
}

type Config struct {
	httpConfig
}

type httpConfig struct {
	Domain  string `json:"domain"`
	Storage string `json:"storage"`
}
