package role

import (
	v1 "NeuroNET/internal/neuronetserver/dto/v1"
	"NeuroNET/internal/neuronetserver/model"
	"NeuroNET/internal/neuronetserver/store"
	"NeuroNET/internal/pkg/code"
	"NeuroNET/internal/pkg/utils"
	"NeuroNET/internal/pkg/utils/mapper"
	"NeuroNET/pkg/errors"
	"NeuroNET/pkg/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(c *gin.Context, args *v1.RoleCreateArgs) (*v1.RoleCreateReply, error)
	Delete(c *gin.Context, args *v1.RoleDeleteArgs) error
	GetList(c *gin.Context, args *v1.RoleListArgs) (*v1.RoleListReply, error)
	Update(c *gin.Context, args *v1.RoleUpdateArgs) error
}

type service struct {
	store store.Factory
	db    *gorm.DB
}

func NewService(db *gorm.DB, store store.Factory) *service {
	return &service{
		store: store,
		db:    db,
	}
}

type queryArgs struct {
	Page     int
	PageSize int
	Name     string
	ParentID int64
	RootIDs  []int64
	RoleIDs  []int64
}

func (s *service) Create(c *gin.Context, args *v1.RoleCreateArgs) (*v1.RoleCreateReply, error) {
	var binds []model.PermissionRole

	tx := s.db.Begin()

	cnt, err := s.store.Role().GetCntBy(c, tx, s.store.WithName(args.Name))
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("get role cnt error: %v", err)
		return nil, err
	}
	if cnt != 0 {
		tx.Rollback()
		log.C(c).Warnf("role existed")
		return nil, errors.WithCode(code.ErrDataExisted, "data existed")
	}

	role := &model.Role{
		Name:     args.Name,
		ParentID: args.ParentID,
	}
	err = s.store.Role().Create(c, tx, role)
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("create role error: %v", err)
		return nil, err
	}

	newRole := new(model.Role)

	if args.ParentID == 0 {
		newRole.RootID = role.ID
		err = s.store.Role().Updates(c, tx, newRole, s.store.WithID(role.ID))
		if err != nil {
			tx.Rollback()
			log.C(c).Warnf("update role error: %v", err)
			return nil, err
		}
	} else {
		parentRole, err := s.store.Role().GetBy(c, tx, s.store.WithID(args.ParentID))
		if err != nil {
			tx.Rollback()
			log.C(c).Warnf("get role error: %v", err)
			return nil, err
		}

		newRole.RootID = parentRole.RootID

		err = s.store.Role().Updates(c, tx, newRole, s.store.WithID(role.ID))
		if err != nil {
			tx.Rollback()
			log.C(c).Warnf("update role error: %v", err)
			return nil, err
		}
	}

	if len(args.PermissionIDs) == 0 {
		tx.Commit()
		return &v1.RoleCreateReply{
			MetaID: v1.MetaID{ID: role.ID},
		}, nil
	}

	for _, permissionID := range args.PermissionIDs {
		binds = append(binds, model.PermissionRole{
			RoleID:       role.ID,
			PermissionID: permissionID,
		})
	}

	err = s.store.RolePermission().Create(c, tx, binds)
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("create role permission error: %v", err)
		return nil, err
	}

	tx.Commit()
	log.C(c).Infof("create role success: %v", role.ID)
	return &v1.RoleCreateReply{
		MetaID: v1.MetaID{ID: role.ID},
	}, nil
}

