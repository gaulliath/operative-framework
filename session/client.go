package session

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type OpfClient struct {
	Client  http.Client
	Header  Headers
	Data    []byte
	WithTor bool
}

type Headers map[string]string

func GetOpfClient() OpfClient {
	return OpfClient{
		WithTor: false,
		Header:  make(Headers),
	}
}

func (c *OpfClient) SetUserAgent(user string) *OpfClient {
	c.Header["User-Agent"] = user
	return c
}

func (c *OpfClient) SetData(data interface{}) (*OpfClient, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return c, err
	}
	c.Data = buf.Bytes()
	return c, nil
}

func (header Headers) Add(key string, value string) {
	header[key] = value
}

func (c *OpfClient) Perform(method string, uri string) (*http.Response, error) {

	if method == "" {
		method = "GET"
	}

	proxy, err := url.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		return nil, err
	}

	if strings.ToLower(os.Getenv("WITH_TOR")) == "true" {
		c.Client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		c.Client.Timeout = time.Second * 30
	}

	_, err = url.Parse(uri)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(strings.ToUpper(method), uri, nil)

	switch method {
	case "POST":
		req, _ = http.NewRequest(strings.ToUpper(method), uri, bytes.NewBuffer(c.Data))
		break
	}

	for header, content := range c.Header {
		req.Header.Set(header, content)
	}
	return c.Client.Do(req)
}
