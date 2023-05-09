package model

import "time"

type BaseModel struct {
	ID        int64     `gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time `gorm:"column:add_time;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime"`
}
