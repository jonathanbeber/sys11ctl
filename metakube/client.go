package metakube

import (
	"io"
	"net"
	"net/http"
	"time"
)

// Client wraps http.Client with metakube config
type Client interface {
	// do receives a http request method, URL and body. It returns a
	// http.Response and possible errors.
	do(method, url string, body io.Reader) (*http.Response, error)
}

// NewClient returns a `metakube.Client` given a metakube base url and an API token.
func NewClient(url, credentials string) Client {
	client := client{
		baseURL:     url,
		credentials: credentials,
	}

	// Sets a client timeout
	client.Timeout = time.Second * 10
	client.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	return client
}

type client struct {
	http.Client
	baseURL     string
	credentials string
}

func (c client) do(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.baseURL+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", "Bearer "+c.credentials)
	req.Header.Add("accept", "application/json")
	return c.Do(req)
}
