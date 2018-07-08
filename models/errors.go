package models

import (
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")

	ErrNameNil      = errors.New("user name is nil")
	ErrUserIDNil    = errors.New("user ID is nil")
	ErrKgIDNil      = errors.New("kindergarten ID is nil")
	ErrBlockIDNil   = errors.New("block ID is nil")
	ErrDevTypeIDNil = errors.New("device type ID is nil")
	ErrTokenNil     = errors.New("token is nil")
	ErrCellUsed     = errors.New("cell number has been used")
)
