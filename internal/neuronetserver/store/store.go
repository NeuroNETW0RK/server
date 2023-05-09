package store

type Factory interface {
	IOptionBase

	User() IUser
	Role() IRole
	Permission() IPermission
	UserRole() IUserRole
	RolePermission() IRolePermission
}
