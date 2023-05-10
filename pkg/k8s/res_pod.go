package bykubernetes

import (
	"context"
	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"net/http"
	"neuronet/internal/pkg/code"
	"neuronet/pkg/errors"
	apiv1 "neuronet/pkg/k8s/api/v1"
	"neuronet/pkg/k8s/informer"
	"neuronet/pkg/k8s/meta"
	"neuronet/pkg/k8s/ws"
)

var _ IPodAction = (*pods)(nil)

type IPod interface {
	Pods() IPodAction
}

type IPodAction interface {
	Get(ctx context.Context, options meta.GetOptions) (*v1.Pod, error)
	List(ctx context.Context, options meta.ListOptions) ([]*v1.Pod, error)
	Top(ctx context.Context, options meta.TopOptions) (*apiv1.TopPod, error)
	Logs(ctx context.Context, w http.ResponseWriter, r *http.Request, options meta.LogOptions) error
	Command(ctx context.Context, w http.ResponseWriter, r *http.Request, options *meta.TerminalOptions, kubeconfigPath string) error
}

type pods struct {
	client        kubernetes.Interface
	metricsClient versioned.Interface
	informer      informer.Storer
}

func newPods(c kubernetes.Interface, metricsClient versioned.Interface, informerStore informer.Storer) *pods {
	return &pods{
		client:        c,
		metricsClient: metricsClient,
		informer:      informerStore,
	}
}

func (p *pods) List(ctx context.Context, options meta.ListOptions) ([]*v1.Pod, error) {
	var (
		list []*v1.Pod
		err  error
	)

	if options.Label != "" {
		list, err = p.informer.InformerPods().ListByLabel(ctx, options.Namespace, options.Label)
		if err != nil {
			return nil, err
		}
		return list, nil

	}

	list, err = p.informer.InformerPods().ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (p *pods) Top(ctx context.Context, options meta.TopOptions) (*apiv1.TopPod, error) {
	var (
		memory, cpu int64
	)

	podTop := new(apiv1.TopPod)
	metrics, err := p.metricsClient.MetricsV1beta1().PodMetricses(options.Namespace).Get(ctx, options.ObjectName, metav1.GetOptions{})
	if err != nil {
		if k8serror.IsNotFound(err) {
			return nil, errors.WithCode(code.ErrDataNotFound, err.Error())
		}
		return nil, err
	}

	for _, c := range metrics.Containers {
		memory += c.Usage.Memory().Value()
		cpu += c.Usage.Cpu().Value()
	}

	podTop.PodName = options.ObjectName
	podTop.Memory = memory
	podTop.Cpu = cpu

	return podTop, nil
}

func (p *pods) Get(ctx context.Context, options meta.GetOptions) (*v1.Pod, error) {
	return p.informer.InformerPods().Get(ctx, options)
}

func (p *pods) Logs(ctx context.Context, w http.ResponseWriter, r *http.Request, options meta.LogOptions) error {

	opts := &v1.PodLogOptions{
		Follow:    true,
		Container: options.ContainerName,
	}

	request := p.client.CoreV1().Pods(options.Namespace).GetLogs(options.ObjectName, opts)
	readCloser, err := request.Stream(ctx)
	if err != nil {
		return err
	}

	session, err := ws.NewPodLogSession(w, r)
	if err != nil {
		return err
	}

	defer func() {
		session.Close()
	}()

	session.SetReaderCloser(readCloser)

	go session.Read()
	go session.Write()

	select {
	case <-session.Done():
		return nil
	}
}

func (p *pods) Command(ctx context.Context, w http.ResponseWriter, r *http.Request, webShellOptions *meta.TerminalOptions, kubeconfigPath string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	session, err := ws.NewTerminalSession(w, r)
	if err != nil {
		return err
	}
	// 处理关闭
	defer func() {
		_ = session.Close()
	}()

	// 组装 POST 请求
	req := p.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(webShellOptions.ObjectName).
		Namespace(webShellOptions.Namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: webShellOptions.Container,
			Command:   []string{"/bin/bash"},
			Stderr:    true,
			Stdin:     true,
			Stdout:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	// remotecommand 主要实现了http 转 SPDY 添加X-Stream-Protocol-Version相关header 并发送请求
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	// 与 kubelet 建立 stream 连接
	if err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout:            session,
		Stdin:             session,
		Stderr:            session,
		TerminalSizeQueue: session,
		Tty:               true,
	}); err != nil {
		_, _ = session.Write([]byte("exec pod command failed," + err.Error()))
		// 标记关闭terminal
		session.Done()
	}

	return nil
}
