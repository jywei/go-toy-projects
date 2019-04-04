package handler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/honestbee/habitat-catalog-cache/handler"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type syncMockService struct{}

func (*syncMockService) GetExternalKeys() ([]string, error) { return []string{}, nil }
func (*syncMockService) SyncCatalog(*model.Catalog) error   { return nil }
func (*syncMockService) GetBrand(int) (*model.Brand, error) { return nil, nil }
func (*syncMockService) SyncBrand(*model.Brand) error       { return nil }
func (*syncMockService) PushJob(job *model.Job) error {
	switch job.Type {
	case model.ExternalKeyJob:
		var key string
		job.GetValue(&key)
		switch key {
		case "cause-error-external-key":
			return errors.New("error occured")
		}
	case model.StoreIDJob:
		v := new(model.StoreIDJobValue)
		job.GetValue(v)
		switch v.ExternalKey {
		case "cause-error-external-key":
			return errors.New("error occured")
		}
	case model.BrandIDJob:
		var v int
		job.GetValue(&v)
		switch v {
		case causePushJobFailedBrandID:
			return errors.New("error occured")
		}
	}
	return nil
}
func (*syncMockService) PopJob() (*model.Job, error)                         { return nil, nil }
func (*syncMockService) SyncProduct(*model.Product) error                    { return nil }
func (*syncMockService) GetProducts([]int) ([]*model.Product, error)         { return nil, nil }
func (*syncMockService) SyncBarcode([]*model.Barcode) error                  { return nil }
func (*syncMockService) GetProductIDs(int, []string) ([]int, error)          { return nil, nil }
func (*syncMockService) GetStoreByBrandID(int) (*model.Store, error)         { return nil, nil }
func (*syncMockService) SelectStores(string, string) ([]*model.Store, error) { return nil, nil }
func (*syncMockService) GetStoreByID(id int) (*model.Store, error)           { return nil, nil }
func (*syncMockService) SyncStore(*model.Store) error                        { return nil }
func (*syncMockService) Close() error                                        { return nil }
func (*syncMockService) UpdateSyncStatus(*model.SyncStatus) error            { return nil }
func (*syncMockService) GetSyncStatus() (*model.SyncStatus, error)           { return nil, nil }

var syncEnv *handler.Env
var syncService *syncMockService

func init() {
	logger := zerolog.New(ioutil.Discard)
	syncService = new(syncMockService)
	syncEnv = &handler.Env{
		Logger:  &logger,
		Service: syncService,
	}
}

const (
	causePushJobFailedBrandID = 449449
)

func TestSyncBrand(t *testing.T) {
	testCases := []struct {
		params httprouter.Params
		want   *handler.Error
	}{
		{
			params: httprouter.Params{httprouter.Param{}},
			want:   &handler.Error{OutputErr: handler.InvalidAttributeErrMsg},
		},
		{
			params: httprouter.Params{httprouter.Param{Key: "id", Value: strconv.Itoa(causePushJobFailedBrandID)}},
			want:   &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
		{
			params: httprouter.Params{httprouter.Param{Key: "id", Value: "3345678"}},
			want:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with brand id:%v", tc.params), func(t *testing.T) {
			_, got := handler.SyncBrand(syncEnv, tc.params, nil)
			switch tc.want {
			case nil:
				if got != nil {
					t.Errorf("error want:nil, got:%v", got)
				}
			default:
				if err, ok := got.(*handler.Error); ok {
					if tc.want.OutputErr != err.OutputErr {
						t.Errorf("error want:%s, got:%s", tc.want.OutputErr, err.OutputErr)
					}
				} else {
					t.Errorf("cast error:%v failed", got)
				}
			}
		})
	}
}

func TestSyncStore(t *testing.T) {
	testCases := []struct {
		request *http.Request
		params  httprouter.Params
		want    *handler.Error
	}{
		{
			request: &http.Request{PostForm: url.Values{}},
			params:  httprouter.Params{httprouter.Param{}},
			want:    &handler.Error{OutputErr: handler.InvalidAttributeErrMsg},
		},
		{
			request: &http.Request{PostForm: url.Values{}},
			params:  httprouter.Params{httprouter.Param{Key: "id", Value: "3345678"}},
			want:    &handler.Error{OutputErr: handler.InvalidAttributeErrMsg},
		},
		{
			request: &http.Request{PostForm: url.Values{"externalKey": []string{"3345678"}}},
			params:  httprouter.Params{httprouter.Param{}},
			want:    &handler.Error{OutputErr: handler.InvalidAttributeErrMsg},
		},
		{
			request: &http.Request{PostForm: url.Values{"externalKey": []string{"3345678"}}},
			params:  httprouter.Params{httprouter.Param{Key: "id", Value: "connot convert id"}},
			want:    &handler.Error{OutputErr: handler.InvalidAttributeErrMsg},
		},
		{
			request: &http.Request{PostForm: url.Values{"externalKey": []string{"cause-error-external-key"}}},
			params:  httprouter.Params{httprouter.Param{Key: "id", Value: "3345678"}},
			want:    &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
		{
			request: &http.Request{PostForm: url.Values{"externalKey": []string{"3345678"}}},
			params:  httprouter.Params{httprouter.Param{Key: "id", Value: "3345678"}},
			want:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with external key:%q, store id:%v",
			tc.request.PostForm["externalKey"], tc.params), func(t *testing.T) {
			_, got := handler.SyncStore(syncEnv, tc.params, tc.request)
			switch tc.want {
			case nil:
				if got != nil {
					t.Errorf("error want:nil, got:%v", got)
				}
			default:
				if err, ok := got.(*handler.Error); ok {
					if tc.want.OutputErr != err.OutputErr {
						t.Errorf("error want:%s, got:%s", tc.want.OutputErr, err.OutputErr)
					}
				} else {
					t.Errorf("cast error:%v failed", got)
				}
			}
		})
	}
}

func TestSyncProducts(t *testing.T) {
	testCases := []struct {
		request *http.Request
		want    *handler.Error
	}{
		{
			request: &http.Request{PostForm: url.Values{}},
			want:    nil,
		},
		{
			request: &http.Request{PostForm: url.Values{
				"externalKeys": []string{"cause-error-external-key"},
			}},
			want: &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
		{
			request: &http.Request{PostForm: url.Values{
				"external_keys": []string{"cause-error-external-key"},
			}},
			want: &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with external keys:%q", tc.request.PostForm["externalKey"]), func(t *testing.T) {
			_, got := handler.SyncProducts(syncEnv, nil, tc.request)
			switch tc.want {
			case nil:
			default:
				if err, ok := got.(*handler.Error); ok {
					if tc.want.OutputErr != err.OutputErr {
						t.Errorf("error want:%s, got:%s", tc.want.OutputErr, err.OutputErr)
					}
				} else {
					t.Errorf("cast error:%v failed", got)
				}
			}
		})
	}
}
