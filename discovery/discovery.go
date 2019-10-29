package discovery

import (
	"github.com/chenqinghe/nacos-go-sdk/api/v1"
	"github.com/chenqinghe/nacos-go-sdk/api/v1/naming"
	"github.com/chenqinghe/nacos-go-sdk/discovery/lb"
)

type Discovery interface {

	// RegisterInstance 注册服务实例
	RegisterInstance(instance *Instance) error

	// UpdateInstance 更新实例信息
	UpdateInstance(instance *Instance) error

	// DeregisterInstance 注销服务实例
	DeregisterInstance(instance *Instance) error

	// QueryServices 查询服务列表
	QueryServices() ([]string, error)

	// QueryInstances 查询服务实例列表
	QueryInstances(serviceName string) ([]*Instance, error)

	// GetInstance 获取一个服务实例，可通过一定的负载均衡策略
	GetInstance(serviceName string) (*Instance, error)
}

type Instance = naming.Instance

type Metadata = naming.Metadata

type nacosDiscovery struct {
	namingService *naming.Client
	snapshot      Snapshoter
	lbStrategy    lb.Strategy
	logger        Logger
}

var _ Discovery = (*nacosDiscovery)(nil)

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func NewNacosDiscovery(c *v1.Client, options ...Option) *nacosDiscovery {
	nd := &nacosDiscovery{
		namingService: naming.NewNamingService(c),
		snapshot:      nil,
		lbStrategy:    lb.NewRandom(),
		logger:        nil,
	}

	for _, opt := range options {
		opt(nd)
	}

	return nd
}

type Option func(discovery *nacosDiscovery)

func SetSnapshot(snapshoter Snapshoter) Option {
	return func(discovery *nacosDiscovery) {
		discovery.snapshot = snapshoter
	}
}

func SetLBStrategy(strategy lb.Strategy) Option {
	return func(discovery *nacosDiscovery) {
		discovery.lbStrategy = strategy
	}
}

func SetLogger(logger Logger) Option {
	return func(discovery *nacosDiscovery) {
		discovery.logger = logger
	}
}

func (d *nacosDiscovery) RegisterInstance(instance *Instance) error {
	err := d.namingService.RegisterInstance(instance)
	if err != nil {
		return err
	}

	go func() {
		for {
			_, err := d.namingService.Heartbeat(instance)
			if err != nil {
				if d.logger != nil {
					d.logger.Errorf("send heartbeat error:%s\n", err)
				}
				return
			}
			// TODO: stop when deregister
		}
	}()

	return nil
}

func (d *nacosDiscovery) DeregisterInstance(instance *Instance) error {
	return d.namingService.DeregisterInstance(instance)
}

func (d *nacosDiscovery) UpdateInstance(instance *Instance) error {
	return d.namingService.UpdateInstance(instance)
}

func (d *nacosDiscovery) QueryInstances(serviceName string) ([]*Instance, error) {
	return d.namingService.GetInstances(serviceName, nil)
}

func (d *nacosDiscovery) GetInstance(serviceName string) (*Instance, error) {
	instances, err := d.namingService.GetInstances(serviceName, nil)
	if err != nil {
		return nil, err
	}
	return d.lbStrategy.Select(instances), nil
}

func (d *nacosDiscovery) QueryServices() ([]string, error) {
	var maxServices = 9999
	return d.namingService.ListService(1, maxServices, "", "")
}
