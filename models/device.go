package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	log "github.com/astaxie/beego/logs"
)

type DevAddReq struct {
	DevTypeID int64
	BlockID   int64
	DevName   string
	RmtDevID  string
}

func (this *DevAddReq) Check() error {

	if this.DevTypeID <= 0 {
		return ErrDevTypeIDNil
	}
	if this.BlockID <= 0 {
		return ErrBlockIDNil
	}
	if this.DevName == "" {
		return errors.New("DevName is nil")
	}
	if this.RmtDevID == "" {
		return errors.New("RmtDevID is nil")
	}

	return nil
}

type Device struct {
	DevTypeID int64
	BlockID   int64
	DevName   string
	DevID     int64
	BlockName string
	Status    []*DevStatus
	CtrlNum   int
	rmtCtrlID string
	rmtDevID  string
	lastTime  time.Time
}

type DevStatus struct {
	Port   int
	Status int //1 打开，-1 关闭， -2 离线
}

func (this *Device) GetRmtDevID() string {
	return this.rmtDevID
}

func (this *Device) SetRmtDevID(id string) {
	this.rmtDevID = id
}

func InsertDev(dev *Device) (int64, error) {

	if dev.DevTypeID == 1 { //如果是灯
		n, err := GetDevNumByBlockID(dev.BlockID, dev.DevTypeID)
		if err != nil {
			log.Error(err)
			return 0, err
		}
		dev.rmtDevID = fmt.Sprint(dev.rmtDevID, "_", "sw", n+1)
	}
	tx, err := db.Begin() //开始事务
	if err != nil {
		log.Error(err)
		return 0, err
	}

	var id int64
	row := tx.QueryRow(`select nextval('device_seq');`)
	err = row.Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	now := time.Now()
	//	dev.Status = -2
	_, err = tx.Exec(`insert into device(devid, devtypeid, blockid,devname, rmtdevid, lasttime ) values($1 ,$2, $3, $4, $5, $6);`, &id, &dev.DevTypeID, &dev.BlockID, &dev.DevName, &dev.rmtDevID, &now)
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return 0, err
	}

	dev.DevID = id
	var num int
	row = tx.QueryRow(`select ctrl_num from devicetype where typeid=$1`)
	err = row.Scan(&num)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	i := num
	for ; i > 0; i-- {
		_, err = tx.Exec(`insert into devstatus(devid, port, status) values($1 ,$2, -2);`, &id, &i)
		if err != nil {
			tx.Rollback()
			log.Error(err)
			return 0, err
		}
	}
	tx.Commit()
	return id, nil
}

