package meta

type CreateOptions struct {
	Namespace string `uri:"namespace" binding:"required"`
}

type UpdateOptions struct {
	Namespace  string `uri:"namespace" binding:"required"`
	ObjectName string `uri:"object_name" binding:"required"`
}

type DeleteOptions struct {
	Namespace  string `uri:"namespace" binding:"required"`
	ObjectName string `uri:"object_name" binding:"required"`
}

type GetOptions struct {
	Namespace  string `uri:"namespace" binding:"required"`
	ObjectName string `uri:"object_name" binding:"required"`
}

type RestartOptions struct {
	Namespace  string `uri:"namespace" binding:"required"`
	ObjectName string `uri:"object_name" binding:"required"`
}

type ListOptions struct {
	Namespace    string `uri:"namespace" binding:"required"`
	ListSelector `json:",inline"`
}

type LogOptions struct {
	Namespace     string `uri:"namespace" binding:"required"`
	ObjectName    string `uri:"object_name" binding:"required"`
	ContainerName string `form:"container_name"`
}

type TopOptions struct {
	Namespace     string `uri:"namespace" binding:"required"`
	ObjectName    string `uri:"object_name" binding:"required"`
	ContainerName string `form:"container_name"`
}

type TerminalOptions struct {
	Namespace  string `uri:"namespace" binding:"required"`
	ObjectName string `uri:"object_name" binding:"required"`
	Container  string `form:"container"`
}

type ListSelector struct {
	Field    string            `form:"field"`
	Label    string            `form:"label"`
	IndexMap map[string]string `form:"index_map"`

	Page     int `form:"page"`      // 页数
	PageSize int `form:"page_size"` // 每页数量
}
