package permission

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/internal/pkg/utils"
	"neuronet/internal/pkg/utils/mapper"
	"neuronet/pkg/errors"
	"neuronet/pkg/log"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(c *gin.Context, args *v1.PermissionCreateArgs) (*v1.PermissionCreateReply, error)
	Delete(c *gin.Context, args *v1.PermissionDeleteArgs) error
	GetList(c *gin.Context, args *v1.PermissionListArgs) (*v1.PermissionListReply, error)
	Update(c *gin.Context, args *v1.PermissionUpdateArgs) error
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
	Page              int
	PageSize          int
	Name              string
	ParentID          int64
	Resource          string
	permissionRootIDs []int64
}

func (s *service) Create(c *gin.Context, args *v1.PermissionCreateArgs) (*v1.PermissionCreateReply, error) {

	tx := s.db.Begin()

	cnt, err := s.store.Permission().GetCntBy(c, tx, s.store.WithName(args.Name))
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("get permission cnt error: %v", err)
		return nil, err
	}
	if cnt != 0 {
		tx.Rollback()
		log.C(c).Warnf("permission existed")
		return nil, errors.WithCode(code.ErrDataExisted, "data existed")
	}

	permission := &model.Permission{
		Name:     args.Name,
		Resource: args.Resource,
		ParentID: args.ParentID,
	}

	err = s.store.Permission().Create(c, tx, permission)
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("create permission error: %v", err)
		return nil, err
	}

	newPermission := new(model.Permission)

	if args.ParentID == 0 {
		log.C(c).Infof("parentID 为 0")
		newPermission.RootID = permission.ID
		err = s.store.Permission().Updates(c, tx, newPermission, s.store.WithID(permission.ID))
		if err != nil {
			tx.Rollback()
			log.C(c).Warnf("update permission error: %v", err)
			return nil, err
		}
		log.C(c).Infof("权限创建成功")
		return &v1.PermissionCreateReply{
			MetaID: v1.MetaID{
				ID: permission.ID,
			},
		}, nil
	}

	log.C(c).Infof("parentID 不为 0")

	parentPermission, err := s.store.Permission().GetBy(c, tx, s.store.WithID(args.ParentID))
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("get parent permission error: %v", err)
		return nil, err
	}

	newPermission.RootID = parentPermission.RootID
	log.C(c).Infof("newPermission.RootID: %v", newPermission.RootID)

	err = s.store.Permission().Updates(c, tx, newPermission, s.store.WithID(permission.ID))
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("update permission error: %v", err)
		return nil, err
	}
	tx.Commit()
	log.C(c).Infof("权限创建成功")

	return &v1.PermissionCreateReply{
		MetaID: v1.MetaID{
			ID: permission.ID,
		},
	}, nil
}

func (s *service) Delete(c *gin.Context, args *v1.PermissionDeleteArgs) error {
	var (
		detailsReply  []v1.PermissionDetailReply
		permissionIDs []int64
		rootID        int64
	)
	tx := s.db.Begin()
	permission, err := s.store.Permission().GetBy(c, tx, s.store.WithID(args.ID))
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf(" 查询 第 %v permission 失败：%v", args.ID, err)
		return err
	}
	rootID = permission.RootID
	permissions, err := s.store.Permission().GetListBy(c, tx, s.store.WithRootID(rootID))
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf(" 查询 rootID %v permission 失败：%v", rootID, err)
		return err
	}
	log.C(c).Infof("permissions: %v", permissions)
	for _, permission := range permissions {
		detailsReply = append(detailsReply, v1.PermissionDetailReply{
			MetaID: v1.MetaID{
				ID: permission.ID,
			},
			ParentID: permission.ParentID,
		})
	}

	detailsReply = getTreeRecursive(detailsReply, args.ID)
	// 得到扁平结构ids
	permissionIDs = flattenTree(detailsReply)

	// 合并 childrenIDs 和 id，方便一起删除
	permissionIDs = append(permissionIDs, args.ID)

	err = s.store.Permission().DeleteBy(c, tx, s.store.WithIDs(permissionIDs))
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf(" 删除 permission %v 失败：%v", args.ID, err)
		return err
	}
	log.C(c).Infof(" 删除 permission %v 成功", args.ID)

	err = s.store.RolePermission().DeleteBy(c, tx, s.store.RolePermission().WithPermissionIDs(permissionIDs))
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf(" 删除 role permission bind 失败：%v", err)
		return err
	}
	log.C(c).Infof(" 删除 role permission bind 成功")

	tx.Commit()
	log.C(c).Infof(" 删除 permission %v 成功", args.ID)
	return nil

}

