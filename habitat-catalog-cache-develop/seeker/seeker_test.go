package seeker_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/honestbee/habitat-catalog-cache/seeker"
)

var (
	normalStoreJSON = []byte(`[
		{
		    "id": 16849,
		    "name": "Habitat Store",
		    "pickupPoint": "",
		    "slug": "habitat-store",
		    "brandId": 9483,
		    "addressId": 36257,
		    "catalogId": 14202,
		    "priority": null,
		    "notes": "",
		    "description": "",
		    "imageUrl": null,
		    "closed": null,
		    "temporarilyClosed": false,
		    "opensAt": null,
		    "estimatedDeliveryTime": null,
		    "bufferTime": 0,
		    "deliveryTypes": [],
		    "shippingModes": [
			"offline"
		    ],
		    "storeType": "habitat",
		    "minimumOrderFreeDelivery": null,
		    "defaultDeliveryFee": null,
		    "freeDeliveryEligible": true,
		    "minimumSpend": null,
		    "minimumSpendExtraFee": "29.99"
		}
	]`)
	normalBrandJSON = []byte(`{
		"id": 9483,
		"name": "Habitat Test",
		"slug": "habitat-test",
		"description": "If you love food, this is your natural habitat",
		"about": "Habitat is a next-generation fresh food and grocery store, but more than that, it's a place to get in touch with, get personal with food. If you love food, this is your natural habitat.",
		"imageUrl": "https://assets-staging.honestbee.com/headers/images/360x120/habitat-test.png",
		"productsImageUrl": "https://assets-staging.honestbee.com/headers/images/360x120/product-list-habitat-test.png",
		"brandColor": "",
		"currency": "SGD",
		"countryId": 1,
		"minimumOrderFreeDelivery": null,
		"defaultDeliveryFee": null,
		"departments": [
		    {
			"id": 29013,
			"name": "Alcohol",
			"description": null,
			"imageUrl": null
		    },
		    {
			"id": 29012,
			"name": "Bakery",
			"description": null,
			"imageUrl": null
		    },
		    {
			"id": 29014,
			"name": "Food Cupboard",
			"description": null,
			"imageUrl": null
		    }
		],
		"brandTraits": [],
		"sameStorePrice": true,
		"brandType": "habitat",
		"promotionText": "",
		"parentBrandId": null,
		"productsCount": 142,
		"storeId": 16849,
		"priceMarkupPercentage": "0.0",
		"freeDeliveryEligible": true,
		"itemReplacementOptions": [
		    "find_best_match",
		    "pick_specific_replacement",
		    "do_not_replace"
		],
		"estimatedDeliveryTime": null,
		"closed": null,
		"tags": [],
		"catalogId": 12131,
		"opensAt": null,
		"minimumSpend": null,
		"minimumSpendExtraFee": "29.99",
		"deliveryTypes": [],
		"shippingModes": [
		    "offline"
		],
		"cashbackAmount": 0,
		"reservedTags": [
		    {
			"id": 101,
			"title": "",
			"key": "habitat",
			"imageUrl": null
		    },
		    {
			"id": 95,
			"title": "Eat More, Earn More",
			"key": "cashback",
			"imageUrl": "https://assets.honestbee.com/food/tags/foodtag-sumo6.jpg"
		    }
		],
		"defaultConciergeFee": null
	    }`)
	normalProductsJSON = []byte(`{
		"products": [
		    {
			"id": 3070177,
			"title": "Dove Hair Fall Rescue Shampoo",
			"description": "",
			"imageUrl": "https://assets.honestbee.com/products/images/480/habitattest_hb90000003ea_hb90000003ea-1.jpg",
			"previewImageUrl": "https://assets.honestbee.com/products/images/480/habitattest_hb90000003ea_hb90000003ea-1.jpg",
			"slug": "",
			"barcodes": [
			    "HB90000003EA"
			],
			"barcode": null,
			"unitType": "unit_type_item",
			"soldBy": "sold_by_item",
			"amountPerUnit": "1.0",
			"size": "",
			"status": "status_available",
			"imageUrlBasename": "habitattest_hb90000003ea_hb90000003ea-1.jpg",
			"currency": "SGD",
			"promotionStartsAt": null,
			"promotionEndsAt": null,
			"maxQuantity": "5183.0",
			"customerNotesEnabled": false,
			"price": "1.0",
			"normalPrice": "3.0",
			"nutritionalInfo": null,
			"productBrand": "",
			"productInfo": null,
			"packingSize": "",
			"descriptionHtml": null,
			"countryOfOrigin": null,
			"tags": [],
			"alcohol": false
		    }
		],
		"meta": {
		    "current_page": 1,
		    "total_pages": 1,
		    "total_count": 1
		}
	    }`)
	causeDecodeFailedStoreJSON = []byte(`[
		{
		    "id": 16849,
		    "name": "Habitat Store",
		    "pickupPoint": "",
		    "slug": "habitat-store",
		    "brandId": 9483,
		    "addressId": 36257,
		    "catalogId": 14202,
		    "priority": null,
		    "notes": "",
		    "description": "",
		    "imageUrl": null,
		    "closed": null,
		    "temporarilyClosed": false,
		    "opensAt": null,
		    "estimatedDeliveryTime": null,
		    "bufferTime": 0,
		    "deliveryTypes": [],
		    "shippingModes": [
			"offline"
		    ],
		    "storeType": "habitat",
		    "minimumOrderFreeDelivery": null,
		    "defaultDeliveryFee": null,
		    "freeDeliveryEligible": true,
		    "minimumSpend": null,
		    "minimumSpendExtraFee": "29.99"
	]`)
)

