package common

import "github.com/spf13/viper"

type Application struct {
	ReadTimeout   int
	WriterTimeout int
	Host          string
	Port          string
	Name          string
	Mode          string
	Domain        string
}

func InitApplication(cfg *viper.Viper) *Application {
	return &Application{
		ReadTimeout:   cfg.GetInt("readTimeout"),
		WriterTimeout: cfg.GetInt("writerTimeout"),
		Host:          cfg.GetString("host"),
		Port:          portDefault(cfg),
		Name:          cfg.GetString("name"),
		Mode:          cfg.GetString("mode"),
		Domain:        cfg.GetString("domain"),
	}
}

var ApplicationConfig = new(Application)

func portDefault(cfg *viper.Viper) string {
	if cfg.GetString("port") == "" {
		return "8000"
	} else {
		return cfg.GetString("port")
	}
}
