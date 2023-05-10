package mapper

import (
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/neuronetserver/model"
)

func RoleBoMapper(role model.RoleBo) v1.RoleDetailReply {

	var permissionsDetail []v1.PermissionDetailReply
	for _, permission := range role.Permissions {
		permissionsDetail = append(permissionsDetail, PermissionBoMapper(permission))
	}

	return v1.RoleDetailReply{
		MetaID: v1.MetaID{
			ID: role.ID,
		},
		MetaName: v1.MetaName{
			Name: role.Name,
		},
		MetaTime: v1.MetaTime{
			CreateTime: role.CreatedAt,
			UpdateTime: role.UpdatedAt,
		},
		ParentID:    role.ParentID,
		Children:    nil,
		Permissions: permissionsDetail,
	}
}
