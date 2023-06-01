package model

import "gorm.io/plugin/soft_delete"

// ResourceType 可以建表管理
type ResourceType int64

const (
	ResourceTrainingJob ResourceType = iota + 1
	ResourceTensorboard
	ResourceJupyter
)

func (r ResourceType) String() string {
	switch r {
	case ResourceTrainingJob:
		return "training_job"
	case ResourceTensorboard:
		return "tensorboard"
	case ResourceJupyter:
		return "jupyter"
	}
	return ""
}

type TaskResource struct {
	BaseModel
	TaskID    int64                 `gorm:"type:int;column:task_id;comment:任务id"`
	ImageID   int64                 `gorm:"type:int;column:image_id;comment:镜像id"`
	Name      string                `gorm:"uniqueIndex:idx_name_deleted;type:varchar(256);not null;column:name;comment:资源名字"`
	Port      int32                 `gorm:"type:int;default 0;column:port;comment:资源端口号"`
	Type      int64                 `gorm:"type:int;not null;column:type;comment:1 表示 training_job，2 表示 tensorboard 3 表示 jupyter"`
	Cpu       int64                 `gorm:"type:int;default:0;column:cpu;comment:cpu核数"`
	Gpu       int64                 `gorm:"type:int;default:0;column:gpu;comment:gpu数量"`
	Memory    int64                 `gorm:"type:int;default:0;column:memory;comment:memory占用大小 Gi为单位"`
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_name_deleted"`
}

func (t TaskResource) TableName() string {
	return "task_resource"
}

type TaskResourceDo struct {
	TaskResource
}
