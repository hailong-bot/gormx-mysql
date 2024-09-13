package gormxmysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/hailong-bot/gormx"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DOBase struct {
	DataObjecter gormx.DataObjecter `gorm:"-" json:"-"`
	ID           int64              `gorm:"type:int(11);primaryKey;autoIncrement" json:"id"`
}

func (d *DOBase) GetIDer() interface{} {
	return d.ID
}

func (d *DOBase) Updates(db *gorm.DB, values gormx.UPO) error {
	if err := db.Model(d.DataObjecter).Updates(map[string]interface{}(values)).Error; err != nil {
		if mySQLDriverErr, ok := err.(*mysql.MySQLError); ok && mySQLDriverErr.Number == DuplicateEntryErrCode {
			return errors.WithStack(ErrDuplicateKey)
		}
		return errors.WithStack(err)
	}
	return nil
}

func (d *DOBase) Insert(db *gorm.DB) error {
	if err := db.Create(d.DataObjecter).Error; err != nil {
		if mySQLDriverErr, ok := err.(*mysql.MySQLError); ok &&
			mySQLDriverErr.Number == DuplicateEntryErrCode {
			return errors.WithStack(ErrDuplicateKey)
		}
		return errors.WithStack(err)
	}
	return nil
}

func (d *DOBase) Delete(db *gorm.DB) error {
	if err := db.Delete(d.DataObjecter).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
