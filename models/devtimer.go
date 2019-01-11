package models

import (
	"sync"
	"time"

	log "github.com/astaxie/beego/logs"
)

type (
	DevCtrlFunc func(devid int64, status []*DevStatus) error
)

type DevTimer struct {
	TimerID    int64 //ID
	DevID      int64 //设备
	TurnHour   int   //开关时间点,小时
	TurnMinute int   //开关时间点,分钟
	Port       int   //线路
	OnOff      int   //开关1 开,-1 关
}
type AddDevTimerReq struct {
	DevID      int64 //设备
	TurnHour   int   //开关时间点,小时
	TurnMinute int   //开关时间点,分钟
	Port       int   //线路
	OnOff      int   //开关1 开,-1 关
}

type DevTimers struct {
	timers map[int64]*DevTimer
	lock   sync.RWMutex
}

func (this *DevTimers) Add(devTimer *DevTimer) {

	defer this.lock.Unlock()
	this.lock.Lock()
	this.timers[devTimer.TimerID] = devTimer

}

func (this *DevTimers) Delete(id int64) {
	defer this.lock.Unlock()
	this.lock.Lock()
	delete(this.timers, id)
}

func (this *DevTimers) Run(ctrlfunc DevCtrlFunc) {

	if ctrlfunc == nil {
		return
	}
	ticker := time.NewTicker(60 * time.Second)
	devStatus := make([]*DevStatus, 1)
	var now time.Time

	for {
		now = <-ticker.C
		this.lock.RLock()
		for _, devTimer := range this.timers {
			if now.Hour() == devTimer.TurnHour && now.Minute() == devTimer.TurnMinute {
				devStatus[0] = &DevStatus{}
				devStatus[0].Port = devTimer.Port
				devStatus[0].Status = devTimer.OnOff
				err := ctrlfunc(devTimer.DevID, devStatus)
				if err != nil {
					log.Error("ctrlfunc:", err)
				}
			}
		}
		this.lock.RUnlock()
	}
}

func NewDevTimer() *DevTimers {
	timers := &DevTimers{}
	timers.timers = make(map[int64]*DevTimer)
	return timers
}

var (
	devTimers *DevTimers
)

func initDevTimer() {
	devTimers = NewDevTimer()
	rows, err := db.Query(`select timerid,devid,on_off,port,turn_hour,turn_minute from devicetimer`)
	if err != nil {
		log.Error(err)
		return
	}

	for rows.Next() {
		devTimer := &DevTimer{}
		err = rows.Scan(&devTimer.TimerID, &devTimer.DevID, &devTimer.OnOff, &devTimer.Port, &devTimer.TurnHour, &devTimer.TurnMinute)
		if err != nil {
			log.Error(err)
			return
		}
		devTimers.Add(devTimer)
	}

}

func AddDevTimer(devTimer *DevTimer) (int64, error) {
	var id int64
	row := db.QueryRow(`select nextval('userid_seq');`)
	err := row.Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	_, err = db.Exec(`insert into devicetimer( timerid,devid,on_off,turn_hour,turn_minute,port ) values($1,$2,$3,$4,$5,$6);`, &id, &devTimer.DevID, &devTimer.OnOff, &devTimer.TurnHour, &devTimer.TurnMinute, &devTimer.Port)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	devTimer.TimerID = id
	devTimers.Add(devTimer)

	return id, nil
}

func DelDevTimer(id int64) error {
	_, err := db.Exec(`delete from devicetimer where timerid=$1;`, id)
	if err != nil {
		log.Error(err)
		return err
	}
	devTimers.Delete(id)

	return nil
}

func RunDevTimer(ctrlfunc DevCtrlFunc) {
	go devTimers.Run(ctrlfunc)
}
