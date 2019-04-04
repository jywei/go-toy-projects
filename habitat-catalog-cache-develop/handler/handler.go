package handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// Handler is a type of function to represent every handling requests function.
type Handler func(e *Env, ps httprouter.Params, r *http.Request) (interface{}, error)

// Middleware pre-handle every incoming request and generate output to the client based on every handler returned error code.
func Middleware(e *Env, fn Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)

		e.Logger.Info().Fields(map[string]interface{}{
			"from":   r.RemoteAddr,
			"path":   r.URL.Path,
			"method": r.Method,
		}).Msgf("receiving data")

		ret, err := fn(e, ps, r)
		if err != nil {
			er, ok := err.(*Error)
			if !ok {
				// Any error types we don't specifically look out for default
				// to serving a HTTP 500.
				er = NewErr(ServerInternalErrCode, err)
			}
			if er.internalErr != nil {
				e.Logger.Error().Fields(map[string]interface{}{
					"from":   r.RemoteAddr,
					"path":   r.URL.Path,
					"method": r.Method,
					"error":  er.Error(),
				}).Msgf("error occurred")
			}
			w.WriteHeader(er.Status)
			ret = er
		}
		encoder.Encode(ret)
	}
}

// BasicAuth is a user and password checking.
func BasicAuth(wantUser, wantPwd string, next Handler) Handler {
	return func(e *Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
		gotUser, gotPwd, ok := r.BasicAuth()
		if !ok {
			return nil, NewErr(UnauthorizedErrCode, errors.New("basic auth parse out not ok"))
		}

		if wantUser != gotUser || wantPwd != gotPwd {
			return nil, NewErr(UnauthorizedErrCode, errors.New("user or password not match"))
		}

		return next(e, ps, r)
	}
}
