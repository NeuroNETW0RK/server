package model

import "gorm.io/plugin/soft_delete"

type TaskStatus int64

const (
	TaskReady TaskStatus = iota + 1
	TaskRunning
	TaskEnd
)

func (t TaskStatus) String() string {
	switch t {
	case TaskReady:
		return "已准备"
	case TaskRunning:
		return "进行中"
	case TaskEnd:
		return "已停止"
	}
	return ""
}

type Task struct {
	BaseModel
	Name           string                `gorm:"uniqueIndex:idx_name_uuid_deleted;type:varchar(256);not null;column:name;comment:任务名字"`
	Uuid           string                `gorm:"uniqueIndex:idx_name_uuid_deleted;type:varchar(256);not null;column:uuid;comment:在k8s中显示的任务id"`
	StatusID       int64                 `gorm:"type:int;not null;column:status;comment:1 表示 ready，2 表示 running，3 表示 end"`
	TrainStartTime int64                 `gorm:"type:int;default:0;column:train_start_time;comment:训练开启时间"`
	TrainEndTime   int64                 `gorm:"type:int;default:0;column:train_end_time;comment:训练结束时间"`
	Desc           string                `gorm:"type:longText;column:desc;comment:任务描述"`
	DeletedAt      soft_delete.DeletedAt `gorm:"uniqueIndex:idx_name_uuid_deleted"`
}

func (t Task) TableName() string {
	return "task"
}

type TaskDo struct {
	Task
}
