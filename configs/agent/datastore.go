package agent

import "github.com/spf13/viper"

type dataDir struct {
	Dir      string
	FileName string
}

func InitLog(cfg *viper.Viper) *dataDir {
	return &dataDir{
		Dir:      cfg.GetString("dir"),
		FileName: cfg.GetString("filename"),
	}
}

var DataConfig = new(dataDir)
