package lb

import (
	"github.com/chenqinghe/nacos-go-sdk/discovery"
	"sync/atomic"
)

type RoundRobin struct {
	index uint64
}


func (rr *RoundRobin)Select(instances []*discovery.Instance)*discovery.Instance {
	instance:=instances[rr.index%uint64(len(instances))]

	atomic.AddUint64(&rr.index,1)

	return instance
}