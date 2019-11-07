package lb

import (
	"math/rand"
	"reflect"
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

func (r *Random) Select(instances interface{}) interface{} {
	v := reflect.ValueOf(instances)

	return v.Index(r.r.Intn(v.Len())).Interface()
}
