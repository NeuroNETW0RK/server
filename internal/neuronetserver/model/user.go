package model

import "gorm.io/plugin/soft_delete"

const (
	// UserID 用户id
	UserID = "user_id"
	// UserRoleGuest 对应数据库的游客id写死为1
	UserRoleGuest = iota + 1
)

type User struct {
	BaseModel
	Account   string                `gorm:"uniqueIndex:idx_account_deleted;type:varchar(64);not null;column:account;comment:账号"`
	Name      string                `gorm:"type:varchar(64);not null;column:name;comment:姓名"`
	Password  string                `gorm:"type:varchar(256);not null;column:password;comment:密码"`
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_account_deleted;column:deleted_at;comment:删除时间"`
}

func (u User) TableName() string {
	return "admin_user"
}

type UserDo struct {
	User
	Roles []RoleDo `gorm:"many2many:admin_user_role;foreignKey:ID;joinForeignKey:user_id;References:ID;JoinReferences:role_id"`
}

type UserRole struct {
	UserID int64 `gorm:"column:user_id"`
	RoleID int64 `gorm:"column:role_id"`
}

func (u UserRole) TableName() string {
	return "admin_user_role"
}

type Role struct {
	BaseModel
	Name      string                `gorm:"uniqueIndex:idx_role_deleted;type:varchar(64);not null;column:name;comment:角色名字"`
	ParentID  int64                 `gorm:"type:int;column:parent_id;comment:父角色id"`
	RootID    int64                 `gorm:"type:int;column:root_id;comment:根角色id"`
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_role_deleted;column:deleted_at;comment:删除时间"`
}

func (r Role) TableName() string {
	return "admin_role"
}

type RoleDo struct {
	Role
	Permissions []PermissionDo `gorm:"many2many:admin_permission_role;foreignKey:ID;joinForeignKey:role_id;References:ID;JoinReferences:permission_id;column:permissions;comment:权限"`
	Users       []User         `gorm:"many2many:admin_user_role;foreignKey:ID;joinForeignKey:role_id;References:ID;JoinReferences:user_id;column:users;comment:用户"`
}

type PermissionRole struct {
	RoleID       int64 `gorm:"type:int;column:role_id"`
	PermissionID int64 `gorm:"type:int;column:permission_id"`
}

func (p PermissionRole) TableName() string {
	return "admin_permission_role"
}

type Permission struct {
	BaseModel
	Name      string                `gorm:"type:varchar(256);not null;column:name;comment:权限名字"`
	Resource  string                `gorm:"uniqueIndex:idx_resource_deleted;type:varchar(256);not null;column:resource;comment:权限资源"`
	ParentID  int64                 `gorm:"type:int;column:parent_id;comment:父权限id"`
	RootID    int64                 `gorm:"type:int;column:root_id;comment:根权限id"`
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_resource_deleted;column:deleted_at;comment:删除时间"`
}

func (p Permission) TableName() string {
	return "admin_permission"
}

type PermissionDo struct {
	Permission
	Roles []Role `gorm:"many2many:admin_permission_role;foreignKey:ID;joinForeignKey:permission_id;References:ID;JoinReferences:role_id;column:roles;comment:角色"`
}
