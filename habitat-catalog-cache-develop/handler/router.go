package handler

import (
	"net/http"

	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

// Env is the application-wid configuration.
type Env struct {
	Logger  *zerolog.Logger
	Service model.Service
}

// New returns a http.Server capable router with all routing handlers.
func New(conf *config.BasicAuth, logger *zerolog.Logger, service model.Service) (*httprouter.Router, error) {
	e := &Env{
		Logger:  logger,
		Service: service,
	}

	mux := httprouter.New()
	mux.PanicHandler = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		e.Logger.Error().Fields(map[string]interface{}{
			"from":   r.RemoteAddr,
			"path":   r.RequestURI,
			"method": r.Method,
		}).Msgf("panic:%v", v)
	}

	mux.POST("/api/v1/products", Middleware(e, BasicAuth(conf.User, conf.Pwd, SyncProducts)))
	mux.GET("/api/v1/products", Middleware(e, GetProducts))
	mux.GET("/api/v1/brands/:id", Middleware(e, GetBrand))
	mux.POST("/api/v1/brands/:id", Middleware(e, BasicAuth(conf.User, conf.Pwd, SyncBrand)))
	mux.GET("/api/v1/stores", Middleware(e, GetStores))
	mux.POST("/api/v1/stores/:id", Middleware(e, BasicAuth(conf.User, conf.Pwd, SyncStore)))

	// status check endpoint for k8s
	mux.GET("/api/v1/status", Middleware(e, GetStatus))

	return mux, nil
}
