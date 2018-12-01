// refer to github.com/astaxie/beego/httplib
package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pingcap/errors"
)

var defaultUserAgent = "DefaultHttpUser"

// Get returns *HttpRequest with GET method.
func Get(url string) *HttpRequest {
	var req http.Request
	req.Method = "GET"
	req.Header = http.Header{}
	req.Header.Set("User-Agent", defaultUserAgent)
	return &HttpRequest{url, &req, map[string]string{}, false, 60 * time.Second, 60 * time.Second, nil}
}

// Post returns *HttpRequest with POST method.
func Post(url string) *HttpRequest {
	var req http.Request
	req.Method = "POST"
	req.Header = http.Header{}
	req.Header.Set("User-Agent", defaultUserAgent)
	return &HttpRequest{url, &req, map[string]string{}, false, 60 * time.Second, 60 * time.Second, nil}
}

// Put returns *HttpRequest with PUT method.
func Put(url string) *HttpRequest {
	var req http.Request
	req.Method = "PUT"
	req.Header = http.Header{}
	req.Header.Set("User-Agent", defaultUserAgent)
	return &HttpRequest{url, &req, map[string]string{}, false, 60 * time.Second, 60 * time.Second, nil}
}

// Delete returns *HttpRequest DELETE GET method.
func Delete(url string) *HttpRequest {
	var req http.Request
	req.Method = "DELETE"
	req.Header = http.Header{}
	req.Header.Set("User-Agent", defaultUserAgent)
	return &HttpRequest{url, &req, map[string]string{}, false, 60 * time.Second, 60 * time.Second, nil}
}

// Head returns *HttpRequest with HEAD method.
func Head(url string) *HttpRequest {
	var req http.Request
	req.Method = "HEAD"
	req.Header = http.Header{}
	req.Header.Set("User-Agent", defaultUserAgent)
	return &HttpRequest{url, &req, map[string]string{}, false, 60 * time.Second, 60 * time.Second, nil}
}

// HttpRequest provides more useful methods for requesting one url than http.Request.
type HttpRequest struct {
	url              string
	req              *http.Request
	params           map[string]string
	showdebug        bool
	connectTimeout   time.Duration
	readWriteTimeout time.Duration
	tlsClientConfig  *tls.Config
}

// Debug sets show debug or not when executing request.
func (b *HttpRequest) Debug(isdebug bool) *HttpRequest {
	b.showdebug = isdebug
	return b
}

// SetTimeout sets connect time out and read-write time out for ArthasRequest.
func (b *HttpRequest) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *HttpRequest {
	b.connectTimeout = connectTimeout
	b.readWriteTimeout = readWriteTimeout
	return b
}

// SetTLSClientConfig sets tls connection configurations if visiting https url.
func (b *HttpRequest) SetTLSClientConfig(config *tls.Config) *HttpRequest {
	b.tlsClientConfig = config
	return b
}

// Header add header item string in request.
func (b *HttpRequest) Header(key, value string) *HttpRequest {
	b.req.Header.Set(key, value)
	return b
}

// SetCookie add cookie into request.
func (b *HttpRequest) SetCookie(cookie *http.Cookie) *HttpRequest {
	b.req.Header.Add("Cookie", cookie.String())
	return b
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (b *HttpRequest) Param(key, value string) *HttpRequest {
	b.params[key] = value
	return b
}

// Body adds request raw body.
// it supports string and []byte.
func (b *HttpRequest) Body(data interface{}) *HttpRequest {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		b.req.Body = ioutil.NopCloser(bf)
		b.req.ContentLength = int64(len(t))
	}
	return b
}

func (b *HttpRequest) getResponse() (*http.Response, error) {
	var paramBody string
	if len(b.params) > 0 {
		var buf bytes.Buffer
		for k, v := range b.params {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
			buf.WriteByte('&')
		}
		paramBody = buf.String()
		paramBody = paramBody[0 : len(paramBody)-1]
	}

	if b.req.Method == "GET" && len(paramBody) > 0 {
		if strings.Index(b.url, "?") != -1 {
			b.url += "&" + paramBody
		} else {
			b.url = b.url + "?" + paramBody
		}
	} else if b.req.Method == "POST" && b.req.Body == nil && len(paramBody) > 0 {
		b.Header("Content-Type", "application/x-www-form-urlencoded")
		b.Body(paramBody)
	}

	url, err := url.Parse(b.url)
	if url.Scheme == "" {
		b.url = "http://" + b.url
		url, err = url.Parse(b.url)
	}
	if err != nil {
		return nil, errors.Trace(err)
	}

	b.req.URL = url
	if b.showdebug {
		dump, err := httputil.DumpRequest(b.req, true)
		if err != nil {
			println(err.Error())
		}
		println(string(dump))
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: b.tlsClientConfig,
			Dial:            TimeoutDialer(b.connectTimeout, b.readWriteTimeout),
		},
	}
	resp, err := client.Do(b.req)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return resp, nil
}

// String returns the body string in response.
// it calls Response inner.
func (b *HttpRequest) String() (string, error) {
	data, err := b.Bytes()
	if err != nil {
		return "", errors.Trace(err)
	}

	return string(data), nil
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (b *HttpRequest) Bytes() ([]byte, error) {
	resp, err := b.getResponse()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return data, nil
}

// ToFile saves the body data in response to one file.
// it calls Response inner.
func (b *HttpRequest) ToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Trace(err)
	}
	defer f.Close()

	resp, err := b.getResponse()
	if err != nil {
		return errors.Trace(err)
	}
	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// ToJson returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (b *HttpRequest) ToJson(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return errors.Trace(err)
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// ToXml returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (b *HttpRequest) ToXML(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return errors.Trace(err)
	}
	err = xml.Unmarshal(data, v)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// Response executes request client gets response mannually.
func (b *HttpRequest) Response() (*http.Response, error) {
	return b.getResponse()
}

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, errors.Trace(err)
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

// Parse body to data
func FromBody(r *http.Request, data interface{}) (string, bool, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", false, errors.Trace(err)
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return string(body), false, errors.Trace(err)
	}

	return string(body), true, nil
}

// Encode data to json
func ToJson(data interface{}) (string, error) {
	value, err := json.Marshal(data)
	if err != nil {
		return "", errors.Trace(err)
	}

	return string(value), nil
}
