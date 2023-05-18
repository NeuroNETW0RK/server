package v1

type ClusterCreateArgs struct {
	MetaName
	Description    string `form:"description" json:"description"`
	KubeConfigPath string `form:"kube_config_path" json:"kube_config_path" binding:"required"`
}

type ClusterCreateReply struct {
	MetaID
}

type ClusterDeleteArgs struct {
	MetaID
}

type ClusterListArgs struct {
	MetaPage
	MetaName
}

type ClusterListReply struct {
	MetaPage
	MetaTotalCnt
	List []ClusterDetailReply `json:"list"`
}

type ClusterDetailReply struct {
	MetaID
	MetaName
	MetaTime
	KubeConfigPath string `json:"kube_config_path"`
}

type ClusterUpdateArgs struct {
	MetaID
	MetaName
	Description    string `form:"description" json:"description"`
	KubeConfigPath string `form:"kube_config_path" json:"kube_config_path"`
}

type ClusterReloadArgs struct {
	MetaID
}
