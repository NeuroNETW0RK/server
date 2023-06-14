package image

import (
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var _ BuildService = (*buildService)(nil)

type BuildService interface {
	Create(c *gin.Context, args *v1.ImageCreateArgs) (*v1.ImageCreateReply, error)
	Delete(c *gin.Context, args *v1.ImageDeleteArgs) error
	GetList(c *gin.Context, args *v1.ImageListArgs) (*v1.ImageListReply, error)
	Update(c *gin.Context, args *v1.ImageUpdateArgs) error
	Reload(c *gin.Context, args *v1.ImageReloadArgs) error
}

type buildService struct {
	store store.Factory
	db    *gorm.DB
}

func NewBuildService(db *gorm.DB, store store.Factory) *buildService {
	return &buildService{
		store: store,
		db:    db,
	}
}

func (s *buildService) Create(c *gin.Context, args *v1.ImageCreateArgs) (*v1.ImageCreateReply, error) {
	cnt, err := s.store.Image().GetCntBy(c, s.db, s.store.WithName(args.Name))
	if err != nil {
		log.C(c).Warnf("get image cnt error: %v", err)
		return nil, err
	}
	if cnt != 0 {
		log.C(c).Warnf("image existed")
		return nil, errors.WithCode(code.ErrDataExisted, "data existed")
	}

	newImage := &model.ImageInfo{
		ImageName: args.Name,
	}

	err = s.store.Image().Create(c, s.db, newImage)
	if err != nil {
		log.C(c).Warnf("create image error: %v", err)
		return nil, err
	}

	return &v1.ImageCreateReply{
		MetaID: v1.MetaID{
			ID: newImage.ID,
		},
	}, nil
}

func (s *buildService) Delete(c *gin.Context, args *v1.ImageDeleteArgs) error {
	// 删除镜像前需要做个判断，镜像是否被引用
	err := s.store.Image().DeleteBy(c, s.db, s.store.WithID(args.ID))
	if err != nil {
		log.Warnf("delete image error: %v", err)
		return err
	}
	return nil
}

func (s *buildService) GetList(c *gin.Context, args *v1.ImageListArgs) (*v1.ImageListReply, error) {
	var imagesDetailReply []v1.ImageDetailReply
	listQueryArgs := queryArgs{
		Page:     args.Page,
		PageSize: args.PageSize,
		Name:     args.Name,
	}
	listQuery := s.listQuery(listQueryArgs)
	imageTagBos, err := s.store.ImageTag().GetListBy(c, s.db, listQuery...)
	if err != nil {
		log.C(c).Warnf("get image list error: %v", err)
		return nil, err
	}

	for _, bo := range imageTagBos {
		imagesDetailReply = append(imagesDetailReply, v1.ImageDetailReply{
			MetaID: v1.MetaID{
				ID: bo.ID,
			},
			MetaName: v1.MetaName{},
			MetaTime: v1.MetaTime{
				CreateTime: bo.CreatedAt,
				UpdateTime: bo.UpdatedAt,
			},
		})
	}

	cntQueryArgs := queryArgs{
		Name: args.Name,
	}
	cntQuery := s.listQuery(cntQueryArgs)
	cnt, err := s.store.Image().GetCntBy(c, s.db, cntQuery...)
	if err != nil {
		log.C(c).Warnf("get image cnt error: %v", err)
		return nil, err
	}

	return &v1.ImageListReply{
		MetaPage: v1.MetaPage{
			Page:     args.Page,
			PageSize: args.PageSize,
		},
		MetaTotalCnt: v1.MetaTotalCnt{
			TotalCnt: cnt,
		},
		List: imagesDetailReply,
	}, nil
}

func (s *buildService) Update(c *gin.Context, args *v1.ImageUpdateArgs) error {
	_, err := s.store.Image().GetBy(c, s.db, s.store.WithID(args.ID))
	if err != nil {
		log.Warnf("get image error: %v", err)
		return err
	}

	newImage := &model.ImageInfo{
		ImageName:    args.Name,
		ImageUserFor: args.Description,
	}

	err = s.store.Image().Updates(c, s.db, newImage, s.store.WithID(args.ID))
	if err != nil {
		log.C(c).Warnf("update image error: %v", err)
		return err
	}

	return nil
}

func (s *buildService) Reload(c *gin.Context, args *v1.ImageReloadArgs) error {
	_, err := s.store.Image().GetBy(c, s.db, s.store.WithID(args.ID))
	if err != nil {
		log.C(c).Warnf("get image error: %v", err)
		return err
	}
	if err != nil {
		log.C(c).Warnf("create clientSets error: %v", err)
		return errors.WithCode(code.ErrInternalServer, "create clientSets error")
	}
	return nil
}

func (s *buildService) listQuery(args queryArgs) []store.DBOptions {
	var query []store.DBOptions

	if args.Page > 0 && args.PageSize > 0 {
		query = append(query, s.store.WithPage(args.Page, args.PageSize))
	}
	if args.Name != "" {
		query = append(query, s.store.WithNameLike(args.Name))
	}
	return query
}
