package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// SyncProducts handles sync products request.
func SyncProducts(e *Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
	exkeys := strings.Split(ValueFetcher("externalKeys", r.PostFormValue), ",")
	for _, exkey := range exkeys {
		strings.TrimSpace(exkey)
		if exkey == "" {
			continue
		}
		err := e.Service.PushJob(&model.Job{
			Type:  model.ExternalKeyJob,
			Value: exkey,
		})
		if err != nil {
			return nil, NewErr(ServerInternalErrCode, errors.Wrapf(
				err,
				"handle: [SyncProducts] PushJob failed on external key:%s, keys:%v",
				exkey, exkeys,
			))
		}
	}

	return fmt.Sprintf("success triggered sync up job external keys:%v", exkeys), nil
}

// SyncStore handles sync store request.
func SyncStore(e *Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
	exkey := ValueFetcher("externalKey", r.PostFormValue)
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil || exkey == "" {
		return nil, NewErr(InvalidAttributeErrCode, errors.Wrapf(
			err,
			"handle: [SyncStore] convert store id:%s from string to int failed or externalKey:%s is empty",
			ps.ByName("id"), exkey,
		))
	}

	err = e.Service.PushJob(&model.Job{
		Type: model.StoreIDJob,
		Value: &model.StoreIDJobValue{
			ExternalKey: exkey,
			ID:          id,
		},
	})
	if err != nil {
		return nil, NewErr(ServerInternalErrCode, errors.Wrapf(
			err,
			"handle: [SyncStore] PushJob failed on external key:%s, id:%d",
			exkey, id,
		))
	}

	return fmt.Sprintf("success triggered sync up job store id:%d, external key:%s", id, exkey), nil
}

// SyncBrand handles sync brand request.
func SyncBrand(e *Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		return nil, NewErr(InvalidAttributeErrCode, errors.Wrapf(err, "handler: [SyncBrand] convert brand id failed"))
	}

	err = e.Service.PushJob(&model.Job{
		Type:  model.BrandIDJob,
		Value: id,
	})
	if err != nil {
		return nil, NewErr(ServerInternalErrCode, errors.Wrapf(
			err,
			"handle: [SyncBrand] PushJob failed on id:%d",
			id,
		))
	}

	return fmt.Sprintf("success triggered sync up job brand id:%d", id), nil
}
