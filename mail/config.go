package mail

import (
	"flag"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/emacsist/go-common/helper/file"
)

var mailFilePath *string

// MailConfig : mail的配置文件内容
var MailConfig mailConfig

// mailConfig : 配置文件
type mailConfig struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
}

// LoadConfig : 加载
func loadConfig() {
	e := file.LoadJSON(*mailFilePath, &MailConfig)
	if e != nil {
		panic("init mail config error : " + e.Error())
	}
	log.Infof("init %v OK. %v", *mailFilePath, MailConfig)
}

func parseCommandLine() {

	pwd, e := os.Getwd()
	if e != nil {
		panic(e.Error())
	}
	log.Infof("working directory: %v", pwd)

	mailFilePath = flag.String("mail.config", pwd+"/mail.json", "-mail.config=/tmp/mail.json")
	log.Infof("mail config file is %v", *mailFilePath)
}

func init() {
	parseCommandLine()
	loadConfig()
}
