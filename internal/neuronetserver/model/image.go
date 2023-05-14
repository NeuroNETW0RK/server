package model

import "gorm.io/plugin/soft_delete"

type Image struct {
	BaseModel
	Repository string                `gorm:"uniqueIndex:idx_rep_tag_deleted;type:varchar(512);not null;column:repository;comment:镜像仓库"`
	Tag        string                `gorm:"uniqueIndex:idx_rep_tag_deleted;type:varchar(512);not null;column:tag;comment:镜像tag"`
	DeletedAt  soft_delete.DeletedAt `gorm:"uniqueIndex:idx_rep_tag_deleted"`
}

func (i Image) TableName() string {
	return "image"
}

type ImageBo struct {
	Image
}
