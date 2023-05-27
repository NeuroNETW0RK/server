package user

import (
	"crypto/sha512"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"neuronet/internal/neuronetserver/configs"
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	auth "neuronet/internal/pkg/jwt"
	"neuronet/internal/pkg/utils/mapper"
	userutils "neuronet/internal/pkg/utils/user"
	"neuronet/pkg/errors"
	"neuronet/pkg/log"

	"strings"
)

var _ Service = (*service)(nil)

type Service interface {
	Register(c *gin.Context, args *v1.UserRegisterArgs) error
	Login(c *gin.Context, args *v1.UserLoginArgs) (*v1.UserLoginReply, error)
	Delete(c *gin.Context, args *v1.UserDeleteArgs) error
	GetDetail(c *gin.Context, args *v1.UserDetailArgs) (*v1.UserDetailReply, error)
	GetList(c *gin.Context, args *v1.UserListArgs) (*v1.UserListReply, error)
	Update(c *gin.Context, args *v1.UserUpdateArgs) error
}

type service struct {
	store store.Factory
	db    *gorm.DB
}

type queryArgs struct {
	Page     int
	PageSize int
	Account  string
}

func NewService(db *gorm.DB, store store.Factory) *service {
	return &service{
		store: store,
		db:    db,
	}
}

func (s *service) Register(c *gin.Context, args *v1.UserRegisterArgs) error {
	var binds []model.UserRole

	tx := s.db.Begin()

	newUser := &model.User{
		Account:  args.Account,
		Name:     args.Name,
		Password: userutils.EncodePassword(args.Password),
	}
	err := s.store.User().Create(c, tx, newUser)
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("create user error: %v", err)
		return err
	}

	if len(args.RoleIDs) == 0 {
		err = s.store.UserRole().Create(c, tx, &model.UserRole{
			UserID: newUser.ID,
			RoleID: model.UserRoleGuest,
		})
		if err != nil {
			tx.Rollback()
			log.C(c).Warnf("create user role error: %v", err)
			return err
		}
		tx.Commit()
		log.C(c).Infof("create user:%v success", newUser)
		return nil
	}

	for _, roleID := range args.RoleIDs {
		binds = append(binds, model.UserRole{
			UserID: newUser.ID,
			RoleID: roleID,
		})
	}

	err = s.store.UserRole().CreateInBatch(c, tx, binds)
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("create user role error: %v", err)
		return err
	}

	tx.Commit()

	log.C(c).Infof("create user:%v success", newUser)
	return nil
}

func (s *service) Login(c *gin.Context, args *v1.UserLoginArgs) (*v1.UserLoginReply, error) {
	var roleIDs []int64

	userInfo, err := s.store.User().GetBy(c, s.db, s.store.User().WithAccount(args.Account))
	if err != nil {
		log.C(c).Warnf("get user cnt error: %v", err)
		return nil, err
	}

	binds, err := s.store.UserRole().GetListBy(c, s.db, s.store.UserRole().WithUserIDs([]int64{userInfo.ID}))
	if err != nil {
		log.C(c).Warnf("get user role binds error: %v", err)
		return nil, err
	}

	for _, bind := range binds {
		roleIDs = append(roleIDs, bind.RoleID)
	}

	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(userInfo.Password, "$")
	correct := password.Verify(args.Password, passwordInfo[2], passwordInfo[3], options)

	if !correct {
		log.C(c).Warnf("password incorrect")
		return nil, errors.WithCode(code.ErrPasswordIncorrect, "password incorrect")
	}

	j := auth.NewJWT()

	token, err := j.CreateToken(&auth.JwtClaims{
		ID:      userInfo.ID,
		Account: userInfo.Account,
		RoleIDs: roleIDs,
	}, configs.LoginTokenExpire)
	if err != nil {
		log.C(c).Warnf("create token error: %v", err)
		return nil, errors.WithCode(code.ErrTokenInvalid, "create token error")
	}

	log.C(c).Infof("user login success, user:%v", userInfo)

	return &v1.UserLoginReply{
		TokenHead: "Bearer ",
		Token:     token,
		MetaID: v1.MetaID{
			ID: userInfo.ID,
		},
	}, nil
}

