package models

import (
	"database/sql"
	"errors"

	log "github.com/astaxie/beego/logs"
)

type DevTypeDelReq struct {
	DevTypeID int64
}

func (this *DevTypeDelReq) Check() error {

	if this.DevTypeID <= 0 {
		return ErrBlockIDNil
	}

	return nil
}

type DevTypeAddReq struct {
	DevTypeName string
}

func (this *DevTypeAddReq) Check() error {

	if this.DevTypeName == "" {
		return errors.New("block name is nil")
	}
	return nil
}

type DevType struct {
	DevTypeID   int64
	DevTypeName string
	PortNum     int //端口数量，一个端口控制一个电器（灯/空调）
}

func InsertDevType(devType *DevType) (int64, error) {
	var id int64
	row := db.QueryRow(`select nextval('devtypeid_seq');`)
	err := row.Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	_, err = db.Exec(`insert into devicetype( typeid,typename ) values($1,$2);`, &id, &devType.DevTypeName)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return id, nil
}

func DelDevType(id int64) error {
	ret, err := db.Exec(`delete from devicetype  where typeid = $1;`, &id)
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

func GetDevTypeByID(id int64) (*DevType, error) {
	var devType DevType
	row := db.QueryRow(`select typeid,typename from devicetype where typeid=$1;`, &id)
	err := row.Scan(&devType.DevTypeID, &devType.DevTypeName)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &devType, nil
}

func GetDevTypeByName(name string) (*DevType, error) {
	var devType DevType
	row := db.QueryRow(`select typeid,typename from devicetype where typename=$1;`, &name)
	err := row.Scan(&devType.DevTypeID, &devType.DevTypeName)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &devType, nil
}

type DevTypesRes struct {
	Num      int
	DevTypes []*DevType
}

func GetDevTypes() ([]*DevType, error) {
	rows, err := db.Query(`select typeid,typename from devicetype ;`)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	devTypes := make([]*DevType, 0)
	for rows.Next() {
		devType := DevType{}
		err = rows.Scan(&devType.DevTypeID, &devType.DevTypeName)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		devTypes = append(devTypes, &devType)
	}

	return devTypes, nil
}
