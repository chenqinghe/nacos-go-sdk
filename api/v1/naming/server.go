package naming

import (
	"encoding/json"
	"fmt"
	"github.com/chenqinghe/nacos-go-sdk/api/v1"
	"io/ioutil"
	"net/http"
)

type SystemService struct {
	c *v1.Client
}

func NewSystemService(c *v1.Client) *SystemService {
	return &SystemService{c: c}
}

type SystemConfig struct {
	Name                                     string                 `json:"name"`
	Masters                                  interface{}            `json:"masters"`
	AdWeightMap                              map[string]interface{} `json:"adWeightMap"`
	DefaultPushCacheMillis                   int                    `json:"defaultPushCacheMillis"`
	ClientBeatInterval                       int                    `json:"clientBeatInterval"`
	DefaultCacheMillis                       int                    `json:"defaultCacheMillis"`
	DistroThreshold                          float64                `json:"distroThreshold"`
	HealthCheckEnabled                       bool                   `json:"healthCheckEnabled"`
	DistroEnabled                            bool                   `json:"distroEnabled"`
	EnableStandalone                         bool                   `json:"enableStandalone"`
	PushEnabled                              bool                   `json:"pushEnabled"`
	CheckTimes                               int                    `json:"checkTimes"`
	HTTPHealthParams                         HealthParams           `json:"httpHealthParams"`
	TCPHealthParams                          HealthParams           `json:"tcpHealthParams"`
	MysqlHealthParams                        HealthParams           `json:"mysqlHealthParams"`
	IncrementalList                          []string               `json:"incrementalList"`
	ServerStatusSynchronizationPeriodMillis  int                    `json:"serverStatusSynchronizationPeriodMillis"`
	ServiceStatusSynchronizationPeriodMillis int                    `json:"serviceStatusSynchronizationPeriodMillis"`
	DisableAddIP                             bool                   `json:"disableAddIP"`
	SendBeatOnly                             bool                   `json:"sendBeatOnly"`
	LimitedUrlMap                            map[string]interface{} `json:"limitedUrlMap     "`
	DistroServerExpiredMillis                int                    `json:"distroServerExpiredMillis"`
	PushGoVersion                            string                 `json:"pushGoVersion"`
	PushJavaVersion                          string                 `json:"pushJavaVersion"`
	PushPythonVersion                        string                 `json:"pushPythonVersion"`
	PushCVersion                             string                 `json:"pushCVersion"`
	EnableAuthentication                     bool                   `json:"enableAuthentication"`
	OverriddenServerStatus                   string                 `json:"overriddenServerStatus"`
	DefaultInstanceEphemeral                 bool                   `json:"defaultInstanceEphemeral"`
	HealthCheckWhiteList                     []string               `json:"healthCheckWhiteList"`
	Checksum                                 string                 `json:"checksum"`
}

type HealthParams struct {
	Max    int     `json:"max"`
	Min    int     `json:"min"`
	Factor float64 `json:"factor"`
}

func (ss *SystemService) QuerySwitches() (*SystemConfig, error) {
	resp, err := ss.c.Get(ss.c.GetUrl(v1.SystemSwitchesPath))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response code not ok: %d, body: %s", resp.StatusCode, v1.ReadResponseBody(resp.Body))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sc SystemConfig
	if err := json.Unmarshal(data, &sc); err != nil {
		return nil, err
	}

	return &sc, nil
}
