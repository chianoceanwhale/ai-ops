package config

import (
	cm "ai-ops/internal/pkg/common"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type LogDir struct {
	Dir string
}

func InitLog(cfg *viper.Viper) *LogDir {
	return &LogDir{
		Dir: cfg.GetString("dir"),
	}
}

var LogConfig = new(LogDir)

func InitLogger() {
	log.SetFormatter(&log.JSONFormatter{FieldMap: log.FieldMap{
		log.FieldKeyTime:  "@timestamp",
		log.FieldKeyLevel: "@level",
		log.FieldKeyMsg:   "@message"}})

	switch cm.Mode(cm.ApplicationConfig.Mode) {
	case cm.ModeDev, cm.ModeTest:
		log.SetOutput(os.Stdout)
		log.SetLevel(log.TraceLevel)
	case cm.ModeProd, cm.ModeDebug:
		log.AddHook(newLfsHook(LogConfig.Dir+"ai-ops.log", 30))
	}

	log.SetReportCaller(true)
}

//日志按天滚动
func newLfsHook(logName string, maxRemainCnt uint) log.Hook {
	writer, err := rotatelogs.New(
		logName+".%Y%m%d",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间，这里设置为24小时分割一次
		rotatelogs.WithRotationTime(24*time.Hour),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(maxRemainCnt),
	)

	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}

	log.SetLevel(log.TraceLevel)

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{DisableColors: true})

	return lfsHook
}
