package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	Insecure       bool          `mapstructure:"insecure"`
	Timeout        time.Duration `mapstructure:"timeout"`
	CertFile       string        `mapstructure:"cert_file"`
	KeyFile        string        `mapstructure:"key_file"`
	Log            bool          `mapstructure:"log"`
	Single         bool          `mapstructure:"single"`
	Duration       string        `mapstructure:"duration"`
	Bytes          string        `mapstructure:"bytes"`
	ResponseStatus string        `mapstructure:"status"`
	Request        string        `mapstructure:"request"`
	Response       string        `mapstructure:"response"`
}

const (
	MethodPost   = "POST"
	MethodGet    = "GET"
	MethodPatch  = "PATCH"
	MethodDelete = "DELETE"
)

var fieldConfig Config
var staticClient *http.Client

func SetClient(c *http.Client) {
	staticClient = c
}

func NewClient(c Config) (*http.Client, error) {
	if len(c.Duration) > 0 {
		fieldConfig.Duration = c.Duration
	} else {
		fieldConfig.Duration = "duration"
	}
	if len(c.Request) > 0 {
		fieldConfig.Request = c.Request
	} else {
		fieldConfig.Request = "request"
	}
	if len(c.Response) > 0 {
		fieldConfig.Response = c.Response
	} else {
		fieldConfig.Response = "response"
	}
	if len(c.Bytes) > 0 {
		fieldConfig.Bytes = c.Bytes
	} else {
		fieldConfig.Bytes = "bytes"
	}
	if len(c.ResponseStatus) > 0 {
		fieldConfig.ResponseStatus = c.ResponseStatus
	} else {
		fieldConfig.ResponseStatus = "status"
	}
	if len(c.CertFile) > 0 && len(c.KeyFile) > 0 {
		return NewTLSClient(c.CertFile, c.KeyFile, c.Timeout)
	} else {
		if c.Timeout > 0 {
			transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: c.Insecure},}
			client0 := &http.Client{Transport: transport, Timeout: c.Timeout * time.Second}
			staticClient = client0
			return client0, nil
		} else {
			transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: c.Insecure}}
			client0 := &http.Client{Transport: transport}
			staticClient = client0
			return client0, nil
		}
	}
}
func NewTLSClient(certFile, keyFile string, timeout time.Duration) (*http.Client, error) {
	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	if timeout <= 0 {
		client0 := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					Certificates:       []tls.Certificate{clientCert},
					MinVersion:         tls.VersionTLS10,
					MaxVersion:         tls.VersionTLS10,
				},
			},
		}
		staticClient = client0
		return client0, nil
	} else {
		client0 := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					Certificates:       []tls.Certificate{clientCert},
					MinVersion:         tls.VersionTLS10,
					MaxVersion:         tls.VersionTLS10,
				},
			},
			Timeout: timeout * time.Second,
		}
		staticClient = client0
		return client0, nil
	}
}
func Do(ctx context.Context, client *http.Client, url string, method string, body *[]byte, headers *map[string]string) (*http.Response, error) {
	if body != nil {
		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(*body))
		if err != nil {
			return nil, err
		}
		if headers != nil {
			for k, v := range *headers {
				req.Header.Add(k, v)
			}
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		return resp, err
	} else {
		req, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, err
		}
		if headers != nil {
			for k, v := range *headers {
				req.Header.Add(k, v)
			}
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		return resp, err
	}
}

func DoGet(ctx context.Context, client *http.Client, url string, headers *map[string]string) (*http.Response, error) {
	return Do(ctx, client, url, MethodGet, nil, headers)
}
func DoDelete(ctx context.Context, client *http.Client, url string, headers *map[string]string) (*http.Response, error) {
	return Do(ctx, client, url, MethodDelete, nil, headers)
}

func DoPost(ctx context.Context, client *http.Client, url string, body []byte, headers *map[string]string) (*http.Response, error) {
	return Do(ctx, client, url, MethodPost, &body, headers)
}
func DoPatch(ctx context.Context, client *http.Client, url string, body []byte, headers *map[string]string) (*http.Response, error) {
	return Do(ctx, client, url, MethodPatch, &body, headers)
}
func Post(ctx context.Context, url string, obj interface{}, headers *map[string]string) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, MethodPost, url, obj, headers)
}
func Patch(ctx context.Context, url string, obj interface{}, headers *map[string]string) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, MethodPatch, url, obj, headers)
}
func PostWithClient(ctx context.Context, client *http.Client, url string, obj interface{}, headers *map[string]string) (*json.Decoder, error) {
	return DoWithClient(ctx, client, MethodPost, url, obj, headers)
}
func PatchWithClient(ctx context.Context, client *http.Client, url string, obj interface{}, headers *map[string]string) (*json.Decoder, error) {
	return DoWithClient(ctx, client, MethodPatch, url, obj, headers)
}
func DoWithClient(ctx context.Context, client *http.Client, method string, url string, obj interface{}, headers *map[string]string) (*json.Decoder, error) {
	b, ok := obj.([]byte)
	if ok {
		res, er1 := Do(ctx, client, url, method, &b, headers)
		if er1 != nil {
			return nil, er1
		}
		buf := new(bytes.Buffer)
		_, er2 := buf.ReadFrom(res.Body)
		if er2 != nil {
			return nil, er2
		}
		s := buf.String()
		return json.NewDecoder(strings.NewReader(s)), nil
	} else {
		rq, er0 := json.Marshal(obj)
		if er0 != nil {
			return nil, er0
		}
		res, er1 := Do(ctx, client, url, method, &rq, headers)
		if er1 != nil {
			return nil, er1
		}
		buf := new(bytes.Buffer)
		_, er2 := buf.ReadFrom(res.Body)
		if er2 != nil {
			return nil, er2
		}
		s := buf.String()
		return json.NewDecoder(strings.NewReader(s)), nil
	}
}
