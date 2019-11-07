package naming

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chenqinghe/nacos-go-sdk/api/v1"
)

type Client struct {
	c *v1.Client
}

func NewNamingService(c *v1.Client) *Client {
	return &Client{c: c}
}

type Instance interface {
	GetId() string
	GetIp() string
	GetPort() int
	GetNamespace() string
	GetWeight() float64
	GetEnable() bool
	GetHealthy() bool
	GetMetadata() Metadata
	GetClusterName() string
	GetServiceName() string
	GetGroupName() string
	GetEphemeral() bool
}

type Metadata map[string]interface{}

func (m Metadata) String() string {
	d, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(d)
}

func (ns *Client) RegisterInstance(instance Instance) error {
	values := make(url.Values)
	values.Set("ip", instance.GetIp())
	values.Set("port", strconv.Itoa(instance.GetPort()))
	values.Set("namespaceId", instance.GetNamespace())
	values.Set("weight", strconv.FormatFloat(instance.GetWeight(), 'f', 2, 64))
	values.Set("enabled", strconv.FormatBool(instance.GetEnable()))
	values.Set("healthy", strconv.FormatBool(instance.GetHealthy()))
	values.Set("metadata", instance.GetMetadata().String())
	values.Set("clusterName", instance.GetClusterName())
	values.Set("serviceName", instance.GetServiceName())
	values.Set("groupName", instance.GetGroupName())
	values.Set("ephemeral", strconv.FormatBool(instance.GetEphemeral()))

	req, err := http.NewRequest(http.MethodPost, v1.JoinUrlQueryString(ns.c.GetUrl(v1.InstancePath), values), nil)
	if err != nil {
		return err
	}
	resp, err := ns.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code not ok: %d, body: %s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if r := string(data); r != "ok" {
		return v1.ErrUnexpectedResponse{Expected: "ok", Data: r}
	}

	return nil
}

func (ns *Client) DeregisterInstance(instance Instance) error {
	values := make(url.Values)
	values.Set("ip", instance.GetIp())
	values.Set("port", strconv.Itoa(instance.GetPort()))
	values.Set("namespaceId", instance.GetNamespace())
	values.Set("clusterName", instance.GetClusterName())
	values.Set("serviceName", instance.GetServiceName())
	values.Set("groupName", instance.GetGroupName())
	values.Set("ephemeral", strconv.FormatBool(instance.GetEphemeral()))

	req, err := http.NewRequest(http.MethodDelete, v1.JoinUrlQueryString(ns.c.GetUrl(v1.InstancePath), values), nil)
	if err != nil {
		return err
	}

	resp, err := ns.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code not ok: %d, body:%s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if r := string(data); r != "ok" {
		return v1.ErrUnexpectedResponse{Expected: "ok", Data: r}
	}

	return nil
}

func (ns *Client) UpdateInstance(instance Instance) error {
	values := make(url.Values)
	values.Set("ip", instance.GetIp())
	values.Set("port", strconv.Itoa(instance.GetPort()))
	values.Set("namespaceId", instance.GetNamespace())
	values.Set("weight", strconv.FormatFloat(instance.GetWeight(), 'f', 2, 64))
	values.Set("enabled", strconv.FormatBool(instance.GetEnable()))
	values.Set("healthy", strconv.FormatBool(instance.GetHealthy()))
	values.Set("metadata", instance.GetMetadata().String())
	values.Set("clusterName", instance.GetClusterName())
	values.Set("serviceName", instance.GetServiceName())
	values.Set("groupName", instance.GetGroupName())
	values.Set("ephemeral", strconv.FormatBool(instance.GetEphemeral()))

	req, err := http.NewRequest(http.MethodPut, v1.JoinUrlQueryString(ns.c.GetUrl(v1.InstancePath), values), nil)
	if err != nil {
		return err
	}
	resp, err := ns.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code not ok: %d, body: %s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if r := string(data); r != "ok" {
		return v1.ErrUnexpectedResponse{Expected: "ok", Data: r}
	}

	return nil
}

type GetInstanceOption struct {
	GroupName   string
	NamespaceId string
	Clusters    []string
	HealthyOnly bool
	Ephemeral   bool
}

type node struct {
	Valid      bool
	Marked     bool
	InstanceId string
	Port       int
	Ip         string
	Weight     float64
	Metadata   Metadata

	clusterName string
	serviceName string
	groupName   string
	namespaceId string
}

