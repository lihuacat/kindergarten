package mqtt

import (
	//	"fmt"

	log "github.com/astaxie/beego/logs"
	"github.com/lihuacat/surgemq/service"
)

func StartMqttServer() error {
	//	glog.CopyStandardLogTo("INFO")
	server := service.Server{
		KeepAlive:      1000,
		ConnectTimeout: 1000,
		TimeoutRetries: 2,
	}

	err := server.ListenAndServe("tcp://:1883")
	if err != nil {
		log.Error(err)

	}
	return err
}
