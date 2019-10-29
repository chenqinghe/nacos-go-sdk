package discovery

import (
	v1 "github.com/chenqinghe/nacos-go-sdk/api/v1"
	"testing"
	"time"
)

func TestNacosDiscovery_UpdateInstancet(t *testing.T) {
	srv := NewNacosDiscovery(v1.NewNacosClient("http://192.168.107.236:8848"))

	instance := &Instance{
		ServiceName: "appdemo",
		Ip:          "127.0.0.1",
		Port:        8848,
		Enable:      true,
		Metadata: Metadata{
			"name": "111",
			"age":  123,
		},
		Ephemeral: true,
	}

	if err := srv.RegisterInstance(instance); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 10)

	instance.Metadata = Metadata{
		"aaa": "bbb",
	}
	if err := srv.UpdateInstance(instance); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 10)

}
