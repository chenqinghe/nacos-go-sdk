package lb

import (
	"reflect"
	"sync/atomic"
)

type RoundRobin struct {
	index uint64
}

func (rr *RoundRobin) Select(instances interface{}) interface{} {
	v := reflect.ValueOf(instances)
	instance := v.Index(int(rr.index % uint64(v.Len()))).Interface()

	atomic.AddUint64(&rr.index, 1)

	return instance
}
