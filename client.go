package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
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
	Size           string        `mapstructure:"size"`
	ResponseStatus string        `mapstructure:"status"`
	Request        string        `mapstructure:"request"`
	Response       string        `mapstructure:"response"`
	Error          string        `mapstructure:"error"`
	Fields         string        `mapstructure:"fields"`
}

type FieldConfig struct {
	Fields *[]string `mapstructure:"fields"`
}

const (
	methodPost   = "POST"
	methodPut    = "PUT"
	methodGet    = "GET"
	methodPatch  = "PATCH"
	methodDelete = "DELETE"
)

const FIELDS = "logFields"

var log *zap.Logger
var fieldConfig FieldConfig

func SetLogger(logger *zap.Logger) {
	log = logger
}

var conf Config
var staticClient *http.Client

func SetClient(c *http.Client) {
	staticClient = c
}

func NewClient(c Config) (*http.Client, error) {
	conf.Log = c.Log
	conf.Single = c.Single
	conf.ResponseStatus = c.ResponseStatus
	conf.Size = c.Size
	if len(c.Duration) > 0 {
		conf.Duration = c.Duration
	} else {
		conf.Duration = "duration"
	}
	if len(c.Request) > 0 {
		conf.Request = c.Request
	} else {
		conf.Request = "request"
	}
	if len(c.Response) > 0 {
		conf.Response = c.Response
	} else {
		conf.Response = "response"
	}
	if len(c.Error) > 0 {
		conf.Error = c.Error
	} else {
		conf.Error = "error"
	}
	if len(c.Fields) > 0 {
		fields := strings.Split(c.Fields, ",")
		fieldConfig.Fields = &fields
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
		b := *body
		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(b))
		if err != nil {
			return nil, err
		}
		return AddHeaderAndDo(client, req, headers)
	} else {
		req, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, err
		}
		return AddHeaderAndDo(client, req, headers)
	}
}
func AddHeaderAndDo(client *http.Client, req *http.Request, headers *map[string]string) (*http.Response, error) {
	if headers != nil {
		for k, v := range *headers {
			req.Header.Add(k, v)
		}
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	return resp, err
}
func DoGet(ctx context.Context, client *http.Client, url string, headers *map[string]string) (*http.Response, error) {
	return Do(ctx, client, url, methodGet, nil, headers)
}
func DoDelete(ctx context.Context, client *http.Client, url string, headers *map[string]string) (*http.Response, error) {
	return Do(ctx, client, url, methodDelete, nil, headers)
}
func DoPost(ctx context.Context, client *http.Client, url string, body []byte, headers *map[string]string) (*http.Response, error) {
	return Do(ctx, client, url, methodPost, &body, headers)
}
func DoPut(ctx context.Context, client *http.Client, url string, body []byte, headers *map[string]string) (*http.Response, error) {
	return Do(ctx, client, url, methodPut, &body, headers)
}
func DoPatch(ctx context.Context, client *http.Client, url string, body []byte, headers *map[string]string) (*http.Response, error) {
	return Do(ctx, client, url, methodPatch, &body, headers)
}
func Get(ctx context.Context, url string) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, methodGet, url, nil, nil)
}
func GetWithHeader(ctx context.Context, url string, headers *map[string]string) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, methodGet, url, nil, headers)
}
func GetAndDecode(ctx context.Context, url string, result interface{}) error {
	return GetWithHeaderAndDecode(ctx, url, nil, nil, result)
}
func GetWithHeaderAndDecode(ctx context.Context, url string, obj interface{}, headers *map[string]string, result interface{}) error {
	decoder, er1 := DoWithClient(ctx, staticClient, methodGet, url, obj, headers)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func Delete(ctx context.Context, url string) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, methodDelete, url, nil, nil)
}
func DeleteWithHeader(ctx context.Context, url string, headers *map[string]string) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, methodDelete, url, nil, headers)
}
func DeleteAndDecode(ctx context.Context, url string, result interface{}) error {
	return DeleteWithHeaderAndDecode(ctx, url, nil, nil, result)
}
func DeleteWithHeaderAndDecode(ctx context.Context, url string, obj interface{}, headers *map[string]string, result interface{}) error {
	decoder, er1 := DoWithClient(ctx, staticClient, methodDelete, url, obj, headers)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func Post(ctx context.Context, url string, obj interface{}) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, methodPost, url, obj, nil)
}
func PostWithHeader(ctx context.Context, url string, obj interface{}, headers *map[string]string) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, methodPost, url, obj, headers)
}
func PostAndDecode(ctx context.Context, url string, obj interface{}, result interface{}) error {
	return PostWithHeaderAndDecode(ctx, url, obj, nil, result)
}
func PostWithHeaderAndDecode(ctx context.Context, url string, obj interface{}, headers *map[string]string, result interface{}) error {
	decoder, er1 := DoWithClient(ctx, staticClient, methodPost, url, obj, headers)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func Put(ctx context.Context, url string, obj interface{}) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, methodPut, url, obj, nil)
}
func PutWithHeader(ctx context.Context, url string, obj interface{}, headers *map[string]string) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, methodPut, url, obj, headers)
}
func PutAndDecode(ctx context.Context, url string, obj interface{}, result interface{}) error {
	return PutWithHeaderAndDecode(ctx, url, obj, nil, result)
}
func PutWithHeaderAndDecode(ctx context.Context, url string, obj interface{}, headers *map[string]string, result interface{}) error {
	decoder, er1 := DoWithClient(ctx, staticClient, methodPut, url, obj, headers)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func Patch(ctx context.Context, url string, obj interface{}) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, methodPatch, url, obj, nil)
}
func PatchWithHeader(ctx context.Context, url string, obj interface{}, headers *map[string]string) (*json.Decoder, error) {
	return DoWithClient(ctx, staticClient, methodPatch, url, obj, headers)
}
func PatchAndDecode(ctx context.Context, url string, obj interface{}, result interface{}) error {
	return PatchWithHeaderAndDecode(ctx, url, obj, nil, result)
}
func PatchWithHeaderAndDecode(ctx context.Context, url string, obj interface{}, headers *map[string]string, result interface{}) error {
	decoder, er1 := DoWithClient(ctx, staticClient, methodPatch, url, obj, headers)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func DoWithClient(ctx context.Context, client *http.Client, method string, url string, obj interface{}, headers *map[string]string) (*json.Decoder, error) {
	b, ok := obj.([]byte)
	if ok {
		return DoAndBuildDecoder(ctx, client, url, method, &b, headers)
	} else {
		s, ok2 := obj.(string)
		if ok2 {
			b2 := []byte(s)
			return DoAndBuildDecoder(ctx, client, url, method, &b2, headers)
		} else {
			rq, er0 := json.Marshal(obj)
			if er0 != nil {
				return nil, er0
			}
			return DoAndBuildDecoder(ctx, client, url, method, &rq, headers)
		}
	}
}
func DoAndBuildDecoder(ctx context.Context, client *http.Client, url string, method string, body *[]byte, headers *map[string]string) (*json.Decoder, error) {
	if conf.Log == true && log != nil {
		if !conf.Single && len(conf.Request) > 0 && body != nil {
			fs1 := make([]zap.Field, 0)
			rq := string(*body)
			if len(rq) > 0 {
				f0 := zap.String(conf.Request, rq)
				fs1 = append(fs1, f0)
			}
			fs1 = AppendFields(ctx, fs1)
			log.Info(method+" "+url, fs1...)
		}
		start := time.Now()
		res, er1 := Do(ctx, client, url, method, body, headers)
		if er1 != nil {
			if conf.Single && len(conf.Request) > 0 {
				fs2 := make([]zap.Field, 0)
				if body != nil {
					rq := string(*body)
					if len(rq) > 0 {
						f0 := zap.String(conf.Request, rq)
						fs2 = append(fs2, f0)
					}
				}
				f1 := zap.String(conf.Error, er1.Error())
				fs2 = append(fs2, f1)
				fs2 = AppendFields(ctx, fs2)
				log.Error(method+" "+url, fs2...)
			}
			return nil, er1
		}
		end := time.Now()
		fs3 := make([]zap.Field, 0)
		f1 := zap.Int64(conf.Duration, end.Sub(start).Milliseconds())
		fs3 = append(fs3, f1)
		if conf.Single && len(conf.Request) > 0 && body != nil {
			rq := string(*body)
			if len(rq) > 0 {
				f2 := zap.String(conf.Request, rq)
				fs3 = append(fs3, f2)
			}
		}
		if len(conf.Size) > 0 {
			f3 := zap.Int64(conf.Size, res.ContentLength)
			fs3 = append(fs3, f3)
		}
		if len(conf.ResponseStatus) > 0 {
			f3 := zap.Int(conf.ResponseStatus, res.StatusCode)
			fs3 = append(fs3, f3)
		}
		buf := new(bytes.Buffer)
		_, er3 := buf.ReadFrom(res.Body)
		if er3 != nil {
			log.Error(method+" "+url, fs3...)
			return nil, er3
		}
		s := buf.String()
		if len(conf.Response) > 0 {
			f3 := zap.String(conf.Response, s)
			fs3 = append(fs3, f3)
		}
		fs3 = AppendFields(ctx, fs3)
		if res.StatusCode == 503 {
			log.Error(method+" "+url, fs3...)
			er2 := errors.New("503 Service Unavailable")
			return nil, er2
		}
		log.Info(method+" "+url, fs3...)
		return json.NewDecoder(strings.NewReader(s)), nil
	} else {
		res, er1 := Do(ctx, client, url, method, body, headers)
		if er1 != nil {
			return nil, er1
		}
		if res.StatusCode == 503 {
			er2 := errors.New("503 Service Unavailable")
			return nil, er2
		}
		return json.NewDecoder(res.Body), nil
	}
}
func AppendFields(ctx context.Context, fields []zap.Field) []zap.Field {
	if logFields, ok := ctx.Value(FIELDS).(map[string]string); ok {
		for k, v := range logFields {
			f := zap.String(k, v)
			fields = append(fields, f)
		}
	}
	if fieldConfig.Fields != nil {
		cfs := *fieldConfig.Fields
		for _, k2 := range cfs {
			if v2, ok := ctx.Value(k2).(string); ok && len(v2) > 0 {
				f := zap.String(k2, v2)
				fields = append(fields, f)
			}
		}
	}
	return fields
}
