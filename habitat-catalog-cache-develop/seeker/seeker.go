package seeker

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/pkg/errors"
)

const (
	// StoresURL is the backend stores endpoint url.
	StoresURL = "/api/stores"
	// BrandsURL is the backend brands endpoint url.
	BrandsURL = "/api/brands"
)

const (
	version1 = "application/vnd.honestbee+json;version=1"
	version2 = "application/vnd.honestbee+json;version=2"
)

// Seeker is for seeking the backend catalog products informations.
type Seeker interface {
	FetchStores(externalKey string) ([]*model.Store, error)
	FetchBrand(brandID int) (*model.Brand, error)
	FetchProducts(storeID, page int) ([]*model.Product, int, error)
}

type seeker struct {
	domain      string
	retryTimes  int
	retryPeriod time.Duration
	client      *http.Client
}

// New returns a Seeker instance.
func New(conf *config.Seeker) (Seeker, error) {
	return &seeker{
		domain:      conf.FetchDomain,
		retryTimes:  conf.RetryTimes,
		retryPeriod: time.Duration(conf.RetryPeriodSec) * time.Second,
		client: &http.Client{
			Timeout: time.Duration(conf.TimeoutSec) * time.Second,
		},
	}, nil
}

func (s *seeker) connect(dest interface{}, expectStatus int, req *http.Request) error {
	var err error
	var resp *http.Response
	for i := 0; i < s.retryTimes; i++ {
		resp, err = s.client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(s.retryPeriod)
	}
	if err != nil {
		return errors.Wrapf(err, "seeker: [connect] url[%s] http client do failed", req.URL.String())
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectStatus {
		return errors.Errorf("seeker: [connect] url[%s] status expect[%v], actual[%v]",
			req.URL.String(),
			expectStatus,
			resp.Status,
		)
	}

	if err = json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return errors.Wrapf(err, "seeker: [connect] url[%s] json decode failed", req.URL.String())
	}

	return nil
}

func (s *seeker) connectGET(dest interface{}, url, version string, expectStatus int, params io.Reader) error {
	req, err := http.NewRequest(http.MethodGet, url, params)
	if err != nil {
		return errors.Wrapf(err, "seeker: [connectGET] url[%s] http NewRequest failed", url)
	}

	req.Header.Set("Accept", version)
	return errors.Wrapf(
		s.connect(dest, expectStatus, req),
		"seeker: [connectGET] url[%s] connect failed",
		url,
	)
}

// FetchStores returns stores information by external key.
func (s *seeker) FetchStores(externalKey string) ([]*model.Store, error) {
	url := s.domain + StoresURL + "?externalKey=" + externalKey

	ret := make([]*model.Store, 0)
	if err := s.connectGET(&ret, url, version1, http.StatusOK, nil); err != nil {
		return nil, errors.Wrapf(err, "seeker: [FetchStores] connect failed")
	}

	return ret, nil
}

// FetchBrand returns brand information by brand id.
func (s *seeker) FetchBrand(brandID int) (*model.Brand, error) {
	url := s.domain + BrandsURL + "/" + strconv.Itoa(brandID)

	ret := new(model.Brand)
	if err := s.connectGET(ret, url, version1, http.StatusOK, nil); err != nil {
		return nil, errors.Wrapf(err, "seeker: [FetchBrand] connect failed")
	}
	return ret, nil
}

// FetchProducts returns products information and total pages by store id.
func (s *seeker) FetchProducts(storeID, page int) ([]*model.Product, int, error) {
	url := s.domain + StoresURL + "/" + strconv.Itoa(storeID) + "?page=" + strconv.Itoa(page)
	ret := struct {
		Products []*model.Product `json:"products"`
		Meta     struct {
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
			TotalCount  int `json:"total_count"`
		}
	}{}
	if err := s.connectGET(&ret, url, version2, http.StatusOK, nil); err != nil {
		return nil, 0, errors.Wrapf(err, "seeker: [FetchProducts] connect failed")
	}

	return ret.Products, ret.Meta.TotalPages, nil
}