func (i *node) GetId() string          { return i.InstanceId }
func (i *node) GetIp() string          { return i.Ip }
func (i *node) GetPort() int           { return i.Port }
func (i *node) GetNamespace() string   { return i.namespaceId }
func (i *node) GetWeight() float64     { return i.Weight }
func (i *node) GetEnable() bool        { return i.Valid }
func (i *node) GetHealthy() bool       { return i.Valid }
func (i *node) GetMetadata() Metadata  { return i.Metadata }
func (i *node) GetClusterName() string { return i.clusterName }
func (i *node) GetServiceName() string { return i.serviceName }
func (i *node) GetGroupName() string   { return i.groupName }
func (i *node) GetEphemeral() bool     { return true }

func (ns *Client) GetInstances(serviceName string, option *GetInstanceOption) ([]Instance, error) {
	values := make(url.Values)
	values.Set("serviceName", serviceName)

	if option != nil {
		if option.GroupName != "" {
			values.Set("groupName", option.GroupName)
		}
		if option.NamespaceId != "" {
			values.Set("namespaceId", option.NamespaceId)
		}
		if len(option.Clusters) > 0 {
			values.Set("clusters", strings.Join(option.Clusters, ","))
		}
		if option.HealthyOnly {
			values.Set("healthyOnly", "true")
		}
	}

	resp, err := ns.c.Get(v1.JoinUrlQueryString(ns.c.GetUrl(v1.InstanceListPath), values))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response code not ok: %d, body:%s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Hosts []*node `json:"hosts"`
	}
	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	instances := make([]Instance, len(r.Hosts))
	for k := range r.Hosts {
		instances[k] = r.Hosts[k]
	}

	return instances, nil
}

func (ns *Client) GetInstance(serviceName string, ip, port string, option *GetInstanceOption) (Instance, error) {
	values := make(url.Values)
	values.Set("serviceName", serviceName)
	values.Set("ip", ip)
	values.Set("port", port)

	if option != nil {
		if option.GroupName != "" {
			values.Set("groupName", option.GroupName)
		}
		if option.NamespaceId != "" {
			values.Set("namespaceId", option.NamespaceId)
		}
		if len(option.Clusters) > 0 {
			values.Set("clusters", strings.Join(option.Clusters, ","))
		}
		if option.HealthyOnly {
			values.Set("healthyOnly", "true")
		}
		if option.Ephemeral {
			values.Set("ephemeral", "true")
		}
	}

	resp, err := ns.c.Get(v1.JoinUrlQueryString(ns.c.GetUrl(v1.InstancePath), values))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response code not ok: %d, body:%s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var nd node
	if err := json.Unmarshal(data, &nd); err != nil {
		return nil, err
	}

	return &nd, nil
}

func (ns *Client) Heartbeat(instance Instance) (time.Duration, error) {
	values := make(url.Values)
	values.Set("serviceName", instance.GetServiceName())
	values.Set("groupName", instance.GetGroupName())
	values.Set("ephemeral", strconv.FormatBool(instance.GetEphemeral()))

	beat, err := json.Marshal(instance)
	if err != nil {
		return 0, err
	}
	values.Set("beat", string(beat))

	req, err := http.NewRequest(http.MethodPut, v1.JoinUrlQueryString(ns.c.GetUrl(v1.InstanceHeartbeatPath), values), nil)
	if err != nil {
		return 0, err
	}

	resp, err := ns.c.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("http response code not ok: %d, body: %s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	type Response struct {
		ClientBeatInterval int `json:"clientBeatInterval"`
	}

	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		return 0, err
	}

	return time.Millisecond * time.Duration(r.ClientBeatInterval), nil

}

type Service struct {
	ServiceName      string    `json:"name"`
	GroupName        string    `json:"groupName"`        // 字符串	否	分组名
	NamespaceId      string    `json:"namespaceId"`      // 字符串	否	命名空间ID
	ProtectThreshold float64   `json:"protectThreshold"` // 浮点数	否	保护阈值,取值0到1,默认0
	Metadata         Metadata  `json:"metadata"`         // 字符串	否	元数据
	Selector         Metadata  `json:"selector"`         // JSON格式字符串	否	访问策略
	Clusters         []Cluster `json:"clusters"`
}

type Cluster struct {
	Name          string        `json:"name"`
	HealthChecker HealthChecker `json:"healthChecker"`
	Metadata      Metadata      `json:"metadata"`
}