func (s *service) Delete(c *gin.Context, args *v1.RoleDeleteArgs) error {
	var (
		detailsReply []v1.RoleDetailReply
		roleIDs      []int64
		rootID       int64
	)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		log.C(c).Infof("开始删除 role id: %v", args.ID)
		role, err := s.store.Role().GetBy(c, tx, s.store.WithID(args.ID))
		if err != nil {
			log.C(c).Warnf("查询 第 %v role 失败：%v", args.ID, err)
			return err
		}

		rootID = role.RootID
		roles, err := s.store.Role().GetListBy(c, tx, s.store.WithRootID(rootID))
		if err != nil {
			log.C(c).Warnf("查询 rootID %v role 失败：%v", rootID, err)
			return err
		}
		log.C(c).Infof("查询 rootID %v role 成功：%v", rootID, roles)
		for _, role := range roles {
			detailsReply = append(detailsReply, v1.RoleDetailReply{
				MetaID: v1.MetaID{
					ID: role.ID,
				},
				ParentID: role.ParentID,
			})
		}

		detailsReply = getTreeRecursive(detailsReply, args.ID)
		// 得到扁平结构ids
		roleIDs = flattenTree(detailsReply)

		// 合并 childrenIDs 和 id，方便一起删除
		roleIDs = append(roleIDs, args.ID)

		err = s.store.Role().DeleteBy(c, tx, s.store.WithIDs(roleIDs))
		if err != nil {
			log.C(c).Warnf("删除 role 失败：%v", err)
			return err
		}

		err = s.store.RolePermission().DeleteBy(c, tx, s.store.RolePermission().WithRoleIDs([]int64{args.ID}))
		if err != nil {
			log.C(c).Warnf("删除 role permission 失败：%v", err)
			return err
		}
		err = s.store.UserRole().DeleteBy(c, tx, s.store.UserRole().WithRoleIDs([]int64{args.ID}))
		if err != nil {
			log.C(c).Warnf("删除 user role 失败：%v", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.C(c).Warnf("删除 role 失败：%v", err)
		return err
	}

	log.C(c).Infof("删除 role 成功 id：%v", args.ID)
	return nil
}

func (s *service) GetList(c *gin.Context, args *v1.RoleListArgs) (*v1.RoleListReply, error) {
	var (
		roleDetails []v1.RoleDetailReply
		roleRootIDs []int64
	)

	query := queryArgs{
		Page:     args.Page,
		PageSize: args.PageSize,
		Name:     args.Name,
		ParentID: args.ParentID,
	}
	roles, err := s.store.Role().GetListBy(c, s.db, s.listQuery(query)...)
	if err != nil {
		log.C(c).Warnf("查询 role 失败：%v", err)
		return nil, err
	}
	// 先取出 rootID
	for _, role := range roles {
		roleRootIDs = append(roleRootIDs, role.RootID)
	}
	roleRootIDs = utils.RemoveRepeatedIDs(roleRootIDs)

	query = queryArgs{
		Name:    args.Name,
		RootIDs: roleRootIDs,
	}
	argsQuery := s.listQuery(query)

	argsQuery = append(argsQuery, s.store.WithPreload("Permissions"))

	roles, err = s.store.Role().GetListBy(c, s.db, argsQuery...)
	if err != nil {
		log.C(c).Warnf("查询 role 失败：%v", err)
		return nil, err
	}

	for _, role := range roles {
		roleDetails = append(roleDetails, mapper.RoleBoMapper(role))
	}

	roleDetails = getTreeRecursive(roleDetails, args.ParentID)

	query = queryArgs{
		Name:     args.Name,
		ParentID: args.ParentID,
	}
	cnt, err := s.store.Role().GetCntBy(c, s.db, s.listQuery(query)...)
	if err != nil {
		log.C(c).Warnf("查询 role 失败：%v", err)
		return nil, err
	}

	return &v1.RoleListReply{
		MetaPage: v1.MetaPage{
			Page:     args.Page,
			PageSize: args.PageSize,
		},
		MetaTotalCnt: v1.MetaTotalCnt{
			TotalCnt: cnt,
		},
		List: roleDetails,
	}, nil
}

func (s *service) Update(c *gin.Context, args *v1.RoleUpdateArgs) error {
	var binds []model.PermissionRole

	err := s.db.Transaction(func(tx *gorm.DB) error {
		newRole := new(model.Role)
		newRole.Name = args.Name

		err := s.store.Role().Updates(c, tx, newRole, s.store.WithID(args.ID))
		if err != nil {
			log.C(c).Warnf("update role error: %v", err)
			return err
		}

		err = s.store.RolePermission().DeleteBy(c, tx, s.store.WithID(args.ID))
		if err != nil {
			log.C(c).Warnf("delete role permission error: %v", err)
			return err
		}

		if len(args.PermissionIDs) == 0 {
			log.C(c).Infof("permission ids is empty")
			return nil
		}

		for _, permissionID := range args.PermissionIDs {
			binds = append(binds, model.PermissionRole{
				RoleID:       args.ID,
				PermissionID: permissionID,
			})
		}

		err = s.store.RolePermission().Create(c, tx, binds)
		if err != nil {
			log.C(c).Warnf("create role permission error: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.C(c).Warnf("update role error: %v", err)
		return err
	}
	log.C(c).Infof("update role success: %v", args.ID)
	return nil
}

func (s *service) listQuery(args queryArgs) []store.DBOptions {
	var query []store.DBOptions

	if args.Name != "" {
		query = append(query, s.store.WithNameLike(args.Name))
	}
	if args.ParentID != 0 {
		query = append(query, s.store.WithParentID(args.ParentID))
	}
	if args.Page > 0 && args.PageSize > 0 {
		query = append(query, s.store.WithPage(args.Page, args.PageSize))
	}
	if len(args.RootIDs) > 0 {
		query = append(query, s.store.WithRootIDs(args.RootIDs))
	}
	if len(args.RoleIDs) > 0 {
		query = append(query, s.store.WithRootIDs(args.RoleIDs))
	}
	return query
}

func getTreeRecursive(list []v1.RoleDetailReply, parentID int64) []v1.RoleDetailReply {
	res := make([]v1.RoleDetailReply, 0)
	for _, v := range list {
		if v.ParentID == parentID {
			v.Children = getTreeRecursive(list, v.ID)
			res = append(res, v)
		}
	}
	return res
}

func flattenTree(list []v1.RoleDetailReply) []int64 {
	if len(list) == 0 {
		return nil
	}
	var resultID []int64
	for _, v := range list {
		resultID = append(resultID, v.ID)
		resultID = append(resultID, flattenTree(v.Children)...)

	}
	return resultID
}
