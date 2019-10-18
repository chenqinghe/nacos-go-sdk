package nacos

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ConfigService struct {
	c *Client
}

func NewConfigService(client *Client) *ConfigService {
	return &ConfigService{
		c: client,
	}
}

func (cs *ConfigService) GetConfig(namespace, group, dataId string) ([]byte, error) {

	vals := make(url.Values)
	vals.Set("tenant", namespace)
	vals.Set("dataId", dataId)
	vals.Set("group", group)
	u := cs.c.GetUrl(config) + "?" + vals.Encode()
	resp, err := cs.c.c.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response code not ok: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (cs *ConfigService) PublishConfig(namespace, group, dataId string, data []byte, typ string) error {
	vals := make(url.Values)
	vals.Set("tenant", namespace)
	vals.Set("group", group)
	vals.Set("dataId", dataId)
	vals.Set("content", string(data))
	vals.Set("type", typ)

	resp, err := cs.c.c.PostForm(cs.c.GetUrl(config), vals)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code not ok: %d", resp.StatusCode)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if r := string(respData); r != "true" {
		return fmt.Errorf("publish config failed, response data:%s", r)
	}

	return nil
}

func (cs *ConfigService) RemoveConfig(namespace, group, dataId string) error {
	vals := make(url.Values)
	vals.Set("tenant", namespace)
	vals.Set("group", group)
	vals.Set("dataId", dataId)
	u := cs.c.GetUrl(config) + "?" + vals.Encode()

	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return err
	}

	resp, err := cs.c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code not ok: %d", resp.StatusCode)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if r := string(respData); r != "true" {
		return fmt.Errorf("publish config failed, response data:%s", r)
	}

	return nil
}

type ListenOption struct {
	RetryInterval  time.Duration
	PullingTimeout time.Duration
}

func (cs *ConfigService) Listen(namespace, group, dataId string, option ...ListenOption) *Listener {
	ch := make(chan []byte)
	errors := make(chan error, 1) // buffered first error
	quit := make(chan struct{})

	go func() {
		var dataMd5 string
		var data []byte
		var timer = time.NewTimer(0)
		var retryInterval = time.Second
		if len(option) > 0 {
			retryInterval = option[0].RetryInterval
		}

		buf := bytes.NewBuffer(nil)

		for {
			select {
			case <-quit:
			default:
			}

			buf.Reset()
			buf.WriteString("Listening-Configs=")
			buf.WriteString(dataId)
			buf.WriteByte(2)
			buf.WriteString(group)
			buf.WriteByte(2)
			buf.WriteString(dataMd5)
			if namespace != "" {
				buf.WriteByte(2)
				buf.WriteString(namespace)
			}
			buf.WriteByte(1)

			u := cs.c.GetUrl(configListener)
			req, _ := http.NewRequest(http.MethodPost, u, bytes.NewReader(buf.Bytes()))
			var timeout string
			if len(option) == 0 {
				timeout = "30000"
			} else {
				timeout = strconv.Itoa(int(option[0].PullingTimeout.Milliseconds()))
			}

			req.Header.Set("Long-Pulling-Timeout", timeout)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := cs.c.c.Do(req)
			if err != nil {
				timer.Reset(retryInterval)
				select {
				case errors <- err:
				case <-quit:
					timer.Stop()
					return
				case <-timer.C:
				}
				continue
			}

			if resp.StatusCode != http.StatusOK {
				timer.Reset(retryInterval)
				select {
				case errors <- fmt.Errorf("response http status code not ok: %d", resp.StatusCode):
				case <-quit:
					timer.Stop()
					return
				case <-timer.C:
				}
				continue
			}

			data, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				timer.Reset(retryInterval)
				select {
				case errors <- err:
				case <-quit:
					timer.Stop()
					return
				case <-timer.C:
				}
				continue
			}
			if len(data) == 0 { // 配置无变化
				timer.Reset(retryInterval)
				select {
				case <-quit:
					timer.Stop()
					return
				case <-timer.C:
				}
				continue
			}

			data, err := cs.GetConfig(namespace, group, dataId)
			if err != nil {
				timer.Reset(retryInterval)
				select {
				case <-quit:
					timer.Stop()
					return
				case <-timer.C:
				}
				continue
			}

			resp.Body.Close()

			// update md5
			h := md5.New()
			h.Write(data)
			dataMd5 = hex.EncodeToString(h.Sum(nil))

			select {
			case ch <- data:
			default:
			}
		}
	}()
	return &Listener{
		data:   ch,
		errors: errors,
		quit:   quit,
	}
}

type Listener struct {
	data   chan []byte
	errors chan error
	quit   chan struct{}
}

func (l *Listener) Data() chan []byte {
	return l.data
}

func (l *Listener) Err() chan error {
	return l.errors
}

func (l *Listener) Stop() {
	close(l.quit)
}
