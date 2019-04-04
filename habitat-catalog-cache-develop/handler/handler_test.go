package handler_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/honestbee/habitat-catalog-cache/handler"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type midMockService struct{}

func (*midMockService) GetExternalKeys() ([]string, error)                  { return []string{}, nil }
func (*midMockService) SyncCatalog(*model.Catalog) error                    { return nil }
func (*midMockService) GetBrand(int) (*model.Brand, error)                  { return nil, nil }
func (*midMockService) SyncBrand(*model.Brand) error                        { return nil }
func (*midMockService) PushJob(*model.Job) error                            { return nil }
func (*midMockService) PopJob() (*model.Job, error)                         { return nil, nil }
func (*midMockService) SyncProduct(*model.Product) error                    { return nil }
func (*midMockService) GetProducts([]int) ([]*model.Product, error)         { return nil, nil }
func (*midMockService) SyncBarcode([]*model.Barcode) error                  { return nil }
func (*midMockService) GetProductIDs(int, []string) ([]int, error)          { return nil, nil }
func (*midMockService) GetStoreByBrandID(int) (*model.Store, error)         { return nil, nil }
func (*midMockService) SelectStores(string, string) ([]*model.Store, error) { return nil, nil }
func (*midMockService) GetStoreByID(id int) (*model.Store, error)           { return nil, nil }
func (*midMockService) SyncStore(*model.Store) error                        { return nil }
func (*midMockService) Close() error                                        { return nil }
func (*midMockService) UpdateSyncStatus(*model.SyncStatus) error            { return nil }
func (*midMockService) GetSyncStatus() (*model.SyncStatus, error)           { return nil, nil }

var midEnv *handler.Env
var midService *midMockService

func init() {
	logger := zerolog.New(ioutil.Discard)
	midService = new(midMockService)
	midEnv = &handler.Env{
		Logger:  &logger,
		Service: midService,
	}
}

func TestMiddleware(t *testing.T) {
	testCases := []struct {
		description string
		fn          handler.Handler
		wantStatus  int
		wantBody    map[string]interface{}
	}{
		{
			description: "status 200",
			fn: func(e *handler.Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
				return nil, handler.NewErr(handler.SuccessCode, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   map[string]interface{}{"error": handler.SuccessMsg},
		},
		{
			description: "status 500",
			fn: func(e *handler.Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
				return nil, handler.NewErr(handler.ServerInternalErrCode, nil)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   map[string]interface{}{"error": handler.ServerInternalErrMsg},
		},
		{
			description: "status 500 with not handler error type",
			fn: func(e *handler.Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
				return nil, fmt.Errorf("error occured")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   map[string]interface{}{"error": handler.ServerInternalErrMsg},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.description), func(t *testing.T) {
			h := handler.Middleware(midEnv, tc.fn)
			r := httptest.NewRequest(http.MethodGet, "http://fake.url.com", nil)
			w := httptest.NewRecorder()
			h(w, r, nil)

			resp := w.Result()
			gotStatus := resp.StatusCode
			if tc.wantStatus != gotStatus {
				t.Errorf("status want:%d, got:%d", tc.wantStatus, gotStatus)
			}
			gotBody := make(map[string]interface{})
			if err := json.NewDecoder(resp.Body).Decode(&gotBody); err != nil {
				t.Fatalf("unmarshal body failed:%v", err)
			}
			if !reflect.DeepEqual(tc.wantBody, gotBody) {
				t.Errorf("body want: %v, got:%v", tc.wantBody, gotBody)
			}
		})
	}
}
