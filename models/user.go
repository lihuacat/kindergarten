package models

import (
	"database/sql"
	"errors"

	log "github.com/astaxie/beego/logs"
)

type ModifyPasswdReq struct {
	//	UserToken
	OldPasswd string
	NewPasswd string
}

func (this *ModifyPasswdReq) Check() error {
	var err error
	if this.OldPasswd == "" {
		err = errors.New("old password is nil")
	}
	if this.NewPasswd == "" {
		err = errors.New("new password is nil")
	}

	return err
}

type ModifyUserReq struct {
	//	UserToken
	NewCellNum  string
	NewUserName string
}

func (this *ModifyUserReq) Check() error {

	if this.NewCellNum == "" && this.NewUserName == "" {
		return errors.New("New cell NO. and new name are nil")
	}
	return nil
}

type ModifyUserRes struct {
	UserID   int64
	CellNum  string
	UserName string
}

type RegRes struct {
	LoginRes
	CellNum string
}

type RegReq struct {
	LoginReq
	UserName string
}

func (this *RegReq) Check() error {
	err := this.LoginReq.Check()
	if err != nil {
		return err
	}
	if this.UserName == "" {
		return ErrNameNil
	}

	return nil
}

type LoginReq struct {
	CellNum  string //手机号码
	Password string //登录密码

}

func (this *LoginReq) Check() error {
	if this.CellNum == "" {
		return errors.New("cell number is needed")
	}
	if this.Password == "" {
		return errors.New("password is needed")
	}

	return nil
}

type LoginRes struct {
	UserID   int64  //用户ID
	UserName string //用户名
	Token    string //返回的Token
}

type UserToken struct {
	UserID int64  //用户ID
	Token  string //Token
}

func (this *UserToken) CheckToken() error {
	if this.UserID <= 0 {
		return ErrUserIDNil
	}
	if this.Token == "" {
		return ErrTokenNil
	}

	return CheckSession(this.UserID, this.Token)
	//	return nil
}

type User struct {
	UserID   int64
	Passwd   string
	UserName string
	CellNum  string
}

func UpdateUser(user *User) error {
	_, err := db.Exec(`update users set username = $1, cellnum = $2 where userid = $3;`, &user.UserName, &user.CellNum, &user.UserID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func UpdateUserPasswd(user *User) error {
	_, err := db.Exec(`update users set passwd = $1 where userid = $2;`, &user.Passwd, &user.UserID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func GetUserByID(id int64) (*User, error) {
	user := User{}
	row := db.QueryRow("select userid,username,passwd,cellnum from users where userid = $1 ", id)
	err := row.Scan(&user.UserID, &user.UserName, &user.Passwd, &user.CellNum)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByCellNum(cellNum string) (*User, error) {
	user := User{}
	row := db.QueryRow("select userid,username,passwd,cellnum from users where cellnum = $1 ", cellNum)
	err := row.Scan(&user.UserID, &user.UserName, &user.Passwd, &user.CellNum)
	if err != nil {
		log.Error(err)
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func InsertUser(user *User) (int64, error) {
	var id int64
	row := db.QueryRow(`select nextval('userid_seq');`)
	err := row.Scan(&id)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	_, err = db.Exec(`insert into users( username,"passwd",userid,cellnum ) values($1,$2,$3,$4);`, &user.UserName, &user.Passwd, &id, &user.CellNum)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return id, nil
}
