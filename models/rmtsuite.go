package models

import (
	"strings"

	log "github.com/astaxie/beego/logs"
)

type RmtSuite struct {
	RmtCtrlID     string
	AmmeterID     []string
	LightSwitchID []string
}

func GetRmtSuiteByRmtCtrlID(id string) (*RmtSuite, error) {
	ret := &RmtSuite{}
	lts := make(map[string]bool)
	ammeters := make(map[string]bool)
	ltList := make([]string, 0)
	ammeterList := make([]string, 0)

	block, err := GetBlockByRmtCtrlID(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	devs, err := GetDevsByBlkID(block.BlockID, 0)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for _, dev := range devs {
		log.Debug(dev.GetRmtDevID())
		if dev.DevTypeID == 1 || dev.DevTypeID == 3 || dev.DevTypeID == 4 || dev.DevTypeID == 5 {
			lts[strings.Split(dev.GetRmtDevID(), "_")[0]] = true
		} else if dev.DevTypeID == 2 || dev.DevTypeID == 6 {
			ammeters[dev.GetRmtDevID()] = true
		}
	}

	for lt, _ := range lts {
		ltList = append(ltList, lt)
	}
	for ac, _ := range ammeters {
		ammeterList = append(ammeterList, ac)
	}

	//	light, err := GetDevByBlockIDDevType(block.BlockID, 1)
	//	if err != nil {
	//		log.Error(err)
	//		return nil, err
	//	}

	//	ac, err := GetDevByBlockIDDevType(block.BlockID, 2)
	//	if err != nil {
	//		log.Error(err)
	//		return nil, err
	//	}

	//	ret.RmtCtrlID = id
	//	//	ret.LightSwitchID = strings.Split(light.GetRmtDevID(), "_")[0]
	//	ret.AmmeterID = ac.GetRmtDevID()
	ret.AmmeterID = ammeterList
	ret.LightSwitchID = ltList
	ret.RmtCtrlID = id
	return ret, nil

}
