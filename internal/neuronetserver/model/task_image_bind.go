package model

type TaskImageBind struct {
	TaskID  int64 `gorm:"column:task_id"`
	ImageID int64 `gorm:"column:image_id"`
}

func (t TaskImageBind) TableName() string {
	return "task_image_bind"
}
