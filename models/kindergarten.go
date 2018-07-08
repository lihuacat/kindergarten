package models

import (
	"database/sql"
	"errors"

	log "github.com/astaxie/beego/logs"
)

//type KgRes struct {
//	KgID   int64  //园区ID
//	KgName string //园区名称
//}

type KgAddReq struct {
	KgName string //园区名称
}

func (this *KgAddReq) Check() error {
	var err error
	if this.KgName == "" {
		err = errors.New("kindergarten name is nil")
	}
	return err
}

type KgModReq struct {
	//	UserToken
	KgID   int64
	KgName string //园区名称
}

func (this *KgModReq) Check() error {
	var err error
	if this.KgName == "" {
		err = errors.New("kindergarten name is nil")
	}
	if this.KgID <= 0 {
		err = errors.New("kindergarten ID is nil")
	}
	return err
}

type KgDelReq struct {
	KgID int64
}

func (this *KgDelReq) Check() error {
	var err error

	if this.KgID <= 0 {
		err = errors.New("kindergarten ID is nil")
	}
	return err
}

type Kindergarten struct {
	KgID   int64
	KgName string
}

type KgsRes struct {
	Num int64
	Kgs []*Kindergarten
}

func InsertKindergarten(kg *Kindergarten) (int64, error) {
	var id int64
	row := db.QueryRow(`select nextval('kgid_seq');`)
	err := row.Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	_, err = db.Exec(`insert into kindergarten( kgid,kgname ) values($1,$2);`, &id, &kg.KgName)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return id, nil
}

func GetKgByID(id int64) (*Kindergarten, error) {
	var kg Kindergarten
	row := db.QueryRow(`select kgid,kgname from kindergarten where kgid=$1;`, &id)
	err := row.Scan(&kg.KgID, &kg.KgName)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &kg, nil
}

func GetKgByName(name string) (*Kindergarten, error) {
	var kg Kindergarten
	row := db.QueryRow(`select kgid,kgname from kindergarten where kgname=$1;`, &name)
	err := row.Scan(&kg.KgID, &kg.KgName)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &kg, nil
}

func UpdateKg(kg *Kindergarten) error {
	_, err := db.Exec(`update kindergarten set kgname = $1 where kgid = $2;`, &kg.KgName, &kg.KgID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func DelKg(id int64) error {
	ret, err := db.Exec(`delete from kindergarten  where kgid = $1;`, &id)
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

func GetKgsByUserID(userID int64) ([]*Kindergarten, error) {

	rows, err := db.Query(`select kg.kgid,kg.kgname from block_user bu left join block b on bu.blockid=b.blockid left join kindergarten kg on kg.kgid=b.kgid where bu.userid=$1 group by kg.kgid,kg.kgname;`, &userID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	kgs := make([]*Kindergarten, 0)
	for rows.Next() {
		kg := Kindergarten{}
		err = rows.Scan(&kg.KgID, &kg.KgName)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		kgs = append(kgs, &kg)
	}

	return kgs, nil

}

func GetAllKgs() ([]*Kindergarten, error) {
	rows, err := db.Query(`select kgid,kgname from kindergarten;`)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	kgs := make([]*Kindergarten, 0)
	for rows.Next() {
		kg := Kindergarten{}
		err = rows.Scan(&kg.KgID, &kg.KgName)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		kgs = append(kgs, &kg)
	}

	return kgs, nil
}
