package mysql

import (
	"fmt"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/store"
)

func (d *DBDatastore) WithSelect(columns ...string) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(columns)
	}
}

func (d *DBDatastore) WithIDs(ids []int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", ids)
	}
}

func (d *DBDatastore) WithID(id int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (d *DBDatastore) WithName(name string) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (d *DBDatastore) WithNameLike(name string) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("name like %q ", "%"+name+"%"))
	}
}

func (d *DBDatastore) WithPage(page, pageSize int) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (d *DBDatastore) WithPreload(columnName string, args ...interface{}) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(columnName, args...)
	}
}

func (d *DBDatastore) WithDescIDOrder() store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("id desc")
	}
}

func (d *DBDatastore) WithAscIDOrder() store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("id asc")
	}
}

func (d *DBDatastore) WithParentID(id int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("parent_id = ?", id)
	}
}

func (d *DBDatastore) WithParentIDs(ids []int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("parent_id IN ?", ids)
	}
}

func (d *DBDatastore) WithRootID(id int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("root_id = ?", id)
	}
}

func (d *DBDatastore) WithRootIDs(ids []int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("root_id IN ?", ids)
	}
}

func (d *DBDatastore) WithSystemID(systemID int64) store.DBOptions {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("system_id = ?", systemID)
	}
}
