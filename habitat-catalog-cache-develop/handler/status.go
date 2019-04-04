package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// GetStatus handles k8s server status checks.
func GetStatus(e *Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
	status, err := e.Service.GetSyncStatus()
	if err != nil {
		// because k8s health check depends on http status code,
		// so no matter what returns status OK.
		return nil, NewErr(SuccessCode, errors.Wrapf(err, "handle: [GetStatus] GetSyncStatus failed"))
	}

	return struct {
		GoVersion  string            `json:"goVersion"`
		AppVersion string            `json:"appVersion"`
		ServerTime string            `json:"serverTime"`
		SyncStatus *model.SyncStatus `json:"syncStatus"`
	}{
		GoVersion:  runtime.Version(),
		AppVersion: config.Version,
		ServerTime: time.Now().UTC().String(),
		SyncStatus: status,
	}, nil
}
