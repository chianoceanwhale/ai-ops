package admin

import (
	"ai-ops/common"
	. "ai-ops/configs"
	"ai-ops/configs/agent"
	"ai-ops/database"
	"ai-ops/internal/job"
	"ai-ops/router"
	"ai-ops/tools/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config   string
	port     string
	mode     string
	StartCmd = &cobra.Command{
		Use:     "ai-ops admin server",
		Short:   "Start API server",
		Example: "ai-ops admin server configs/settings.yml",
		PreRun: func(cmd *cobra.Command, args []string) {
			usage()
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "configs/admin/admin-settings.yml", "Start admin server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8000", "Tcp port server listening on")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func usage() {
	usageStr := `starting admin api server`
	log.Printf("%s\n", usageStr)
}

func setup() {

	//1. 读取配置
	ConfigSetup(config)
	//2. 设置日志
	InitLogger()
	//3. 初始化数据库链接
	database.Setup()
}

func run() error {
	if mode != "" {
		SetConfig(config, "settings.application.mode", mode)
	}
	if viper.GetString("settings.application.mode") == string(common.ModeProd) {
		gin.SetMode(gin.ReleaseMode)
	}

	r := router.InitRouter()

	db, _ := database.GormDB.DB()
	defer db.Close()

	if port != "" {
		SetConfig(config, "settings.application.port", port)
	}

	srv := &http.Server{
		Addr:    ApplicationConfig.Host + ":" + ApplicationConfig.Port,
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Printf("%s Server Run http://127.0.0.1:%s/ \r\n", utils.GetCurrntTimeStr(), ApplicationConfig.Port)
	fmt.Printf("%s Swagger URL http://127.0.0.1:%s/swagger/index.html \r\n", utils.GetCurrntTimeStr(), ApplicationConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", utils.GetCurrntTimeStr())
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", utils.GetCurrntTimeStr())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	return nil
}
