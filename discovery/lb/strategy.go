package lb

import "github.com/chenqinghe/nacos-go-sdk/discovery"

type Strategy interface {
	Select(instances []*discovery.Instance) *discovery.Instance
}
