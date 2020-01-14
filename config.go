package screws

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//IConfiguration 配置接口
type IConfiguration interface {
	LoadConfig()
	GetConfig(key string) interface{}
	Applog()
	LoadAlarmWithEmail(host, port, username, password string)
	SendAlarmWithEmail(from, to, subject, content string)
}

//NewConfiguration 初始化配置(配置文件，日志文件)
func NewConfiguration(configFile, LogFile string) IConfiguration {
	return &configuration{
		ConfigFile: configFile,
		LogFile:    LogFile,
		AlarmChan:  make(chan [4]string),
	}
}

//configuration 配置
type configuration struct {
	Parameters map[string]interface{}
	ConfigFile string //json
	LogFile    string
	AlarmChan  chan [4]string
}

//LoadConfig 加载配置
func (c *configuration) LoadConfig() {
	data, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(data, &c.Parameters); err != nil {
		log.Fatal(err)
	}
}

//GetConfig 获取配置
func (c *configuration) GetConfig(key string) interface{} {
	if v, ok := c.Parameters[key]; ok {
		return v
	}
	return nil
}

//Applog  应用日志
func (c *configuration) Applog() {
	logFile, err := os.OpenFile(c.LogFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		sig := <-c
		log.Fatalln("System signal:", sig)
	}()
}

//LoadAlarmWithEmail  安装邮件警报
func (c *configuration) LoadAlarmWithEmail(host, port, username, password string) {
	mailSender := NewMailSender(host, port, username, password)
	go func() {
		for {
			a := <-c.AlarmChan
			go func() {
				if err := mailSender.SendWithTLS(a[0], a[1], a[2], a[3]); err != nil {
					log.Println(err)
				}
			}()
		}
	}()
}

//SendAlarmWithEmail 发送邮件警报
func (c *configuration) SendAlarmWithEmail(from, to, subject, content string) {
	c.AlarmChan <- [4]string{from, to, subject, content}
}
