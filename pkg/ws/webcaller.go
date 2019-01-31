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

const (
	ContentJSON  = "application/json"
	ContentXML   = "application/xml"
	ContentPlain = "text/plain"
)

type Caller struct {
	Method      string
	URL         string
	contentType string
	accept      string
	user        string
	pwd         string
	headers     map[string]string
	//outHeaders []keyValue
	//webErr      *WebError
}

// type keyValue struct {
// 	key   string
// 	value string
// }

type input struct {
	contentType string
	buf         *bytes.Buffer
}

func NewCaller(method, url string) *Caller {
	return &Caller{
		Method:  method,
		URL:     url,
		headers: make(map[string]string),
	}
}

func (c *Caller) SetHeader(name, value string) *Caller {
	//c.headers = append(c.headers, keyValue{key: name, value: value})
	c.headers[name] = value
	return c
}

func (c *Caller) GetHeader(name string) (string, bool) {
	// c.headers = append(c.headers, keyValue{key: name, value: value})
	v, found := c.headers[name]
	return v, found
	//return "", c.headers[name]
}

func (c *Caller) Accept(t string) *Caller {
	c.accept = t
	return c
}
func (c *Caller) Content(t string) *Caller {
	c.contentType = t
	return c
}

func (c *Caller) JSON() *Caller {
	c.contentType = ContentJSON
	c.accept = ContentJSON
	return c
}

func (c *Caller) XML() *Caller {
	c.contentType = ContentXML
	c.accept = ContentXML
	return c
}

func (c *Caller) Plain() *Caller {
	c.contentType = ContentPlain
	c.accept = ContentPlain
	return c
}

func (c *Caller) Auth(user, pwd string) *Caller {
	c.user = user
	c.pwd = pwd
	return c
}

func (c *Caller) getInBuffer(in interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if in == nil {
		return buf, nil
	}
	if c.contentType == ContentJSON {
		if err := json.NewEncoder(buf).Encode(in); err != nil {
			return nil, errors.Wrapf(NewWebError(err, c.URL, http.StatusInternalServerError), "failed encode json")
		}
	} else if c.contentType == ContentXML {
		if err := xml.NewEncoder(buf).Encode(in); err != nil {
			return nil, errors.Wrapf(NewWebError(err, c.URL, http.StatusInternalServerError), "failed encode xml")
		}
	}
	return buf, nil
}

func (c *Caller) Call(client *http.Client, in interface{}, out interface{}) error {
	buf, webErr := c.getInBuffer(in)
	if webErr != nil {
		return webErr
	}

	req, err := http.NewRequest(c.Method, c.URL, buf)
	if err != nil {
		return errors.Wrapf(NewWebError(err, c.URL, http.StatusInternalServerError), "failed to create request")
		//return NewWebError(errors.Wrap(err, "failed to create request"), c.URL, http.StatusInternalServerError)
	}
	if in != nil && c.contentType != "" {
		req.Header.Set("Content-Type", c.contentType)
	}
	if out != nil && c.accept != "" {
		req.Header.Set("Accept", c.accept)
	}
	if c.user != "" || c.pwd != "" {
		req.SetBasicAuth(c.user, c.pwd)
	}
	for k, v := range c.headers {
		fmt.Printf("add header: %s=%s\n", k, v)
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(NewWebError(err, c.URL, http.StatusBadGateway), "failed to call url")
	}
	defer resp.Body.Close()

	// save out headers
	for k, v := range resp.Header {
		if len(v) != 1 {
			continue // todo: currently only handles single value headers
		}
		c.headers[k] = v[0]
		//c.headers = append(c.headers, keyValue{key: k, value: v[0]})
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		DiscardBody(resp)
		// todo: log reply body
		return errors.Wrapf(NewWebError(fmt.Errorf("http error code reply; code=%d", resp.StatusCode), c.URL, resp.StatusCode), "")
	}

	if out == nil {
		DiscardBody(resp)
		return nil
	}

	if c.accept == ContentXML {
		if err := xml.NewDecoder(resp.Body).Decode(out); err != nil {
			return errors.Wrapf(NewWebError(err, c.URL, http.StatusInternalServerError), "failed decode response")
		}
	} else if c.accept == ContentJSON {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return errors.Wrapf(NewWebError(err, c.URL, http.StatusInternalServerError), "failed decode response")
		}
	} else if c.accept == ContentPlain {
		if v, ok := out.(*string); ok {
			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrapf(NewWebError(err, c.URL, http.StatusInternalServerError), "failed read text/plain response")
			}
			*v = string(buf)
		}
	}
	return nil
}

func DiscardBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}