func DelDev(id int64) error {
	tx, err := db.Begin() //开始事务
	ret, err := tx.Exec(`delete from devstatus  where devid = $1;`, &id)
	if err != nil {
		tx.Rollback()
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	n, err := ret.RowsAffected()
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	if n == 0 {
		return ErrNotFound
	}

	ret, err = tx.Exec(`delete from device  where devid = $1;`, &id)
	if err != nil {
		tx.Rollback()
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	n, err = ret.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	tx.Commit()
	return nil
}

func GetDevByID(id int64) (*Device, error) {
	var (
		dev    Device
		status int
	)
	log.Debug("id", id)
	row := db.QueryRow(`select d.devid,d.devtypeid, d.devname, d.blockid, b.blockname, d.status, d.lasttime, b.rmtctrlid, d.rmtdevid, dt.ctrl_num from device d left join devicetype dt on dt.typeid=d.devtypeid left join block b on b.blockid=d.blockid where d.devid=$1;`, &id)
	err := row.Scan(&dev.DevID, &dev.DevTypeID, &dev.DevName, &dev.BlockID, &dev.BlockName, &dev.Status, &dev.lastTime, &dev.rmtCtrlID, &dev.rmtDevID, &dev.CtrlNum)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if time.Now().Sub(dev.lastTime).Seconds() > 70 {
		status = -2
	}

	rows, err := db.Query(`select port,status from devstatus where devid = $1`, &dev.DevID)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	dev.Status = make([]*DevStatus, 0)
	for rows.Next() {
		sta := DevStatus{}
		err := rows.Scan(&sta.Port, &sta.Status)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if status == -2 {
			sta.Status = status
		}
		dev.Status = append(dev.Status, &sta)
	}

	return &dev, nil
}

type DevsRes struct {
	Num     int
	Devices []*Device
}

func GetDevsByKgID(kgid int64, status int) ([]*Device, error) {
	var (
		err  error
		rows *sql.Rows
		stat int
	)
	if status == 0 {
		rows, err = db.Query(`select device.devtypeid, device.devid, device.devname, device.blockid,block.blockname,device.status,device.lasttime,dt.ctrl_num from device left join devicetype dt on dt.typeid = device.devtypeid left join block on device.blockid=block.blockid where block.kgid=$1 order by device.devname;`, &kgid)
	} else {
		rows, err = db.Query(`select device.devtypeid, device.devid, device.devname, device.blockid,block.blockname,device.status,device.lasttime,dt.ctrl_num from device left join devicetype dt on dt.typeid = device.devtypeid left join block on device.blockid=block.blockid where block.kgid=$1 and device.status=$2 order by device.devname;`, &kgid, &status)
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	devs := make([]*Device, 0)
	for rows.Next() {
		dev := Device{}
		err = rows.Scan(&dev.DevTypeID, &dev.DevID, &dev.DevName, &dev.BlockID, &dev.BlockName, &dev.Status, &dev.lastTime, &dev.CtrlNum)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if time.Now().Sub(dev.lastTime).Seconds() > 70 {
			stat = -2
		}

		statusRows, err := db.Query(`select port,status from devstatus where devid = $1`, &dev.DevID)
		if err != nil {
			log.Error(err)
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			}
			return nil, err
		}
		dev.Status = make([]*DevStatus, 0)
		for statusRows.Next() {
			sta := DevStatus{}
			err := statusRows.Scan(&sta.Port, &sta.Status)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			if status == -2 {
				sta.Status = stat
			}
			dev.Status = append(dev.Status, &sta)
		}

		devs = append(devs, &dev)
	}

	return devs, nil
}

func GetDevsByBlkID(blkID int64, status int) ([]*Device, error) {
	var (
		err  error
		rows *sql.Rows
	)

	rows, err = db.Query(`select d.devid, d.devname, d.devtypeid, d.blockid,b.blockname,d.status,d.lasttime, dt.ctrl_num  from device d left join devicetype dt on dt.typeid = d.devtypeid left join block b on d.blockid=b.blockid where d.blockid =$1 ;`, &blkID)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	devs := make([]*Device, 0)
	for rows.Next() {
		dev := Device{}
		err = rows.Scan(&dev.DevID, &dev.DevName, &dev.DevTypeID, &dev.BlockID, &dev.BlockName, &dev.Status, &dev.lastTime, &dev.CtrlNum)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		//		if time.Now().Sub(dev.lastTime).Seconds() > 70 {
		//			dev.Status = -2
		//		}
		var statusRows *sql.Rows
		if status == 0 {
			statusRows, err = db.Query(`select port,status from devstatus where devid = $1`, &dev.DevID)
		} else {
			statusRows, err = db.Query(`select port,status from devstatus where devid = $1 and status=$2`, &dev.DevID, status)
		}
		if err != nil {
			log.Error(err)
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			}
			return nil, err
		}
		dev.Status = make([]*DevStatus, 0)
		for statusRows.Next() {
			sta := DevStatus{}
			err := statusRows.Scan(&sta.Port, &sta.Status)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			//			if status == -2 {
			//				sta.Status = status
			//			}
			dev.Status = append(dev.Status, &sta)
		}

		devs = append(devs, &dev)
	}

	return devs, nil
}

type DevStat struct {
	KgID             int64
	KgName           string
	LightOn          int
	AirConditionerOn int
}

type DevStatRes struct {
	Num      int
	DevStats []*DevStat
}

func KgDevOnNum(kgID int64, devTypeName string) (int, error) {
	var (
		id  int64
		num int
	)
	row := db.QueryRow(`select id from devclass where typename = $1;`, &devTypeName)
	err := row.Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	row = db.QueryRow(`select count(1) from device d left join devicetype dt on d.devtypeid = dt.typeid left join devclass dc on dc.id = dt.devclassid left join block b on d.blockid=b.blockid left join kindergarten k on k.kgid=b.kgid where dt.devclassid=$1 and k.kgid=$2 and d.status=1`, &id, &kgID)
	err = row.Scan(&num)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return num, nil
}

type Devices struct {
	Devs []*Device
}

func DevSetStatus(devID int64, status int) error {
	ret, err := db.Exec(`update device set status = $1  where devid = $2;`, &status, &devID)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	n, err := ret.RowsAffected()
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func SetDevStatusbyRmtDevID(rmtdevid string, status int) error {

	ret, err := db.Exec(`update device set status = $1  where rmtdevid = $2;`, &status, &rmtdevid)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	n, err := ret.RowsAffected()
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func SetDevStatusbyDevIDPort(devid int64, port, status int) error {

	ret, err := db.Exec(`update devstatus set status = $1  where devid = $2;`, &status, &devid, &port)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	n, err := ret.RowsAffected()
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func GetDevPathbyID(id int64) (string, error) {
	var path string
	log.Debug("id", id)
	row := db.QueryRow(`select rmtpath from device where devid=$1;`, &id)
	err := row.Scan(&path)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return "", ErrNotFound
		}
		return "", err
	}

	return path, nil
}

func GetDevNumByBlockID(blockid int64, devtype int64) (int, error) {
	var ret int
	row := db.QueryRow(`select count(1) from device  where blockid = $1 and devtypeid=$2;`, &blockid, &devtype)
	err := row.Scan(&ret)
	if err != nil {
		log.Error(err)

		return 0, err
	}

	return ret, nil
}

func GetDevByBlockIDDevType(blockid int64, devtype int64) (*Device, error) {
	var dev Device
	row := db.QueryRow(`select d.devid,d.devtypeid, d.devname, d.blockid,b.blockname, d.status, b.rmtctrlid, d.rmtdevid,d.lasttime,dt.ctrl_num from device d left join devicetype dt on dt.typeid=d.devtypeid left join block b on b.blockid=d.blockid where b.blockid=$1 and d.devtypeid=$2;`, &blockid, &devtype)
	err := row.Scan(&dev.DevID, &dev.DevTypeID, &dev.DevName, &dev.BlockID, &dev.BlockName, &dev.Status, &dev.rmtCtrlID, &dev.rmtDevID, &dev.lastTime, &dev.CtrlNum)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var status int
	if time.Now().Sub(dev.lastTime).Seconds() > 70 {
		status = -2
	}

	statusRows, err := db.Query(`select port,status from devstatus where devid = $1`, &dev.DevID)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	dev.Status = make([]*DevStatus, 0)
	for statusRows.Next() {
		sta := DevStatus{}
		err := statusRows.Scan(&sta.Port, &sta.Status)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if status == -2 {
			sta.Status = status
		}
		dev.Status = append(dev.Status, &sta)
	}

	return &dev, nil
}

func DevUpdateLastTimeByBlkID(blockid int64) error {

	now := time.Now()
	ret, err := db.Exec(`update device set lasttime = $1  where blockid = $2;`, &now, &blockid)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	n, err := ret.RowsAffected()
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func GetDevByRmtID(rmtID string) (*Device, error) {
	var dev Device
	row := db.QueryRow(`select devid,devtypeid,devname,blockid from device  where rmtdevid = $1;`, &rmtID)
	err := row.Scan(&dev.DevID, &dev.DevTypeID, &dev.DevName, &dev.BlockID)
	if err != nil {
		log.Error(err)

		return nil, err
	}

	return &dev, nil
}
