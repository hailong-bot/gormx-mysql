package gormxmysql

import (
	"reflect"

	"github.com/go-sql-driver/mysql"
	"github.com/hailong-bot/gormx"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ModelBase struct {
	DataObjecter gormx.DataObjecter
}

func (m *ModelBase) GetByID(db *gorm.DB, id int64) (gormx.DataObjecter, error) {
	dataObjectType := reflect.TypeOf(m.DataObjecter)
	for dataObjectType.Kind() == reflect.Ptr {
		dataObjectType = dataObjectType.Elem()
	}

	dataObjecterValue := reflect.New(dataObjectType)
	result := dataObjecterValue.Interface().(gormx.DataObjecter)

	if dataObjecterValue.Elem().Kind() == reflect.Struct {
		doerField := dataObjecterValue.Elem().FieldByName("gormx.DataObjecter")
		if doerField.IsValid() && doerField.CanSet() {
			doerField.Set(dataObjecterValue)
		}
	}

	if err := db.Take(result, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (m *ModelBase) GetByIDWithLock(db *gorm.DB, id int64, lock gormx.Lock) (gormx.DataObjecter, error) {
	// 1.生成该对象
	dataObjecterType := reflect.TypeOf(m.DataObjecter)
	for dataObjecterType.Kind() == reflect.Ptr {
		dataObjecterType = dataObjecterType.Elem()
	}
	dataObjecterValue := reflect.New(dataObjecterType)
	result := dataObjecterValue.Interface().(gormx.DataObjecter)

	// 2.为生成的对象设置 gormx.DataObjecter 值
	if dataObjecterValue.Elem().Kind() == reflect.Struct {
		doerField := dataObjecterValue.Elem().FieldByName("gormx.DataObjecter")
		if doerField.IsValid() && doerField.CanSet() {
			doerField.Set(dataObjecterValue)
		}
	}

	// 3.查找该对象
	query := db
	switch lock {
	case gormx.NoLock:
	case gormx.IS:
		query = query.Clauses(clause.Locking{Strength: clause.LockingStrengthShare})
	case gormx.IX:
		query = query.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
	}
	if err := query.Take(result, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (m *ModelBase) InsertBatch(db *gorm.DB, doList interface{}) error {
	if doList == nil {
		return nil
	}
	doListType := reflect.TypeOf(doList)
	if doListType.Kind() != reflect.Ptr {
		return errors.New("param expect a pointer of slice")
	} else if doListType.Elem().Kind() != reflect.Slice && doListType.Elem().Kind() != reflect.Array {
		return errors.New("param expect a pointer of slice")
	} else {
		doListValue := reflect.ValueOf(doList)
		if doListValue.Elem().Len() == 0 {
			return nil
		}
	}
	if err := db.Create(doList).Error; err != nil {
		if mySQLDriverErr, ok := err.(*mysql.MySQLError); ok &&
			mySQLDriverErr.Number == DuplicateEntryErrCode {
			return errors.WithStack(ErrDuplicateKey)
		}
		return errors.WithStack(err)
	}
	return nil
}

func (m *ModelBase) GetByConditions(db *gorm.DB, where string, values ...interface{}) (gormx.DataObjecter, error) {
	dataObjectType := reflect.TypeOf(m.DataObjecter)
	for dataObjectType.Kind() == reflect.Ptr {
		dataObjectType = dataObjectType.Elem()
	}
	dataObjecterValue := reflect.New(dataObjectType)
	result := dataObjecterValue.Interface().(gormx.DataObjecter)

	if dataObjecterValue.Elem().Kind() == reflect.Struct {
		doerField := dataObjecterValue.Elem().FieldByName("gormx.DataObjecter")
		if doerField.IsValid() && doerField.CanSet() {
			doerField.Set(dataObjecterValue)
		}
	}

	if err := db.Model(m.DataObjecter).Where(where, values...).Take(result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return result, nil
}

func (m *ModelBase) GetByConditionsWithLock(db *gorm.DB, lock gormx.Lock, where string, values ...interface{}) (gormx.DataObjecter, error) {
	// 1.生成该对象
	dataObjecterType := reflect.TypeOf(m.DataObjecter)
	for dataObjecterType.Kind() == reflect.Ptr {
		dataObjecterType = dataObjecterType.Elem()
	}
	dataObjecterValue := reflect.New(dataObjecterType)
	result := dataObjecterValue.Interface().(gormx.DataObjecter)

	// 2.为生成的对象设置 gormx.DataObjecter 值
	if dataObjecterValue.Elem().Kind() == reflect.Struct {
		doerField := dataObjecterValue.Elem().FieldByName("gormx.DataObjecter")
		if doerField.IsValid() && doerField.CanSet() {
			doerField.Set(dataObjecterValue)
		}
	}

	// 3.查找该对象
	query := db.Model(m.DataObjecter).Where(where, values...)
	switch lock {
	case gormx.NoLock:
	case gormx.IS:
		query = query.Clauses(clause.Locking{Strength: clause.LockingStrengthShare})
	case gormx.IX:
		query = query.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
	}
	if err := query.Take(result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func (m *ModelBase) List(db *gorm.DB, offset int, limit int, sortField string, sort gormx.Sort, where string, values ...interface{}) (
	gormx.DataObjecterList, error,
) {
	var result gormx.DataObjecterList

	// 1.生成该对象
	dataObjecterType := reflect.TypeOf(m.DataObjecter)
	dataObjecterPtrType := dataObjecterType
	for dataObjecterPtrType.Kind() == reflect.Ptr && dataObjecterPtrType.Elem().Kind() == reflect.Ptr {
		dataObjecterPtrType = dataObjecterPtrType.Elem()
	}
	resultValue := reflect.New(reflect.SliceOf(dataObjecterPtrType))
	resultSlice := resultValue.Interface()

	// 2.查找到该列表
	sortSQL := sortField + " " + sort.ToString()
	if err := db.Order(sortSQL).Offset(offset).Limit(limit).
		Model(m.DataObjecter).Where(where, values...).
		Find(resultSlice).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	// 3.将结果做类型转换，转为 dataObjecterType
	for i := 0; i < resultValue.Elem().Len(); i++ {
		result = append(
			result,
			resultValue.Elem().Index(i).Interface().(gormx.DataObjecter),
		)
	}

	// 4.为生成的对象设置 gormx.DataObjecter 值
	for i := range result {
		dataObjecterValue := reflect.ValueOf(result[i])
		if dataObjecterValue.Elem().Kind() == reflect.Struct {
			doerField := dataObjecterValue.Elem().FieldByName("gormx.DataObjecter")
			if doerField.IsValid() && doerField.CanSet() {
				doerField.Set(dataObjecterValue)
			}
		}
	}
	return result, nil
}