func (s *service) GetList(c *gin.Context, args *v1.PermissionListArgs) (*v1.PermissionListReply, error) {
	var (
		permissionDetails []v1.PermissionDetailReply
		permissionRootIDs []int64
	)

	query := s.listQuery(queryArgs{
		Page:     args.Page,
		PageSize: args.PageSize,
		Name:     args.Name,
		ParentID: args.ParentID,
		Resource: args.Resource,
	})
	permissions, err := s.store.Permission().GetListBy(c, s.db, query...)
	if err != nil {
		log.C(c).Warnf("get permission list error: %v", err)
		return nil, err
	}

	for _, permission := range permissions {
		permissionRootIDs = append(permissionRootIDs, permission.RootID)
	}
	permissionRootIDs = utils.RemoveRepeatedIDs(permissionRootIDs)

	query = s.listQuery(queryArgs{
		Name:              args.Name,
		Resource:          args.Resource,
		permissionRootIDs: permissionRootIDs,
	})

	permissions, err = s.store.Permission().GetListBy(c, s.db, query...)
	if err != nil {
		return nil, err
	}

	for _, permission := range permissions {
		detail := mapper.PermissionBoMapper(permission)
		permissionDetails = append(permissionDetails, detail)
	}
	permissionDetails = getTreeRecursive(permissionDetails, args.ParentID)

	query = s.listQuery(queryArgs{
		Name:     args.Name,
		Resource: args.Resource,
		ParentID: args.ParentID,
	})
	cnt, err := s.store.Permission().GetCntBy(c, s.db, query...)
	if err != nil {
		log.C(c).Warnf("get permission cnt error: %v", err)
		return nil, err
	}

	return &v1.PermissionListReply{
		MetaPage: v1.MetaPage{
			Page:     args.Page,
			PageSize: args.PageSize,
		},
		MetaTotalCnt: v1.MetaTotalCnt{
			TotalCnt: cnt,
		},
		List: permissionDetails,
	}, nil
}

func (s *service) Update(c *gin.Context, args *v1.PermissionUpdateArgs) error {

	newPermission := new(model.Permission)
	newPermission.Name = args.Name
	newPermission.Resource = args.Resource
	return s.store.Permission().Updates(c, s.db, newPermission, s.store.WithID(args.ID))

}

func (s *service) listQuery(args queryArgs) []store.DBOptions {
	var query []store.DBOptions
	if args.Resource != "" {
		query = append(query, s.store.Permission().WithResource(args.Resource))
	}
	if args.Name != "" {
		query = append(query, s.store.WithNameLike(args.Name))
	}
	if args.ParentID != 0 {
		query = append(query, s.store.WithParentID(args.ParentID))
	}
	if args.Page > 0 && args.PageSize > 0 {
		query = append(query, s.store.WithPage(args.Page, args.PageSize))
	}
	if len(args.permissionRootIDs) > 0 {
		query = append(query, s.store.WithRootIDs(args.permissionRootIDs))
	}
	return query
}

func getTreeRecursive(list []v1.PermissionDetailReply, parentId int64) []v1.PermissionDetailReply {
	res := make([]v1.PermissionDetailReply, 0)
	for _, v := range list {
		if v.ParentID == parentId {
			v.Children = getTreeRecursive(list, v.ID)
			res = append(res, v)
		}
	}
	return res
}

func flattenTree(list []v1.PermissionDetailReply) []int64 {
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
