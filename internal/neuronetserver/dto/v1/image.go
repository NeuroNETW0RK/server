package v1

type ImageCreateArgs struct {
	MetaName
	Description    string `form:"description" json:"description"`
	KubeConfigPath string `form:"kube_config_path" json:"kube_config_path" binding:"required"`
}

type ImageCreateReply struct {
	MetaID
}

type ImageDeleteArgs struct {
	MetaID
}

type ImageListArgs struct {
	MetaPage
	MetaName
}

type ImageListReply struct {
	MetaPage
	MetaTotalCnt
	List []ImageDetailReply `json:"list"`
}

type ImageDetailReply struct {
	MetaID
	MetaName
	MetaTime
}

type ImageUpdateArgs struct {
	MetaID
	MetaName
	Description string `form:"description" json:"description"`
}

type ImageReloadArgs struct {
	MetaID
}
