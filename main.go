package main

import (
	"kindergarten/mqtt"
	_ "kindergarten/routers"
	"time"

	"github.com/astaxie/beego"
	log "github.com/astaxie/beego/logs"
)

func main() {
	log.SetLogFuncCallDepth(3)
	log.SetLevel(log.LevelDebug)
	log.EnableFuncCallDepth(true)
	log.SetLevel(log.LevelDebug)

	//	log.SetLevel(3)
	go mqtt.StartMqttServer() //mqtt服务
	time.Sleep(2 * time.Second)
	log.Info("Starting Client")
	go mqtt.StartClient()
	log.Info("Start Client complete")
	beego.SetStaticPath("/swagger", "swagger")
	beego.Run()
}
