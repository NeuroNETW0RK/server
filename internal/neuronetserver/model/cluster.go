package model

type Cluster struct {
	BaseModel
	Name        string `gorm:"uniqueIndex:idx_name;type:varchar(255);not null;column:name;comment:集群名称"`
	ConfigPath  string `gorm:"type:varchar(255);not null;column:config_path;comment:集群配置文件路径"`
	Description string `gorm:"type:longText;column:description;comment:集群描述"`
}

func (c Cluster) TableName() string {
	return "cluster"
}

type ClusterBo struct {
	Cluster
}
