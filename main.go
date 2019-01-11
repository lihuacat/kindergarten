package main

import (
	"kindergarten/mqtt"
	_ "kindergarten/routers"
	"time"

	"github.com/astaxie/beego"
	log "github.com/astaxie/beego/logs"
)

func main() {

	loglevel, err := beego.AppConfig.Int("log_level")
	if err != nil {
		log.Error(err)
		loglevel = 4
	}

	log.SetLogFuncCallDepth(3)
	log.SetLevel(log.LevelDebug)
	log.EnableFuncCallDepth(true)
	log.SetLevel(loglevel)

	//	log.SetLevel(3)
	go mqtt.StartMqttServer() //mqtt服务
	time.Sleep(2 * time.Second)
	log.Info("Starting Client")
	go mqtt.StartClient()
	log.Info("Start Client complete")
	beego.SetStaticPath("/swagger", "swagger")
	beego.SetStaticPath("/download", "download")
	beego.Run()
}
