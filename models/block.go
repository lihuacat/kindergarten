package models

import (
	"database/sql"
	"errors"

	log "github.com/astaxie/beego/logs"
)

type BlockDelReq struct {
	BlockID int64
}

func (this *BlockDelReq) Check() error {

	if this.BlockID <= 0 {
		return ErrBlockIDNil
	}

	return nil
}

type BolckAddReq struct {
	KgID      int64
	BlockName string
	RmtCtrlID string
}

func (this *BolckAddReq) Check() error {

	if this.KgID <= 0 {
		return ErrKgIDNil
	}
	if this.BlockName == "" {
		return errors.New("BlockName is nil")
	}
	if this.RmtCtrlID == "" {
		return errors.New("RmtCtrlID name is nil")
	}
	return nil
}

type Block struct {
	BlockID   int64
	BlockName string
	KgID      int64
	RmtCtrlID string
}

func InsertBlock(block *Block) (int64, error) {
	var id int64
	row := db.QueryRow(`select nextval('blockid_seq');`)
	err := row.Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	_, err = db.Exec(`insert into block( blockid,blockname,kgid, rmtctrlid ) values($1,$2,$3,$4);`, &id, &block.BlockName, &block.KgID, &block.RmtCtrlID)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return id, nil
}

func DelBlock(id int64) error {
	ret, err := db.Exec(`delete from block  where blockid = $1;`, &id)
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

func GetBlockByID(id int64) (*Block, error) {
	var block Block
	row := db.QueryRow(`select blockid,blockname,kgid,rmtctrlid from block where blockid=$1;`, &id)
	err := row.Scan(&block.BlockID, &block.BlockName, &block.KgID, &block.RmtCtrlID)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &block, nil
}

type BlocksRes struct {
	Num    int
	Blocks []*Block
}

func GetBlocksByKgID(kgID int64) ([]*Block, error) {
	rows, err := db.Query(`select blockid,blockname,kgid from block where kgid =$1 ;`, &kgID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	blocks := make([]*Block, 0)
	for rows.Next() {
		block := Block{}
		err = rows.Scan(&block.BlockID, &block.BlockName, &block.KgID)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		blocks = append(blocks, &block)
	}

	return blocks, nil
}

func GetBlockByRmtCtrlID(rmtctrlid string) (*Block, error) {
	var block Block
	row := db.QueryRow(`select blockid,blockname,kgid,rmtctrlid from block where rmtctrlid=$1;`, &rmtctrlid)
	err := row.Scan(&block.BlockID, &block.BlockName, &block.KgID, &block.RmtCtrlID)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &block, nil
}