const (
	storeNotExistExternalKey          = "store-not-exist-external-key"
	storeCauseDecodeFailedExternalKey = "store-cause-decode-failed-external-key"
	storeSuccessExternalKey           = "store-success-external-key"
	brandNotExistID                   = 449449
	brandSuccessID                    = 3345678
	productNotExistStoreID            = 450450
	productSuccessStoreID             = 3345679
)

func TestFetchProducts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storeID, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(r.URL.String(), "/api/stores/"), "?page=1"))
		if err != nil {
			t.Fatalf("receive an unparsable store id")
		}
		switch storeID {
		case productNotExistStoreID:
			w.WriteHeader(http.StatusNotFound)
		case productSuccessStoreID:
			w.WriteHeader(http.StatusOK)
			w.Write(normalProductsJSON)
		}
	}))
	defer ts.Close()

	sk, err := seeker.New(&config.Seeker{
		FetchDomain:    ts.URL,
		TimeoutSec:     3,
		RetryTimes:     3,
		RetryPeriodSec: 3,
	})
	if err != nil {
		t.Fatalf("new seeker failed:%v", err)
	}

	testCases := []struct {
		description string
		storeID     int
		wantErr     bool
		want        interface{}
	}{
		{
			description: "query by not exist store id",
			storeID:     productNotExistStoreID,
			wantErr:     true,
			want:        nil,
		},
		{
			description: "query data success",
			storeID:     productSuccessStoreID,
			wantErr:     false,
			want: []*model.Product{
				&model.Product{
					ID:              3070177,
					Title:           "Dove Hair Fall Rescue Shampoo",
					Description:     "",
					ImageURL:        "https://assets.honestbee.com/products/images/480/habitattest_hb90000003ea_hb90000003ea-1.jpg",
					PreviewImageURL: "https://assets.honestbee.com/products/images/480/habitattest_hb90000003ea_hb90000003ea-1.jpg",
					Slug:            "",
					Barcodes: []string{
						"HB90000003EA",
					},
					Barcode:              "",
					UnitType:             "unit_type_item",
					SoldBy:               "sold_by_item",
					AmountPerUnit:        "1.0",
					Size:                 "",
					Status:               "status_available",
					ImageURLBasename:     "habitattest_hb90000003ea_hb90000003ea-1.jpg",
					Currency:             "SGD",
					MaxQuantity:          "5183.0",
					CustomerNotesEnabled: false,
					Price:                "1.0",
					NormalPrice:          "3.0",
					Alcohol:              false,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.description), func(t *testing.T) {
			got, _, gotErr := sk.FetchProducts(tc.storeID, 1)
			if tc.wantErr {
				if gotErr == nil {
					t.Errorf("want an error, got nil error")
				}
			} else {
				if gotErr != nil {
					t.Errorf("error want nil, got:%v", gotErr)
				}
				if !reflect.DeepEqual(tc.want, got) {
					t.Errorf("want:%v, got:%v", tc.want, got)
				}
			}
		})
	}
}

