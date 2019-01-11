package models

import (
	//	"database/sql"
	//	"errors"

	log "github.com/astaxie/beego/logs"
)

type BlockUser struct {
	UserID  int64
	BlockID int64
}

func (this *BlockUser) Check() error {
	if this.BlockID <= 0 {
		return ErrBlockIDNil
	}
	if this.UserID <= 0 {
		return ErrUserIDNil
	}
	return nil
}

func InserBlockUser(bu *BlockUser) error {
	_, err := db.Exec(`insert into block_user(blockid, user_id) values($1,$2);`, &bu.BlockID, &bu.UserID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
