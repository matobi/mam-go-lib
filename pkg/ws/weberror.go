package ws

import (
	"fmt"
	"net/http"
)

type WebError struct {
	Cause error
	Msg   string
	URL   string
	Code  int
}

func NewWebErrorMsg(cause error, msg string, code int) error {
	return &WebError{Cause: cause, Msg: msg, Code: code}
}

func NewWebError(cause error, url string, code int) error {
	return &WebError{Cause: cause, URL: url, Code: code}
}

func (e *WebError) Error() string {
	return fmt.Sprintf("error calling url; code=%d; msg=%s; url=%s; cause=%v", e.Code, e.Msg, e.URL, e.Cause)
}

func (e *WebError) GetErrCode() int {
	return e.Code
}

type errcode interface {
	GetErrCode() int
}

// IsTemporary returns true if err is temporary.
func GetErrCode(err error) int {
	if err == nil {
		return 0
	}

	if ec, ok := err.(errcode); ok {
		return ec.GetErrCode()
	}
	return http.StatusInternalServerError // unknown error
}