func (s *service) Delete(c *gin.Context, args *v1.UserDeleteArgs) error {

	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.store.User().DeleteBy(c, tx, s.store.WithID(args.ID))
		if err != nil {
			log.C(c).Warnf("delete user error: %v", err)
			return err
		}
		err = s.store.UserRole().DeleteBy(c, tx, s.store.UserRole().WithUserIDs([]int64{args.ID}))
		if err != nil {
			log.C(c).Warnf("delete user role error: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.C(c).Warnf("delete user error: %v", err)
		return err
	}

	log.C(c).Infof("delete user success, id:%v", args.ID)
	return nil
}

func (s *service) GetDetail(c *gin.Context, args *v1.UserDetailArgs) (*v1.UserDetailReply, error) {

	cnt, err := s.store.User().GetCntBy(c, s.db, s.store.WithID(args.ID))
	if err != nil {
		log.C(c).Warnf("get user cnt error: %v", err)
		return nil, err
	}
	if cnt == 0 {
		log.C(c).Warnf("user not existed")
		return nil, errors.WithCode(code.ErrDataNotFound, "data not existed")
	}
	userInfo, err := s.store.User().GetBy(c, s.db,
		s.store.WithID(args.ID),
		s.store.WithPreload("Roles"),
		s.store.WithPreload("Roles.Permissions"),
		s.store.WithPreload("System"),
	)
	if err != nil {
		log.C(c).Warnf("get user error: %v", err)
		return nil, err
	}

	userDetail := mapper.UserBoMapper(*userInfo)

	return &userDetail, nil

}

func (s *service) GetList(c *gin.Context, args *v1.UserListArgs) (*v1.UserListReply, error) {
	var usersDetail []v1.UserDetailReply

	queryArgs := queryArgs{
		Page:     args.Page,
		PageSize: args.PageSize,
	}

	query := s.listQuery(queryArgs)
	query = append(query, s.store.WithPreload("Roles"), s.store.WithPreload("Roles.Permissions"))

	users, err := s.store.User().GetListBy(c, s.db, query...)
	if err != nil {
		log.C(c).Warnf("get user list error: %v", err)
		return nil, err
	}

	for _, user := range users {
		userReply := mapper.UserBoMapper(user)
		usersDetail = append(usersDetail, userReply)
	}

	cnt, err := s.store.User().GetCntBy(c, s.db)
	if err != nil {
		log.C(c).Warnf("get user cnt error: %v", err)
		return nil, err
	}

	return &v1.UserListReply{
		MetaPage: v1.MetaPage{
			Page:     args.Page,
			PageSize: args.PageSize,
		},
		List: usersDetail,
		MetaTotalCnt: v1.MetaTotalCnt{
			TotalCnt: cnt,
		},
	}, nil
}

func (s *service) Update(c *gin.Context, args *v1.UserUpdateArgs) error {
	var binds []model.UserRole
	tx := s.db.Begin()
	log.C(c).Infof("update user info, args:%v", args)

	newUser := new(model.User)

	newUser.Name = args.Name
	newUser.Account = args.Account
	if args.Password != "" {
		newUser.Password = userutils.EncodePassword(args.Password)
	}
	err := s.store.User().Updates(c, tx, newUser, s.store.WithID(args.ID))
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("update user error: %v", err)

		return err
	}

	err = s.store.UserRole().DeleteBy(c, tx, s.store.UserRole().WithUserIDs([]int64{args.ID}))
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("delete user role error: %v", err)
		return err
	}

	if len(args.RoleIDs) == 0 {
		err = s.store.UserRole().Create(c, tx, &model.UserRole{

			UserID: args.ID,
			RoleID: model.UserRoleGuest,
		})
		if err != nil {
			tx.Rollback()
			log.C(c).Warnf("create user role error: %v", err)
			return err
		}
		tx.Commit()
		log.C(c).Infof("update user success, id:%v", args.ID)
		return nil
	}

	for _, roleID := range args.RoleIDs {
		binds = append(binds, model.UserRole{
			UserID: args.ID,
			RoleID: roleID,
		})
	}
	err = s.store.UserRole().CreateInBatch(c, tx, binds)
	if err != nil {
		tx.Rollback()
		log.C(c).Warnf("create user role error: %v", err)
		return err
	}

	tx.Commit()
	log.C(c).Infof("update user success, id:%v", args.ID)
	return nil
}

func (s *service) listQuery(args queryArgs) []store.DBOptions {
	var query []store.DBOptions

	if args.Page > 0 && args.PageSize > 0 {
		query = append(query, s.store.WithPage(args.Page, args.PageSize))
	}
	return query
}
