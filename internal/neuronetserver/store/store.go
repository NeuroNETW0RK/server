package store

type Factory interface {
	IOptionBase

	User() IUser
	Role() IRole
	Permission() IPermission
	UserRole() IUserRole
	RolePermission() IRolePermission

	Task() ITask
	TaskResource() ITaskResource
	Image() IImage

	Cluster() ICluster
}
