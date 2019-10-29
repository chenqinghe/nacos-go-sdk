package naming

import (
	"github.com/chenqinghe/nacos-go-sdk/api/v1"
	"testing"
)

func TestSystemService_QuerySwitches(t *testing.T) {
	c := v1.NewNacosClient("http://127.0.0.1:8848")

	srv := NewSystemService(c)

	if sc, err := srv.QuerySwitches(); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%#v", *sc)
	}
}

func TestSerivce_DeregisterInstance(t *testing.T) {
	c := v1.NewNacosClient("http://192.168.107.236:8848")

	srv := NewNamingService(c)

	err := srv.DeregisterInstance(&Instance{
		ServiceName: "adserver",
		GroupName:   "adserver",
		Ip:          "127.0.0.2",
		Port:        8899,
	})
	if err != nil {
		t.Fatal(err)
	}
}
