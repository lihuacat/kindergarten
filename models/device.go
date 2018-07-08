package models

import (
	"database/sql"
	"errors"

	log "github.com/astaxie/beego/logs"
)

type DevAddReq struct {
	DevTypeID int64
	BlockID   int64
	DevName   string
}

func (this *DevAddReq) Check() error {

	if this.DevTypeID <= 0 {
		return ErrDevTypeIDNil
	}
	if this.BlockID <= 0 {
		return ErrBlockIDNil
	}
	if this.DevName == "" {
		return errors.New("block name is nil")
	}
	return nil
}

type Device struct {
	DevTypeID int64
	BlockID   int64
	DevName   string
	DevID     int64
	BlockName string
	Status    int
}

func InsertDev(dev *Device) (int64, error) {
	var id int64
	row := db.QueryRow(`select nextval('device_seq');`)
	err := row.Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	_, err = db.Exec(`insert into device(devid, devtypeid, blockid,devname ) values($1,$2,$3,$4);`, &id, &dev.DevTypeID, &dev.BlockID, &dev.DevName)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return id, nil
}

func DelDev(id int64) error {
	ret, err := db.Exec(`delete from device  where devid = $1;`, &id)
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

func GetDevByID(id int64) (*Device, error) {
	var dev Device
	log.Debug("id", id)
	row := db.QueryRow(`select d.devid, d.devname, d.blockid,b.blockname, d.status from device d left join block b on b.blockid=d.blockid where d.devid=$1;`, &id)
	err := row.Scan(&dev.DevID, &dev.DevName, &dev.BlockID, &dev.BlockName, &dev.Status)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
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
	)
	if status == 0 {
		rows, err = db.Query(`select device.devtypeid, device.devid, device.devname, device.blockid,block.blockname,device.status from device left join block on device.blockid=block.blockid where block.kgid=$1 order by device.devname;`, &kgid)
	} else {
		rows, err = db.Query(`select device.devtypeid, device.devid, device.devname, device.blockid,block.blockname,device.status from device left join block on device.blockid=block.blockid where block.kgid=$1 and device.status=$2 order by device.devname;`, &kgid, &status)
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	devs := make([]*Device, 0)
	for rows.Next() {
		dev := Device{}
		err = rows.Scan(&dev.DevTypeID, &dev.DevID, &dev.DevName, &dev.BlockID, &dev.BlockName, &dev.Status)
		if err != nil {
			log.Error(err)
			return nil, err
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
	if status == 0 {
		rows, err = db.Query(`select d.devid, d.devname, d.devtypeid, d.blockid,b.blockname,d.status from device d left join block b on d.blockid=b.blockid where d.blockid =$1 ;`, &blkID)
	} else {
		rows, err = db.Query(`select d.devid, d.devname, d.devtypeid, d.blockid,b.blockname,d.status from device d left join block b on d.blockid=b.blockid where d.blockid =$1 and d.status=$2;`, &blkID, &status)
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	devs := make([]*Device, 0)
	for rows.Next() {
		dev := Device{}
		err = rows.Scan(&dev.DevID, &dev.DevName, &dev.DevTypeID, &dev.BlockID, &dev.BlockName, &dev.Status)
		if err != nil {
			log.Error(err)
			return nil, err
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
		lightID int64
		num     int
	)
	row := db.QueryRow(`select typeid from devicetype where typename = $1;`, &devTypeName)
	err := row.Scan(&lightID)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	row = db.QueryRow(`select count(1) from device d left join block b on d.blockid=b.blockid left join kindergarten k on k.kgid=b.kgid where d.devtypeid=$1 and k.kgid=$2 and d.status=1`, &lightID, &kgID)
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

func SetDevStatusbyRmtPath(path string, status int) error {
	ret, err := db.Exec(`update device set status = $1  where rmtpath = $2;`, &status, &path)
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
