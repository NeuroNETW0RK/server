package v1

type PermissionCreateArgs struct {
	MetaName
	Resource string
	ParentID int64
}

type PermissionCreateReply struct {
	MetaID
}

type PermissionDeleteArgs struct {
	MetaID
}

type PermissionListArgs struct {
	MetaPage
	MetaName
	ParentID int64
	Resource string
}

type PermissionListReply struct {
	MetaPage
	MetaTotalCnt
	List []PermissionDetailReply `json:"list"`
}

type PermissionDetailReply struct {
	MetaID
	MetaName
	MetaTime
	ParentID int64                   `json:"parent_id"`
	Children []PermissionDetailReply `json:"children,omitempty"`
	Resource string                  `json:"resource"`
}

type PermissionUpdateArgs struct {
	MetaID
	MetaName
	Resource string
}
