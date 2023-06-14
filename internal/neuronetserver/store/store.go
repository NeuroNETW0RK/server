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

	Repository() IRepository
	RepositoryImage() IRepositoryImage
	RepositoryImageTag() IRepositoryImageTag

	Cluster() ICluster
	Image() IImage
	ImageTag() IImageTag
	ImageBuild() IImageBuild
}
