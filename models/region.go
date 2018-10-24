package models

import (
	//	"database/sql"
	//	"errors"

	log "github.com/astaxie/beego/logs"
)

type Region struct {
	ID int64
	RegionReq
}

type RegionsRes struct {
	Num     int
	Regions []*Region
}

type RegionReq struct {
	Name string
	KgID int64
}

func (this *RegionReq) Check() error {

	if this.KgID <= 0 {
		return ErrBlockIDNil
	}
	if this.Name == "" {
		return ErrNameNil
	}

	return nil
}

func AddRegion(region *Region) (int64, error) {

	var id int64
	row := db.QueryRow(`select nextval('region_seq');`)
	err := row.Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	_, err = db.Exec(`insert into region( id,name,kg_id ) values($1,$2,$3);`, &id, &region.Name, &region.KgID)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return id, nil

}

func AddRegionBlock(regionID int64, blockID int64) error {
	_, err := db.Exec(`update block set region_id=$1 where blockid=$2;`, regionID, blockID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func DelRegion(regionID int64) error {
	_, err := db.Exec(`delete from region where region.id=$1;`, regionID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func GetRegionbyKgID(kgID int64) ([]*Region, error) {
	rows, err := db.Query(`select id,name,kg_id from region where kg_id =$1 order by name;`, &kgID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	regions := make([]*Region, 0)
	for rows.Next() {
		region := Region{}
		err = rows.Scan(&region.ID, &region.Name, &region.KgID)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		regions = append(regions, &region)
	}

	return regions, nil
}
