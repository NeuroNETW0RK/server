package store

import "gorm.io/gorm"

type DBOptions = func(db *gorm.DB) *gorm.DB

type IOptionBase interface {
	WithID(id int64) DBOptions
	WithIDs(ids []int64) DBOptions
	WithName(name string) DBOptions
	WithNameLike(name string) DBOptions
	WithPreload(columnName string, args ...interface{}) DBOptions
	WithPage(page, pageSize int) DBOptions
	WithDescIDOrder() DBOptions
	WithAscIDOrder() DBOptions
	WithParentID(id int64) DBOptions
	WithParentIDs(ids []int64) DBOptions
	WithRootID(id int64) DBOptions
	WithRootIDs(ids []int64) DBOptions
	WithSystemID(systemID int64) DBOptions
	WithSelect(columns ...string) DBOptions
}
