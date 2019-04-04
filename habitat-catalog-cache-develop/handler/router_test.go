package handler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/handler"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/rs/zerolog"
)

type routerMockService struct{}

func (*routerMockService) GetExternalKeys() ([]string, error)          { return []string{}, nil }
func (*routerMockService) SyncCatalog(*model.Catalog) error            { return nil }
func (*routerMockService) GetBrand(int) (*model.Brand, error)          { return nil, nil }
func (*routerMockService) SyncBrand(*model.Brand) error                { return nil }
func (*routerMockService) PushJob(*model.Job) error                    { return nil }
func (*routerMockService) PopJob() (*model.Job, error)                 { return nil, nil }
func (*routerMockService) SyncProduct(*model.Product) error            { return nil }
func (*routerMockService) GetProducts([]int) ([]*model.Product, error) { return nil, nil }
func (*routerMockService) SyncBarcode([]*model.Barcode) error          { return nil }
func (*routerMockService) GetProductIDs(int, []string) ([]int, error)  { return nil, nil }
func (*routerMockService) GetStoreByBrandID(int) (*model.Store, error) { return nil, nil }
func (*routerMockService) SelectStores(string, string) ([]*model.Store, error) {
	return make([]*model.Store, 2), nil
}
func (*routerMockService) GetStoreByID(id int) (*model.Store, error) { return nil, nil }
func (*routerMockService) SyncStore(*model.Store) error              { return nil }
func (*routerMockService) Close() error                              { return nil }
func (*routerMockService) UpdateSyncStatus(*model.SyncStatus) error  { return nil }
func (*routerMockService) GetSyncStatus() (*model.SyncStatus, error) { return nil, nil }

var routerService *routerMockService
var logger *zerolog.Logger

func init() {
	l := zerolog.New(ioutil.Discard)
	logger = &l
	routerService = new(routerMockService)
}

func TestNewRouter(t *testing.T) {
	authUser := "tester"
	authPwd := "testing"
	router, err := handler.New(&config.BasicAuth{
		User: authUser,
		Pwd:  authPwd,
	}, logger, routerService)
	if err != nil {
		t.Fatalf("new router failed:%v", err)
	}
	ts := httptest.NewServer(router)
	defer ts.Close()

	testCases := []struct {
		description string
		endpoint    string
		wantStatus  int
		method      string
		setAuthUser string
		setAuthPwd  string
	}{
		{
			description: "sync products endpoint",
			endpoint:    "/api/v1/products",
			wantStatus:  http.StatusOK,
			method:      http.MethodPost,
			setAuthUser: authUser,
			setAuthPwd:  authPwd,
		},
		{
			description: "sync products endpoint with wrong passwrod",
			endpoint:    "/api/v1/products",
			wantStatus:  http.StatusUnauthorized,
			method:      http.MethodPost,
			setAuthUser: authUser,
			setAuthPwd:  "wrong pwd",
		},
		{
			description: "sync store endpoint with wrong passwrod",
			endpoint:    "/api/v1/stores/3345678",
			wantStatus:  http.StatusUnauthorized,
			method:      http.MethodPost,
			setAuthUser: authUser,
			setAuthPwd:  "wrong pwd",
		},
		{
			description: "sync brand endpoint with wrong passwrod",
			endpoint:    "/api/v1/brands/3345678",
			wantStatus:  http.StatusUnauthorized,
			method:      http.MethodPost,
			setAuthUser: authUser,
			setAuthPwd:  "wrong pwd",
		},
		{
			description: "get products endpoint",
			endpoint:    "/api/v1/products?externalKey=3345678",
			wantStatus:  http.StatusOK,
			method:      http.MethodGet,
		},
		{
			description: "get brand endpoint",
			endpoint:    "/api/v1/brands/3345678",
			wantStatus:  http.StatusOK,
			method:      http.MethodGet,
		},
		{
			description: "get stores endpoint",
			endpoint:    "/api/v1/stores",
			wantStatus:  http.StatusOK,
			method:      http.MethodGet,
		},
		{
			description: "get status endpoint",
			endpoint:    "/api/v1/status",
			wantStatus:  http.StatusOK,
			method:      http.MethodGet,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.description), func(t *testing.T) {
			client := ts.Client()
			req, err := http.NewRequest(tc.method, ts.URL+tc.endpoint, nil)
			if err != nil {
				t.Fatalf("new request failed:%v", err)
			}
			if tc.setAuthUser != "" && tc.setAuthPwd != "" {
				req.SetBasicAuth(tc.setAuthUser, tc.setAuthPwd)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("requesting failed:%v", err)
			}
			if resp.StatusCode != tc.wantStatus {
				t.Errorf("status want:%d, got:%d", tc.wantStatus, resp.StatusCode)
			}
		})
	}
}
