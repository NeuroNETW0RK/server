package mapper

import (
	v1 "NeuroNET/internal/neuronetserver/dto/v1"
	"NeuroNET/internal/neuronetserver/model"
)

func PermissionBoMapper(permission model.PermissionBo) v1.PermissionDetailReply {
	return v1.PermissionDetailReply{
		MetaID: v1.MetaID{
			ID: permission.ID,
		},
		MetaName: v1.MetaName{
			Name: permission.Name,
		},
		MetaTime: v1.MetaTime{
			CreateTime: permission.CreatedAt,
			UpdateTime: permission.UpdatedAt,
		},
		ParentID: permission.ParentID,
		Children: nil,
		Resource: permission.Resource,
	}
}