type HealthChecker struct {
	Type string
}

func (ns *Client) CreateService(service *Service) error {
	values := make(url.Values)
	values.Set("serviceName", service.ServiceName)

	if service.GroupName != "" {
		values.Set("groupName", service.GroupName)
	}
	if service.NamespaceId != "" {
		values.Set("namespaceId", service.NamespaceId)
	}
	if service.ProtectThreshold != 0 {
		values.Set("protectThreshold", strconv.FormatFloat(service.ProtectThreshold, 'f', 2, 64))
	}
	if len(service.Metadata) > 0 {
		values.Set("metadata", service.Metadata.String())
	}
	if len(service.Selector) > 0 {
		values.Set("selector", service.Selector.String())
	}

	req, err := http.NewRequest(http.MethodPost, v1.JoinUrlQueryString(ns.c.GetUrl(v1.ServicePath), values), nil)
	if err != nil {
		return err
	}

	resp, err := ns.c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code not ok: %d, body: %s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if r := string(data); r != "ok" {
		return v1.ErrUnexpectedResponse{Expected: "ok", Data: r}
	}

	return nil
}

func (ns *Client) DeleteService(service *Service) error {
	values := make(url.Values)
	values.Set("serviceName", service.ServiceName)

	if service.GroupName != "" {
		values.Set("groupName", service.GroupName)
	}
	if service.NamespaceId != "" {
		values.Set("namespaceId", service.NamespaceId)
	}

	req, err := http.NewRequest(http.MethodDelete, v1.JoinUrlQueryString(ns.c.GetUrl(v1.ServicePath), values), nil)
	if err != nil {
		return err
	}

	resp, err := ns.c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code not ok: %d, body: %s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if r := string(data); r != "ok" {
		return v1.ErrUnexpectedResponse{Expected: "ok", Data: r}
	}

	return nil
}

func (ns *Client) UpdateService(service *Service) error {
	values := make(url.Values)
	values.Set("serviceName", service.ServiceName)

	if service.GroupName != "" {
		values.Set("groupName", service.GroupName)
	}
	if service.NamespaceId != "" {
		values.Set("namespaceId", service.NamespaceId)
	}
	if service.ProtectThreshold != 0 {
		values.Set("protectThreshold", strconv.FormatFloat(service.ProtectThreshold, 'f', 2, 64))
	}
	if len(service.Metadata) > 0 {
		values.Set("metadata", service.Metadata.String())
	}
	if len(service.Selector) > 0 {
		values.Set("selector", service.Selector.String())
	}

	req, err := http.NewRequest(http.MethodPut, v1.JoinUrlQueryString(ns.c.GetUrl(v1.ServicePath), values), nil)
	if err != nil {
		return err
	}

	resp, err := ns.c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code not ok: %d, body: %s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if r := string(data); r != "ok" {
		return v1.ErrUnexpectedResponse{Expected: "ok", Data: r}
	}

	return nil
}

func (ns *Client) QueryService(serviceName, groupName, namespace string) (*Service, error) {
	values := make(url.Values)
	values.Set("serviceName", serviceName)
	if groupName != "" {
		values.Set("groupName", groupName)
	}
	if namespace != "" {
		values.Set("namespaceId", namespace)
	}

	req, err := http.NewRequest(http.MethodGet, v1.JoinUrlQueryString(ns.c.GetUrl(v1.ServicePath), values), nil)
	if err != nil {
		return nil, err
	}

	resp, err := ns.c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response code not ok: %d, body: %s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var srv Service
	if err := json.Unmarshal(data, &srv); err != nil {
		return nil, err
	}

	return &srv, nil

}

func (ns *Client) ListService(pageNo, pageSize int, namespace, groupName string) ([]string, error) {
	values := make(url.Values)
	values.Set("pageNo", strconv.Itoa(pageNo))
	values.Set("pageSize", strconv.Itoa(pageSize))
	if groupName != "" {
		values.Set("groupName", groupName)
	}
	if namespace != "" {
		values.Set("namespaceId", namespace)
	}

	resp, err := ns.c.Get(v1.JoinUrlQueryString(ns.c.GetUrl(v1.ServiceListPath), values))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response code not ok: %d, body: %s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("services:", string(data))

	type Response struct {
		Count int      `json:"count"`
		Doms  []string `json:"doms"`
	}

	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}

	return r.Doms, nil
}
