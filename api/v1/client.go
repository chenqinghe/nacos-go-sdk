package v1

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	baseUrl string
	*http.Client
}

func NewNacosClient(baseUrl string) *Client {
	if baseUrl[len(baseUrl)-1] == '/' {
		baseUrl = string(baseUrl[:len(baseUrl)-1])
	}
	return &Client{
		baseUrl: baseUrl,
		Client:  &http.Client{},
	}
}

type PathType int

const (
	pathStart PathType = iota
	ConfigPath
	ConfigListenerPath
	InstancePath
	InstanceListPath
	InstanceHeartbeatPath
	ServicePath
	ServiceListPath
	SystemSwitchesPath
	pathEnd
)

var paths = [...]string{
	ConfigPath:            "/nacos/v1/cs/configs",
	ConfigListenerPath:    "/nacos/v1/cs/configs/listener",
	InstancePath:          "/nacos/v1/ns/instance",
	InstanceListPath:      "/nacos/v1/ns/instance/list",
	InstanceHeartbeatPath: "/nacos/v1/ns/instance/beat",
	ServicePath:           "/nacos/v1/ns/service",
	ServiceListPath:       "/nacos/v1/ns/service/list",
	SystemSwitchesPath:    "/nacos/v1/ns/operator/switches",
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

func ReadResponseBody(body io.Reader) string {
	if body == nil {
		return ""
	}

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(data)
}

type ErrUnexpectedResponse struct {
	Expected string
	Data     string
}

func (e ErrUnexpectedResponse) Error() string {
	return fmt.Sprintf("expect response body: %s, but got: %s", e.Expected, e.Data)
}

func JoinUrlQueryString(url string, values url.Values) string {
	return url + "?" + values.Encode()
}
