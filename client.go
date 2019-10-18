package nacos

import "net/http"

type Client struct {
	baseUrl string
	c       *http.Client
}

func NewNacosClient(baseUrl string) *Client {
	if baseUrl[len(baseUrl)-1] == '/' {
		baseUrl = string(baseUrl[:len(baseUrl)-1])
	}
	return &Client{
		baseUrl: baseUrl,
		c:       &http.Client{},
	}
}

type PathType int

const (
	pathStart PathType = iota
	config
	configListener
	pathEnd
)

var paths = [...]string{
	config:         "/nacos/v1/cs/configs",
	configListener: "/nacos/v1/cs/configs/listener",
}

var pathMap map[PathType]string

func init() {
	pathMap = make(map[PathType]string)
	for i := pathStart; i < pathEnd; i++ {
		pathMap[i] = paths[i]
	}
}

func (c *Client) GetUrl(typ PathType) string {
	return c.baseUrl + pathMap[typ]
}
