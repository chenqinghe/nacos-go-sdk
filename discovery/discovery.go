package discovery

import (
	v1 "github.com/chenqinghe/nacos-go-sdk/api/v1"
	"github.com/chenqinghe/nacos-go-sdk/api/v1/naming"
	"github.com/chenqinghe/nacos-go-sdk/discovery/lb"
	"github.com/rfyiamcool/go-timewheel"
	"time"
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

type Instance struct {
	Id          string   `json:"instanceId"`
	Ip          string   `json:"ip"`
	Port        int      `json:"port"`
	Namespace   string   `json:"namespace"`
	Weight      float64  `json:"weight"`
	Enable      bool     `json:"enable"`
	Healthy     bool     `json:"healthy"`
	Metadata    Metadata `json:"metadata"`
	ClusterName string   `json:"clusterName"`
	ServiceName string   `json:"serviceName"`
	GroupName   string   `json:"groupName"`
	Ephemeral   bool     `json:"ephemeral"`
}

func (i *Instance) GetId() string          { return i.Id }
func (i *Instance) GetIp() string          { return i.Ip }
func (i *Instance) GetPort() int           { return i.Port }
func (i *Instance) GetNamespace() string   { return i.Namespace }
func (i *Instance) GetWeight() float64     { return i.Weight }
func (i *Instance) GetEnable() bool        { return i.Enable }
func (i *Instance) GetHealthy() bool       { return i.Healthy }
func (i *Instance) GetMetadata() Metadata  { return i.Metadata }
func (i *Instance) GetClusterName() string { return i.ClusterName }
func (i *Instance) GetServiceName() string { return i.ServiceName }
func (i *Instance) GetGroupName() string   { return i.GroupName }
func (i *Instance) GetEphemeral() bool     { return i.Ephemeral }

type Metadata = naming.Metadata

type nacosDiscovery struct {
	namingService *naming.Client
	lbStrategy    lb.Strategy
	logger        Logger

	tw *timewheel.TimeWheel

	// TODO: concurrent access protect
	registeredInstances map[string]*Instance
	tasks               map[string]*timewheel.Task
}

var _ Discovery = (*nacosDiscovery)(nil)

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func NewNacosDiscovery(c *v1.Client, options ...Option) *nacosDiscovery {
	tw, _ := timewheel.NewTimeWheel(time.Second, 3600)
	tw.Start()

	nd := &nacosDiscovery{
		namingService: naming.NewNamingService(c),
		lbStrategy:    lb.NewRandom(),
		logger:        nil,
		tw:            tw,
	}

	for _, opt := range options {
		opt(nd)
	}

	return nd
}

type Option func(discovery *nacosDiscovery)

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
	d.registeredInstances[instance.ServiceName] = instance
	err := d.namingService.RegisterInstance(instance)
	if err != nil {
		return err
	}

	// TODO: 根据服务端返回的时间间隔发送心跳
	task := d.tw.Add(time.Second*5, func() {
		_, err := d.namingService.Heartbeat(instance)
		if err != nil {
			if d.logger != nil {
				d.logger.Errorf("send heartbeat error:%s\n", err)
			}
		}
	})

	d.registeredInstances[instance.GetId()] = instance
	d.tasks[instance.GetId()] = task

	return nil
}

func (d *nacosDiscovery) DeregisterInstance(instance *Instance) error {
	key := instance.GetId()
	task := d.tasks[key]
	delete(d.tasks, key)
	delete(d.registeredInstances, key)

	d.tw.Remove(task)

	return d.namingService.DeregisterInstance(instance)
}

func (d *nacosDiscovery) UpdateInstance(instance *Instance) error {
	return d.namingService.UpdateInstance(instance)
}

func (d *nacosDiscovery) QueryInstances(serviceName string) ([]*Instance, error) {
	is,err:= d.namingService.GetInstances(serviceName,nil)
	if err!=nil {
		return  nil,err
	}

	instances:=make([]*Instance,len(is))
	for k,v:=range is {
		instances[k] = &Instance{
			Id:          v.GetId(),
			Ip:          v.GetIp(),
			Port:        v.GetPort(),
			Namespace:   v.GetNamespace(),
			Weight:      v.GetWeight(),
			Enable:      v.GetEnable(),
			Healthy:     v.GetHealthy(),
			Metadata:    v.GetMetadata(),
			ClusterName: v.GetClusterName(),
			ServiceName: v.GetServiceName(),
			GroupName:   v.GetGroupName(),
			Ephemeral:   v.GetEphemeral(),
		}
	}

	return instances,nil
}

func (d *nacosDiscovery) GetInstance(serviceName string) (*Instance, error) {
	instances, err := d.namingService.GetInstances(serviceName, nil)
	if err != nil {
		return nil, err
	}

	instance := d.lbStrategy.Select(instances).(naming.Instance)

	return &Instance{
		Id:          instance.GetId(),
		Ip:          instance.GetIp(),
		Port:        instance.GetPort(),
		Namespace:   instance.GetNamespace(),
		Weight:      instance.GetWeight(),
		Enable:      instance.GetEnable(),
		Healthy:     instance.GetHealthy(),
		Metadata:    instance.GetMetadata(),
		ClusterName: instance.GetClusterName(),
		ServiceName: instance.GetServiceName(),
		GroupName:   instance.GetGroupName(),
		Ephemeral:   instance.GetEphemeral(),
	}, nil
}

func (d *nacosDiscovery) QueryServices() ([]string, error) {
	var maxServices = 9999
	return d.namingService.ListService(1, maxServices, "", "")
}
