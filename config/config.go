package config

import (
	"flag"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/emacsist/go-common/helper/file"
)

var configPath *string

// Configuration : 配置文件内容
var Configuration configuration

// mailConfig : 配置文件
type configuration struct {
	Callback callback `json:"callback"`
	Email    email    `json:"email"`
}

type email struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
}

// callback : 配置文件
type callback struct {
	AMQPURL   string `json:"url"`
	IsOK      bool   `json:"-"`
	QueueName string `json:"mail.queue"`
	Worker    int    `json:"worker"`
}

// LoadConfig : 加载
func loadConfig() {
	e := file.LoadJSON(*configPath, &Configuration)
	if e != nil {
		panic("init config error : " + e.Error())
	}
	log.Infof("init %v OK. %v", *configPath, Configuration)
	if len(Configuration.Callback.AMQPURL) > 0 {
		Configuration.Callback.IsOK = true
		log.Infof("use rabbitmq as callback queue")
	}
}

func parseCommandLine() {

	pwd, e := os.Getwd()
	if e != nil {
		panic(e.Error())
	}
	log.Infof("working directory: %v", pwd)

	configPath = flag.String("config", pwd+"/config.json", "-config=/tmp/config.json")
	log.Infof("config file is %v", *configPath)
}

func init() {
	parseCommandLine()
	loadConfig()
}
