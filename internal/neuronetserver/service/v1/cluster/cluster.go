package cluster

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	v1 "neuronet/internal/neuronetserver/dto/v1"
	"neuronet/internal/neuronetserver/model"
	"neuronet/internal/neuronetserver/store"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	"neuronet/pkg/k8s"
	"neuronet/pkg/k8s/informer"
	"neuronet/pkg/log"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(c *gin.Context, args *v1.ClusterCreateArgs) (*v1.ClusterCreateReply, error)
	Delete(c *gin.Context, args *v1.ClusterDeleteArgs) error
	GetList(c *gin.Context, args *v1.ClusterListArgs) (*v1.ClusterListReply, error)
	Update(c *gin.Context, args *v1.ClusterUpdateArgs) error
}

type service struct {
	store store.Factory
	db    *gorm.DB
}

type queryArgs struct {
	Page     int
	PageSize int
	Name     string
}

func NewService(db *gorm.DB, store store.Factory) *service {
	return &service{
		store: store,
		db:    db,
	}
}

func (s *service) Create(c *gin.Context, args *v1.ClusterCreateArgs) (*v1.ClusterCreateReply, error) {
	cnt, err := s.store.Cluster().GetCntBy(c, s.db, s.store.WithName(args.Name))
	if err != nil {
		log.C(c).Warnf("get cluster cnt error: %v", err)
		return nil, err
	}
	if cnt != 0 {
		log.C(c).Warnf("cluster existed")
		return nil, errors.WithCode(code.ErrDataExisted, "data existed")
	}

	newCluster := &model.Cluster{
		Name:        args.Name,
		ConfigPath:  args.KubeConfigPath,
		Description: args.Description,
	}

	clusterSets := k8s.GetClusterSets()
	config, err := clientcmd.BuildConfigFromFlags("", newCluster.ConfigPath)
	if err != nil {
		log.Warnf("Can't create config %v", err)
		return nil, errors.WithCode(code.ErrInternalServer, "can't create config")
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Warnf("Can't create clientSet %v", err)
		return nil, errors.WithCode(code.ErrInternalServer, "can't create clientSet")
	}

	metricSet, err := versioned.NewForConfig(config)
	if err != nil {
		log.Warnf("Can't create metricSet %v", err)
	}

	stop := make(chan struct{})
	informerStore, err := informer.NewInformerStore(stop, clientSet)
	if err != nil {
		log.C(c).Warnf("create informer store error: %v", err)
		return nil, errors.WithCode(code.ErrInternalServer, "create informer store error")
	}

	err = s.store.Cluster().Create(c, s.db, newCluster)
	if err != nil {
		close(stop)
		log.C(c).Warnf("create cluster error: %v", err)
		return nil, err
	}

	clientSets := &k8s.ClientSets{
		K8sClient:      clientSet,
		MetricsClient:  metricSet,
		InformerClient: informerStore,
	}
	clusterSets.Add(newCluster.Name, clientSets)

	return &v1.ClusterCreateReply{
		MetaID: v1.MetaID{
			ID: newCluster.ID,
		},
	}, nil
}

func (s *service) Delete(c *gin.Context, args *v1.ClusterDeleteArgs) error {
	clusterBo, err := s.store.Cluster().GetBy(c, s.db, s.store.WithID(args.ID))
	if err != nil {
		log.C(c).Warnf("get cluster error: %v", err)
		return err
	}
	err = s.store.Cluster().DeleteBy(c, s.db, s.store.WithID(args.ID))
	if err != nil {
		log.Warnf("delete cluster error: %v", err)
		return err
	}
	clusterSets := k8s.GetClusterSets()
	clusterSets.Delete(clusterBo.Name)
	return nil
}

func (s *service) GetList(c *gin.Context, args *v1.ClusterListArgs) (*v1.ClusterListReply, error) {
	var clustersDetailReply []v1.ClusterDetailReply
	listQueryArgs := queryArgs{
		Page:     args.Page,
		PageSize: args.PageSize,
		Name:     args.Name,
	}
	listQuery := s.listQuery(listQueryArgs)
	clusterBos, err := s.store.Cluster().GetListBy(c, s.db, listQuery...)
	if err != nil {
		log.C(c).Warnf("get cluster list error: %v", err)
		return nil, err
	}

	for _, bo := range clusterBos {
		clustersDetailReply = append(clustersDetailReply, v1.ClusterDetailReply{
			MetaID: v1.MetaID{
				ID: bo.ID,
			},
			MetaName: v1.MetaName{
				Name: bo.Name,
			},
			MetaTime: v1.MetaTime{
				CreateTime: bo.CreatedAt,
				UpdateTime: bo.UpdatedAt,
			},
			KubeConfigPath: bo.ConfigPath,
		})
	}

	cntQueryArgs := queryArgs{
		Name: args.Name,
	}
	cntQuery := s.listQuery(cntQueryArgs)
	cnt, err := s.store.Cluster().GetCntBy(c, s.db, cntQuery...)
	if err != nil {
		log.C(c).Warnf("get cluster cnt error: %v", err)
		return nil, err
	}

	return &v1.ClusterListReply{
		MetaPage: v1.MetaPage{
			Page:     args.Page,
			PageSize: args.PageSize,
		},
		MetaTotalCnt: v1.MetaTotalCnt{
			TotalCnt: cnt,
		},
		List: clustersDetailReply,
	}, nil
}

func (s *service) Update(c *gin.Context, args *v1.ClusterUpdateArgs) error {
	var clientSets *k8s.ClientSets

	clusterBo, err := s.store.Cluster().GetBy(c, s.db, s.store.WithID(args.ID))
	if err != nil {
		log.Warnf("get cluster error: %v", err)
		return err
	}

	newCluster := &model.Cluster{
		Name:        args.Name,
		ConfigPath:  args.KubeConfigPath,
		Description: args.Description,
	}

	stop := make(chan struct{})
	clusterSets := k8s.GetClusterSets()

	if args.KubeConfigPath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", newCluster.ConfigPath)
		if err != nil {
			log.Warnf("Can't create config %v", err)
			return errors.WithCode(code.ErrInternalServer, "can't create config")
		}

		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Warnf("Can't create clientSet %v", err)
			return errors.WithCode(code.ErrInternalServer, "can't create clientSet")
		}

		metricSet, err := versioned.NewForConfig(config)
		if err != nil {
			log.Warnf("Can't create metricSet %v", err)
		}

		stop := make(chan struct{})
		informerStore, err := informer.NewInformerStore(stop, clientSet)
		if err != nil {
			log.C(c).Warnf("create informer store error: %v", err)
			return errors.WithCode(code.ErrInternalServer, "create informer store error")
		}

		clientSets = &k8s.ClientSets{
			K8sClient:      clientSet,
			MetricsClient:  metricSet,
			InformerClient: informerStore,
		}
	} else {
		clientSets = clusterSets.Get(clusterBo.Name)
	}

	err = s.store.Cluster().Updates(c, s.db, newCluster, s.store.WithID(args.ID))
	if err != nil {
		close(stop)
		log.C(c).Warnf("update cluster error: %v", err)
		return err
	}

	clusterSets.Update(newCluster.Name, clientSets)

	return nil
}

func (s *service) listQuery(args queryArgs) []store.DBOptions {
	var query []store.DBOptions

	if args.Page > 0 && args.PageSize > 0 {
		query = append(query, s.store.WithPage(args.Page, args.PageSize))
	}
	if args.Name != "" {
		query = append(query, s.store.WithNameLike(args.Name))
	}
	return query
}
