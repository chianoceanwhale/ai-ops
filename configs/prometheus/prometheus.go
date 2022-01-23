package prometheus

import (
	"flag"

	"github.com/prometheus/client_golang/api"
	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	PrometheusClient apiv1.API
	PrometheusConfig string
	PrometheusURL    string
)

func InitPrometheusConfig() {
	prometheusAddr := flag.String("prometheusConfig", PrometheusConfig, "absolute path to the prometheusConfig file")
	client, err := api.NewClient(api.Config{
		Address: *prometheusAddr,
	})

	if err != nil {
		log.Panicf("InitPrometheusConfig addr:%s err:%v", *prometheusAddr, err)
	}

	PrometheusClient = apiv1.NewAPI(client)
}

//init promethues
func InitPrometheus(cfg *viper.Viper) {
	PrometheusURL = cfg.GetString("url")
}
