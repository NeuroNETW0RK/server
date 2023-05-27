package v1

type UserRegisterArgs struct {
	Account string
	MetaName
	Password string
	RoleIDs  []int64
}

type UserLoginArgs struct {
	Account  string
	Password string
}

type UserLoginReply struct {
	TokenHead string `json:"token_head"`
	Token     string `json:"token"`
	MetaID
}

type UserDeleteArgs struct {
	MetaID
}

type UserDetailArgs struct {
	MetaID
}

type UserDetailReply struct {
	MetaTime
	ID          int64                   `json:"id"`
	Name        string                  `json:"name"`
	Account     string                  `json:"account"`
	Roles       []RoleDetailReply       `json:"roles,omitempty"`
	Permissions []PermissionDetailReply `json:"permissions,omitempty"`
}

type UserListArgs struct {
	MetaPage
}

type UserListReply struct {
	MetaPage
	List []UserDetailReply `json:"list"`
	MetaTotalCnt
}

type UserUpdateArgs struct {
	MetaID
	Account string
	MetaName
	Password string
	RoleIDs  []int64
}