func TestFetchBrand(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		brandID, err := strconv.Atoi(strings.TrimPrefix(r.URL.String(), "/api/brands/"))
		if err != nil {
			t.Fatalf("receive an unparsable brand id")
		}
		switch brandID {
		case brandNotExistID:
			w.WriteHeader(http.StatusNotFound)
		case brandSuccessID:
			w.WriteHeader(http.StatusOK)
			w.Write(normalBrandJSON)
		}
	}))
	defer ts.Close()

	sk, err := seeker.New(&config.Seeker{
		FetchDomain:    ts.URL,
		TimeoutSec:     3,
		RetryTimes:     3,
		RetryPeriodSec: 3,
	})
	if err != nil {
		t.Fatalf("new seeker failed:%v", err)
	}

	testCases := []struct {
		description string
		brandID     int
		wantErr     bool
		want        interface{}
	}{
		{
			description: "query not exist brand id",
			brandID:     brandNotExistID,
			wantErr:     true,
			want:        nil,
		},
		{
			description: "query data success",
			brandID:     brandSuccessID,
			wantErr:     false,
			want: &model.Brand{
				ID:                       9483,
				Name:                     "Habitat Test",
				Slug:                     "habitat-test",
				Description:              "If you love food, this is your natural habitat",
				About:                    "Habitat is a next-generation fresh food and grocery store, but more than that, it's a place to get in touch with, get personal with food. If you love food, this is your natural habitat.",
				ImageURL:                 "https://assets-staging.honestbee.com/headers/images/360x120/habitat-test.png",
				ProductsImageURL:         "https://assets-staging.honestbee.com/headers/images/360x120/product-list-habitat-test.png",
				BrandColor:               "",
				Currency:                 "SGD",
				CountryID:                1,
				MinimumOrderFreeDelivery: "",
				DefaultDeliveryFee:       "",
				SameStorePrice:           true,
				BrandType:                "habitat",
				PromotionText:            "",
				ParentBrandID:            0,
				ProductsCount:            142,
				StoreID:                  16849,
				PriceMarkupPercentage:    "0.0",
				FreeDeliveryEligible:     true,
				EstimatedDeliveryTime:    0,
				Closed:                   false,
				CatalogID:                12131,
				OpensAt:                  time.Time{},
				MinimumSpend:             "",
				MinimumSpendExtraFee:     "29.99",
				DefaultConciergeFee:      "",
				ShippingModes:            []string{"offline"},
				DeliveryTypes:            []string{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.description), func(t *testing.T) {
			got, gotErr := sk.FetchBrand(tc.brandID)
			if tc.wantErr {
				if gotErr == nil {
					t.Errorf("want an error, got nil error")
				}
			} else {
				if gotErr != nil {
					t.Errorf("error want nil, got:%v", gotErr)
				}
				if !reflect.DeepEqual(tc.want, got) {
					t.Errorf("want:%v, got:%v", tc.want, got)
				}
			}
		})
	}
}

func TestFetchStores(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exKey := r.FormValue("externalKey")
		switch exKey {
		case storeNotExistExternalKey:
			w.WriteHeader(http.StatusNotFound)
		case storeCauseDecodeFailedExternalKey:
			w.WriteHeader(http.StatusOK)
			w.Write(causeDecodeFailedStoreJSON)
		case storeSuccessExternalKey:
			w.WriteHeader(http.StatusOK)
			w.Write(normalStoreJSON)
		}
	}))
	defer ts.Close()

	sk, err := seeker.New(&config.Seeker{
		FetchDomain:    ts.URL,
		TimeoutSec:     3,
		RetryTimes:     3,
		RetryPeriodSec: 3,
	})
	if err != nil {
		t.Fatalf("new seeker failed:%v", err)
	}

	testCases := []struct {
		description string
		externalKey string
		wantErr     bool
		want        interface{}
	}{
		{
			description: "query not exist external key",
			externalKey: storeNotExistExternalKey,
			wantErr:     true,
			want:        nil,
		},
		{
			description: "query data but json decode failed",
			externalKey: storeCauseDecodeFailedExternalKey,
			wantErr:     true,
			want:        nil,
		},
		{
			description: "query data success",
			externalKey: storeSuccessExternalKey,
			wantErr:     false,
			want: []*model.Store{
				&model.Store{
					ID:                       16849,
					Name:                     "Habitat Store",
					PickupPoint:              "",
					Slug:                     "habitat-store",
					BrandID:                  9483,
					AddressID:                36257,
					CatalogID:                14202,
					Priority:                 "",
					Notes:                    "",
					Description:              "",
					ImageURL:                 "",
					Closed:                   false,
					TemporarilyClosed:        false,
					OpensAt:                  time.Time{},
					EstimatedDeliveryTime:    0,
					BufferTime:               0,
					StoreType:                "habitat",
					MinimumOrderFreeDelivery: "",
					DefaultDeliveryFee:       "",
					FreeDeliveryEligible:     true,
					MinimumSpend:             "",
					MinimumSpendExtraFee:     "29.99",
					ShippingModes:            []string{"offline"},
					DeliveryTypes:            []string{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.description), func(t *testing.T) {
			got, gotErr := sk.FetchStores(tc.externalKey)
			if tc.wantErr {
				if gotErr == nil {
					t.Errorf("want an error, got nil error")
				}
			} else {
				if gotErr != nil {
					t.Errorf("error want nil, got:%v", gotErr)
				}
				if !reflect.DeepEqual(tc.want, got) {
					t.Errorf("want:%v, got:%v", tc.want, got)
				}
			}
		})
	}
}
