package agent

import (
	"ai-ops/pkg/agent/protos"
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type agent struct {
	host     string
	port     string
	LabelKey string
	Kubecfg  string
}

func InitAgent(cfg *viper.Viper) {
	c := cfg.AllSettings()

	for k, v := range c {
		vmap := v.(map[string]interface{})

		AgentConfig[k] = &agent{
			host:     vmap["host"].(string),
			port:     vmap["port"].(string),
			LabelKey: vmap["label-key"].(string),
			Kubecfg:  vmap["kubecfg"].(string),
		}
	}
}

var AgentConfig = make(map[string]*agent)

//k8s 指定地区 Client
var K8sClient *kubernetes.Clientset
var K8sMetrics *metrics.Clientset
var K8sRestConfig *rest.Config
var IstioClient *versionedclient.Clientset

//grpc 不同地区Client
var K8sGrpcClientMap = make(map[string]protos.IK8SServiceClient)
var K8sGrpcClient protos.IK8SServiceClient

var AgentGrpcClientMap = make(map[string]protos.IAgentServiceClient)
var AgentGrpcClient protos.IAgentServiceClient

var K8sLabelKey string
var K8sConfigPath string

func InitKubeConfig() {
	kubeconfig := flag.String("kubeconfig", K8sConfigPath, "absolute path to the kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Panic("Build kubeconfig err", err)
	}
	K8sRestConfig = config
	K8sClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic("New kubernetes Clientset err", err)
	}

	//istio init
	IstioClient, err = versionedclient.NewForConfig(config)
	if err != nil {
		log.Panic("New kubernetes Clientset err", err)
	}

	K8sMetrics, err = metrics.NewForConfig(config)
	if err != nil {
		log.Panic("New kubernetes Metrics err", err)
	}
}

//通过地区获取k8s配置
func InitKubeByRegion(region string) (clientset *kubernetes.Clientset, restcfg *rest.Config) {
	kubecfg := AgentConfig[region].Kubecfg
	restcfg, err := clientcmd.BuildConfigFromFlags("", kubecfg)
	if err != nil {
		log.Panic("Build kubeconfig err", err)
	}

	clientset, err = kubernetes.NewForConfig(restcfg)
	if err != nil {
		log.Panic("New kubernetes Clientset err", err)
	}

	return
}

//初始化GRPC客户端
func InitK8sGrpcClient() {
	for item := range AgentConfig {
		addr := AgentConfig[item].host + ":" + AgentConfig[item].port
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatal("failed to connect : ", err)
		}
		K8sGrpcClientMap[item] = protos.NewIK8SServiceClient(conn)
		AgentGrpcClientMap[item] = protos.NewIAgentServiceClient(conn)
	}
}
