package session

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type OpfClient struct {
	Client  http.Client
	Header  Headers
	Data    interface{}
	WithTor bool
}

type Headers map[string]string

func GetOpfClient() OpfClient {
	return OpfClient{
		WithTor: false,
		Header:  make(Headers),
	}
}

func (header Headers) Add(key string, value string) {
	header[key] = value
}

func (c *OpfClient) Perform(method string, uri string, body io.Reader) (*http.Response, error) {

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

	req, _ := http.NewRequest(strings.ToUpper(method), uri, body)

	for header, content := range c.Header {
		req.Header.Set(header, content)
	}
	return c.Client.Do(req)
}
