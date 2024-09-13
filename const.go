package gormxmysql

import "errors"

const (
	DuplicateEntryErrCode = 1062
)

var ErrDuplicateKey = errors.New("duplicate entry for key error")
