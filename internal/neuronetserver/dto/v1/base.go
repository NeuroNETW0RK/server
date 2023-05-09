package v1

import "time"

type MetaID struct {
	ID int64 `form:"id" json:"id" binding:"required"`
}

type MetaName struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type MetaSystemID struct {
	SystemID int64 `form:"system_id" json:"system_id"`
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

type MetaPipelineStatus struct {
	ScenarioStatusID   int64  `json:"scenario_status_id"`
	ScenarioStatus     string `json:"scenario_status"`
	AgentStatusID      int64  `json:"agent_status_id"`
	AgentStatus        string `json:"agent_status"`
	TrainingStatusID   int64  `json:"training_status_id"`
	TrainingStatus     string `json:"training_status"`
	EvaluationStatusID int64  `json:"evaluation_status_id"`
	EvaluationStatus   string `json:"evaluation_status"`
}
