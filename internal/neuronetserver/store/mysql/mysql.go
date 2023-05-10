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

func (d *DBDatastore) Task() store.ITask {
	return newTask()
}

func (d *DBDatastore) TaskResource() store.ITaskResource {
	return newTaskResource()
}

func (d *DBDatastore) TaskImage() store.ITaskImage {
	return newTaskImage()
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
