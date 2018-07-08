package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"kindergarten/models"

	log "github.com/astaxie/beego/logs"
	"github.com/lihuacat/surgemq/service"
	"github.com/surgemq/message"
)

var (
	client = &service.Client{}
)

func keepAlive() {
	//	ticker := time.NewTicker(10 * time.Second)
	for {
		time.Sleep(10 * time.Second)
		log.Debug("pinging...")
		err := client.Ping(onPing)
		if err != nil {
			log.Error(err)
			//				StartClient()
		}
	}

}

func StartClient() error {
	msg := message.NewConnectMessage()
	//	msg.SetWillQos(1)
	msg.SetVersion(4)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte("surgemq"))
	msg.SetKeepAlive(1000)
	//	msg.SetWillTopic([]byte("will"))
	//	msg.SetWillMessage([]byte("send me home"))
	msg.SetUsername([]byte("surgemq"))
	msg.SetPassword([]byte("verysecret"))
	err := client.Connect("tcp://47.94.4.51:444", msg)
	if err != nil {
		log.Error(err)
		return err
	}
	submsg := message.NewSubscribeMessage()
	//	submsg.Qos()
	err = submsg.AddTopic([]byte("remote/station/+/cmd"), 1)
	if err != nil {
		log.Error(err)
		return err
	}

	err = submsg.AddTopic([]byte("hello"), 0)
	if err != nil {
		log.Error(err)
		return err
	}

	err = client.Subscribe(submsg, nil, onPublish)
	if err != nil {
		log.Error(err)
		return err
	}

	pubmsg := message.NewPublishMessage()
	pubmsg.SetTopic([]byte("hello"))
	//	pubmsg.SetPayload(make([]byte, 1024))
	pubmsg.SetPayload([]byte("hahaha"))
	pubmsg.SetQoS(0)
	for {
		err = client.Publish(pubmsg, nil)
		if err != nil {
			log.Error(err)
			return err
		}
		time.Sleep(100 * time.Second)
		//		c.Ping(nil)
	}
	//	keepAlive()
	//	for {
	//		time.Sleep(10 * time.Second)
	//		log.Debug("pinging...")
	//		err := client.Ping(onPing)
	//		if err != nil {
	//			log.Error(err)
	//			//				StartClient()
	//		}
	//	}
	//	select {}

	return nil
}

func onPing(msg, ack message.Message, err error) error {
	log.Debug("Ping Server")

	return nil
}

func onPublish(msg *message.PublishMessage) error {
	var err error
	log.Debug(string(msg.Topic()))
	log.Debug(string(msg.Payload()))
	if len(msg.Payload()) == 0 {
		err = errors.New("payload is nil")
		log.Error(err)
		return nil
	}

	if strings.Contains(string(msg.Payload()), "ltstate") { //灯

		ltStatus := LigthStatus{}
		err = json.Unmarshal(msg.Payload(), &ltStatus)
		if err != nil {
			log.Error(err)
			return nil
		}
		dev, err := models.GetDevByRmtID(ltStatus.SwitchID)
		if err != nil {
			log.Error(err)
			return nil
		}
		//	fmt.Sprintf("rmid:%s/switchid:%s/sw1", ltStatus.RoomID, ltStatus.SwitchID)
		status, err := strconv.Atoi(ltStatus.LtStatus.Sw1)
		if err != nil {
			log.Debug(ltStatus.LtStatus.Sw2)
			log.Error(err)
			return nil
		} else {
			if status == 0 {
				status = -1
			}
			models.SetDevStatusbyDevIDPort(dev.DevID, 1, status)
		}
		status, err = strconv.Atoi(ltStatus.LtStatus.Sw2)
		if err != nil {
			log.Debug(ltStatus.LtStatus.Sw2)
			log.Error(err)
			return nil
		} else {
			if status == 0 {
				status = -1
			}
			models.SetDevStatusbyDevIDPort(dev.DevID, 2, status)
		}
		status, err = strconv.Atoi(ltStatus.LtStatus.Sw3)
		if err != nil {
			log.Debug(ltStatus.LtStatus.Sw2)
			log.Error(err)
			return nil
		} else {
			if status == 0 {
				status = -1
			}
			models.SetDevStatusbyDevIDPort(dev.DevID, 3, status)
		}

		block, err := models.GetBlockByRmtCtrlID(ltStatus.RoomID)
		if err != nil {
			log.Error("kaCmd.RmtCtrlID", ltStatus.RoomID)
			log.Error(err)
			return nil
		}
		err = models.DevUpdateLastTimeByBlkID(block.BlockID)
		if err != nil {
			log.Error(err)
			return nil
		}

	} else if strings.Contains(string(msg.Payload()), `"cmd":"keepalive"`) {
		kaCmd := KeepAliveCmd{}
		//		log.Debug("msg.Payload()", msg.Payload())
		err = json.Unmarshal(msg.Payload(), &kaCmd)
		if err != nil {
			log.Error(err)
			return nil
		}
		block, err := models.GetBlockByRmtCtrlID(kaCmd.RmtCtrlID)
		if err != nil {
			log.Error("kaCmd.RmtCtrlID", kaCmd.RmtCtrlID)
			log.Error(err)
			return nil
		}
		err = models.DevUpdateLastTimeByBlkID(block.BlockID)
		if err != nil {
			log.Error(err)
			return nil
		}
	} else { //空调

		acStatus := AirconditionStatus{}
		err = json.Unmarshal(msg.Payload(), &acStatus)
		if err != nil {
			log.Error(err)
			return nil
		}
		dev, err := models.GetDevByRmtID(acStatus.AmmeterID)
		if err != nil {
			log.Error(err)
			return nil
		}
		status, err := strconv.Atoi(acStatus.AcStatus)
		if err != nil {
			log.Error(err)
			return nil
		}
		if status == 0 {
			status = -1
		}
		err = models.SetDevStatusbyDevIDPort(dev.DevID, 1, status)
		if err != nil {
			log.Error(err)
			return nil
		}

		block, err := models.GetBlockByRmtCtrlID(acStatus.RmID)
		if err != nil {
			log.Error("kaCmd.RmtCtrlID", acStatus.RmID)
			log.Error(err)
			return nil
		}
		err = models.DevUpdateLastTimeByBlkID(block.BlockID)
		if err != nil {
			log.Error(err)
			return nil
		}
	}
	return nil
}

