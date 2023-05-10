package model

import "gorm.io/plugin/soft_delete"

type TaskImage struct {
	BaseModel
	Repository string                `gorm:"uniqueIndex:idx_rep_tag_deleted;type:varchar(512);not null;column:repository;comment:镜像仓库"`
	Tag        string                `gorm:"uniqueIndex:idx_rep_tag_deleted;type:varchar(512);not null;column:tag;comment:镜像tag"`
	Cpu        int64                 `gorm:"type:int;default:0;column:cpu;comment:cpu核数"`
	Gpu        int64                 `gorm:"type:int;default:0;column:gpu;comment:gpu数量"`
	Memory     int64                 `gorm:"type:int;default:0;column:memory;comment:memory占用大小 Gi为单位"`
	DeletedAt  soft_delete.DeletedAt `gorm:"uniqueIndex:idx_rep_tag_deleted"`
}

func (i TaskImage) TableName() string {
	return "task_image"
}

type TaskImageBo struct {
	TaskImage
}
