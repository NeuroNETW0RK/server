package v1

import "time"

type MetaID struct {
	ID int64 `form:"id" json:"id" binding:"required"`
}

type MetaName struct {
	Name string `form:"name" json:"name"`
}

type MetaPage struct {
	Page     int `form:"page,default=1" json:"page"`
	PageSize int `form:"page_size,default=10" json:"page_size"`
}

type MetaTotalCnt struct {
	TotalCnt int64 `json:"total_cnt"`
}

type MetaTime struct {
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type MetaTypeID struct {
	TypeID int64 `form:"type_id" json:"type_id"`
}

type MetaUserName struct {
	UserName string `form:"user_name" json:"user_name"`
}

type MetaUserID struct {
	UserID int64 `form:"user_id" json:"user_id"`
}

type MetaIDName struct {
	ID   int64  `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
}
