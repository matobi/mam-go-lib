package ws

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ReplyJSON Return a json object to http.
func ReplyJSON(w http.ResponseWriter, data interface{}) error {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return NewWebErrorMsg(err, "failed marshal json", http.StatusInternalServerError)
		//errw := errors.Wrap(err, "")
		//log.Error().Err(err).Msg("failed marshal json")
		//return http.StatusInternalServerError, errw
	}
	return ReplyRawJSON(w, b)
}

// ReplyRawJSON Return a json string to http
func ReplyRawJSON(w http.ResponseWriter, rawJSON []byte) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(rawJSON))
	return nil
}

// ReplyError Return a error message and status code to http
func ReplyError(w http.ResponseWriter, err error) {
	log.Info().Err(errors.Cause(err)).Msg("ws error")
	http.Error(w, err.Error(), GetErrCode(err))
}
