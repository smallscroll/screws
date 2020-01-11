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

//Configuration 配置
type Configuration struct {
	Parameters map[string]interface{}
	ConfigFile string //json
	Logfile    string
}

//LoadConfig 加载配置
func (c *Configuration) LoadConfig() {
	data, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(data, &c.Parameters); err != nil {
		log.Fatal(err)
	}
}

//Applog  应用日志
func (c *Configuration) Applog() {
	logFile, err := os.OpenFile(c.Logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
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
