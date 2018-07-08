package models

import (
	"strings"

	log "github.com/astaxie/beego/logs"
)

type RmtSuite struct {
	RmtCtrlID     string
	AmmeterID     string
	LightSwitchID string
}

func GetRmtSuiteByRmtCtrlID(id string) (*RmtSuite, error) {
	ret := &RmtSuite{}
	block, err := GetBlockByRmtCtrlID(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	light, err := GetDevByBlockIDDevType(block.BlockID, 1)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	ac, err := GetDevByBlockIDDevType(block.BlockID, 2)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	ret.RmtCtrlID = id
	ret.LightSwitchID = strings.Split(light.GetRmtDevID(), "_")[0]
	ret.AmmeterID = ac.GetRmtDevID()

	return ret, nil

}
