package mysql

import (
	"neuronet/internal/neuronetserver/store"
)

var datastore *DBDatastore

var _ store.Factory = (*DBDatastore)(nil)

func GetDataStore() *DBDatastore {
	return datastore
}

func NewMysqlDatastore() *DBDatastore {
	datastore = &DBDatastore{}
	return datastore
}

type DBDatastore struct {
}

func (d *DBDatastore) Repository() store.IRepository {
	return newRepository()
}

func (d *DBDatastore) RepositoryImage() store.IRepositoryImage {
	return newRepositoryImage()
}

func (d *DBDatastore) RepositoryImageTag() store.IRepositoryImageTag {
	return newRepositoryImageTag()
}

func (d *DBDatastore) Cluster() store.ICluster {
	return newCluster()
}

func (d *DBDatastore) Task() store.ITask {
	return newTask()
}

func (d *DBDatastore) TaskResource() store.ITaskResource {
	return newTaskResource()
}

func (d *DBDatastore) UserRole() store.IUserRole {
	return newUserRole()
}

func (d *DBDatastore) RolePermission() store.IRolePermission {
	return newRolePermission()
}

func (d *DBDatastore) User() store.IUser {
	return newUser()
}

func (d *DBDatastore) Role() store.IRole {
	return newRole()
}

func (d *DBDatastore) Permission() store.IPermission {
	return newPermission()
}
