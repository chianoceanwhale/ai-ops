package admin

import (
	"ai-ops/configs/agent"
	"ai-ops/configs/db"
	"ai-ops/configs/prometheus"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	cfgDatabase    *viper.Viper
	cfgApplication *viper.Viper
	cfgLog         *viper.Viper
	cfgDataLog     *viper.Viper
	cfgPrometheus  *viper.Viper
	cfgCronTask    *viper.Viper
	cfDataStore    *viper.Viper
)

//载入配置文件
func ConfigSetup(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	cfgDatabase = viper.Sub("settings.database")
	if cfgDatabase == nil {
		panic("config not found settings.database")
	}
	db.DatabaseConfig = db.InitDatabase(cfgDatabase)

	cfgApplication = viper.Sub("settings.application")
	if cfgApplication == nil {
		panic("config not found settings.application")
	}
	ApplicationConfig = InitApplication(cfgApplication)

	cfgLog = viper.Sub("settings.log")
	if cfgLog == nil {
		panic("config not found settings.log")
	}
	LogConfig = InitLog(cfgLog)

	cfgPrometheus = viper.Sub("settings.prometheus")
	if cfgPrometheus == nil {
		panic("config not found settings.prometheus")
	}
	prometheus.InitPrometheus(cfgPrometheus)

	//grpc config
	agent.InitAgent(viper.Sub("settings.agent"))

}

var (
	cfgAlert   *viper.Viper
	ServerHost string
)

//载入Agent配置
func ConfigSetupAgent(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	cfgApplication = viper.Sub("settings.application")
	if cfgApplication == nil {
		panic("config not found settings.application")
	}
	ApplicationConfig = InitApplication(cfgApplication)

	cfgLog = viper.Sub("settings.log")
	if cfgLog == nil {
		panic("config not found settings.log")
	}
	LogConfig = InitLog(cfgLog)

	cfgAlert = viper.Sub("settings.alert")
	if cfgAlert == nil {
		panic("config not found settings.alert")
	}
	alert.AlertConfig = alert.ConfigAlert(cfgAlert)

	agent.K8sConfigPath = viper.GetString("settings.k8sconfig")
	prometheus.PrometheusConfig = viper.GetString("settings.prometheus")

	cfgDataLog = viper.Sub("settings.datadir")
	if cfgDataLog == nil {
		panic("config not found settings.datadir")
	}
	agent.DataConfig = agent.InitLog(cfgDataLog)

	ServerHost = viper.GetString("settings.server-host")
}

//monitor config load
func ConfigSetupMonitor(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	cfgApplication = viper.Sub("settings.application")
	if cfgApplication == nil {
		panic("config not found settings.application")
	}
	ApplicationConfig = InitApplication(cfgApplication)

	cfgDatabase = viper.Sub("settings.database")
	if cfgDatabase == nil {
		panic("config not found settings.database")
	}
	db.DatabaseConfig = db.InitDatabase(cfgDatabase)

	cfgLog = viper.Sub("settings.log")
	if cfgLog == nil {
		panic("config not found settings.log")
	}
	LogConfig = InitLog(cfgLog)
}

func SetConfig(configPath string, key string, value interface{}) {
	viper.AddConfigPath(configPath)
	viper.Set(key, value)
	viper.WriteConfig()
}
