package v1

type RoleCreateArgs struct {
	MetaName
	ParentID      int64
	PermissionIDs []int64
}

type RoleCreateReply struct {
	MetaID
}

type RoleDeleteArgs struct {
	MetaID
}

type RoleListArgs struct {
	MetaPage
	MetaName
	ParentID int64
}

type RoleListReply struct {
	MetaPage
	MetaTotalCnt
	List []RoleDetailReply `json:"list"`
}

type RoleDetailReply struct {
	MetaID
	MetaName
	MetaTime
	ParentID    int64                   `json:"parent_id"`
	Children    []RoleDetailReply       `json:"children,omitempty"`
	Permissions []PermissionDetailReply `json:"permissions,omitempty"`
}

type RoleUpdateArgs struct {
	MetaID
	MetaName
	PermissionIDs []int64
}
