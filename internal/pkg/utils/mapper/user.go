package mapper

import (
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/neuronetserver/model"
)

func UserBoMapper(user model.UserBo) v1.UserDetailReply {
	var (
		permissionIDMap   = make(map[int64]model.PermissionBo)
		permissionsDetail []v1.PermissionDetailReply
		rolesDetail       []v1.RoleDetailReply
	)
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			permissionIDMap[permission.ID] = permission
		}
	}

	for _, permission := range permissionIDMap {
		permissionsDetail = append(permissionsDetail, PermissionBoMapper(permission))
	}

	for _, role := range user.Roles {
		// 不需要在role里面返回权限
		role.Permissions = nil
		rolesDetail = append(rolesDetail, RoleBoMapper(role))
	}

	return v1.UserDetailReply{
		MetaTime: v1.MetaTime{
			CreateTime: user.CreatedAt,
			UpdateTime: user.UpdatedAt,
		},
		ID:          user.ID,
		Name:        user.Name,
		Account:     user.Account,
		Roles:       rolesDetail,
		Permissions: permissionsDetail,
	}

}