type SwitchList struct {
	Sw1 string `json:"sw1"`
	Sw2 string `json:"sw2"`
	Sw3 string `json:"sw3"`
}

type Switch1 struct {
	Sw1 string `json:"sw1"`
}
type Switch2 struct {
	Sw2 string `json:"sw2"`
}
type Switch3 struct {
	Sw3 string `json:"sw3"`
}

type LigthStatus struct {
	RoomID   string     `json:"rmid"`
	SwitchID string     `json:"switchid"`
	LtStatus SwitchList `json:"ltstate"`
}

type LigthCmd struct {
	RmtCtrlID string      `json:"rmid"`
	SwitchID  string      `json:"switchid"`
	LtCmd     interface{} `json:"ltcmd"`
}

type AirconditionStatus struct {
	RmID      string `json:"rmid"`
	AmmeterID string `json:"Ammeterid"`
	AcStatus  string `json:"acstate"`
}

type AirconditionProperty struct {
	Mode        string `json:"mode"`
	OnOff       string `json:"key"`
	WindSpeed   string `json:"windspeed"`
	WindScan    string `json:"windscan"`
	Temperature string `json:"temperature"`
}

type AirconditionCmd struct {
	RmtCtrlID string      `json:"rmid"`
	Cmd       interface{} `json:"rmcmd"`
}

type KeepAliveCmd struct {
	RmtCtrlID string `json:"rmid"`
	Cmd       string `json:"cmd"`
}

func ControlDevice(id int64, statuses []*models.DevStatus) error {
	var (
		topic   []byte
		payload []byte
	)
	//	if status == -1 {
	//		status = 0
	//	}
	dev, err := models.GetDevByID(id)
	if err != nil {
		log.Error(err)
		return err
	}

	if dev.DevTypeID == 1 || dev.DevTypeID == 3 || dev.DevTypeID == 4 || dev.DevTypeID == 5 { //设备类型是灯

		for _, status := range statuses {
			ltCmd := LigthCmd{}
			strs := strings.Split(dev.GetRmtDevID(), "_")
			ltCmd.RmtCtrlID = dev.GetRmtDevID()
			ltCmd.SwitchID = strs[0]
			if status.Port == 1 {
				ltCmd.LtCmd = &Switch1{
					Sw1: fmt.Sprint(status.Status),
				}
			} else if status.Port == 2 {
				ltCmd.LtCmd = &Switch2{
					Sw2: fmt.Sprint(status.Status),
				}
			} else if status.Port == 3 {
				ltCmd.LtCmd = &Switch3{
					Sw3: fmt.Sprint(status.Status),
				}
			}

			payload, err = json.Marshal(&ltCmd)
			if err != nil {
				log.Error(err)
				return err
			}
			topic = []byte(fmt.Sprintf("remote/command/%s/cmd", ltCmd.RmtCtrlID))
			pubmsg := message.NewPublishMessage()
			pubmsg.SetTopic(topic)
			//	pubmsg.SetPayload(make([]byte, 1024))
			pubmsg.SetPayload(payload)
			pubmsg.SetQoS(1)
			err = client.Publish(pubmsg, nil)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	} else if dev.DevTypeID == 2 {
		for _, status := range statuses {
			acPro := AirconditionProperty{
				Mode:        "0",
				OnOff:       fmt.Sprint(status.Status),
				WindSpeed:   "0",
				WindScan:    "0",
				Temperature: "25",
			}

			acCmd := AirconditionCmd{
				RmtCtrlID: dev.GetRmtDevID(),
				Cmd:       &acPro,
			}
			topic = []byte(fmt.Sprintf("remote/command/%s/cmd", acCmd.RmtCtrlID))
			payload, err = json.Marshal(&acCmd)
			pubmsg := message.NewPublishMessage()
			pubmsg.SetTopic(topic)
			//	pubmsg.SetPayload(make([]byte, 1024))
			pubmsg.SetPayload(payload)
			pubmsg.SetQoS(1)
			err = client.Publish(pubmsg, nil)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}
	//	log.Debug("topic:", topic)
	//	log.Debug("payload:", string(payload))

	return nil
}
