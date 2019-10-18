package nacos

type NamingSerivce struct {
	c *Client
}

func NewNamingService(c *Client) *NamingSerivce {
	return &NamingSerivce{c: c}
}

type Instance struct {
	Ip          string
	Port        int
	Namespace   string
	Weight      float64
	Enable      bool
	Healthy     bool
	Metadata    string
	ClusterName string
	ServiceName string
	GroupName   string
	Ephemeral   bool
}

func (ns *NamingSerivce) RegisterInstance(instance *Instance) error {

}


func (ns *NamingSerivce)DeregisterInstance(instance *Instance) error {

}