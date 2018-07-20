package ws

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type WebError struct {
	Cause  error
	URL    string
	Status int
}

func (e *WebError) Error() string {
	msg := "cause is nil"
	if e.Cause != nil {
		msg = e.Cause.Error()
	}
	return fmt.Sprintf("error calling url; url=%s; status=%d; msg=%s", e.URL, e.Status, msg)
}

func NewWebError(cause error, url string, status int) *WebError {
	return &WebError{Cause: cause, URL: url, Status: status}
}

const (
	contentJSON = "application/json"
	contentXML  = "application/xml"
)

type Caller struct {
	Method      string
	URL         string
	contentType string
	user        string
	pwd         string
	headers     []keyValue
	webErr      *WebError
}

type keyValue struct {
	key   string
	value string
}

type input struct {
	contentType string
	buf         *bytes.Buffer
}

func NewCaller(method, url string) *Caller {
	return &Caller{
		Method:  method,
		URL:     url,
		webErr:  nil,
		headers: []keyValue{},
	}
}

func (c *Caller) Header(name, value string) *Caller {
	c.headers = append(c.headers, keyValue{key: name, value: value})
	return c
}

func (c *Caller) JSON() *Caller {
	c.contentType = contentJSON
	return c
}

func (c *Caller) XML() *Caller {
	c.contentType = contentXML
	return c
}

func (c *Caller) Auth(user, pwd string) *Caller {
	c.user = user
	c.pwd = pwd
	return c
}

func (c *Caller) getInBuffer(in interface{}) (*bytes.Buffer, *WebError) {
	buf := new(bytes.Buffer)
	if in == nil {
		return buf, nil
	}
	if c.contentType == contentJSON {
		if err := json.NewEncoder(buf).Encode(in); err != nil {
			return nil, NewWebError(errors.Wrapf(err, "failed encode json"), c.URL, http.StatusInternalServerError)
		}
	} else if c.contentType == contentXML {
		if err := xml.NewEncoder(buf).Encode(in); err != nil {
			return nil, NewWebError(errors.Wrapf(err, "failed encode xml"), c.URL, http.StatusInternalServerError)
		}
	}
	return buf, nil
}

func (c *Caller) Call(client *http.Client, in interface{}, out interface{}) *WebError {
	if c.webErr != nil {
		return c.webErr
	}

	buf, webErr := c.getInBuffer(in)
	if webErr != nil {
		return webErr
	}

	req, err := http.NewRequest(c.Method, c.URL, buf)
	if err != nil {
		return NewWebError(errors.Wrap(err, "failed to create request"), c.URL, http.StatusInternalServerError)
	}
	if in != nil && c.contentType != "" {
		req.Header.Set("Content-Type", c.contentType)
	}
	if out != nil && c.contentType != "" {
		req.Header.Set("Accept", c.contentType)
	}
	if c.user != "" || c.pwd != "" {
		req.SetBasicAuth(c.user, c.pwd)
	}
	for _, kv := range c.headers {
		fmt.Printf("add header: %s=%s\n", kv.key, kv.value)
		req.Header.Set(kv.key, kv.value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return NewWebError(errors.Wrap(err, "failed to call url"), c.URL, http.StatusBadGateway)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		DiscardBody(resp)
		// todo: log error
		return NewWebError(errors.Errorf("error code response"), c.URL, resp.StatusCode)
	}

	if out == nil {
		DiscardBody(resp)
		return nil
	}

	if c.contentType == contentXML {
		if err := xml.NewDecoder(resp.Body).Decode(out); err != nil {
			resp.Body.Close()
			return NewWebError(errors.Wrapf(err, "failed decode response"), c.URL, http.StatusInternalServerError)
		}
	} else if c.contentType == contentJSON {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			resp.Body.Close()
			return NewWebError(errors.Wrapf(err, "failed decode response"), c.URL, http.StatusInternalServerError)
		}
	}
	resp.Body.Close()
	return nil
}

func DiscardBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}
