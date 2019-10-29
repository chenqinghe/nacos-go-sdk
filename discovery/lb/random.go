package lb

import (
	"github.com/chenqinghe/nacos-go-sdk/discovery"
	"math/rand"
	"time"
)

type Random struct {
	r *rand.Rand
}

func NewRandom(seed ...int64) *Random {
	if len(seed) == 0 {
		return &Random{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
	}
	return &Random{r: rand.New(rand.NewSource(seed[0]))}
}

func (r *Random) Select(instances []*discovery.Instance) *discovery.Instance {
	return instances[r.r.Intn(len(instances))]
}
