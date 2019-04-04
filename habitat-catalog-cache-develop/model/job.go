package model

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

const (
	jobSetsKey = "catalog_cache_service_sync_jobs"
)

type jobService interface {
	PushJob(*Job) error
	PopJob() (*Job, error)
}

type jobOps struct {
	getConnTimeout time.Duration
	pool           *redis.Pool
}

func (j *jobOps) do(cmd string, args ...interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), j.getConnTimeout)
	defer cancel()

	conn, err := j.pool.GetContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [do] get redis connection failed")
	}
	defer conn.Close()

	reply, err := conn.Do(cmd, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [do] failed on cmd:%s, args:%v", cmd, args)
	}

	return reply, nil
}

func (j *jobOps) PushJob(job *Job) error {
	dj, err := json.Marshal(job)
	if err != nil {
		return errors.Wrapf(err, "model: [PushJob] json marshal failed")
	}

	_, err = j.do("SADD", jobSetsKey, base64.StdEncoding.EncodeToString(dj))
	return errors.Wrapf(err, "model: [PushJob] SADD failed")
}

var (
	// ErrNoJob means there is no job in queue.
	ErrNoJob = errors.New("no job")
)

func (j *jobOps) PopJob() (*Job, error) {
	jobs, err := redis.Strings(j.do("SPOP", jobSetsKey, 1))
	if err != nil {
		return nil, errors.Wrapf(err, "model: [PopJob] SPOP failed")
	}

	job := ""
	switch len(jobs) {
	case 0:
		return nil, ErrNoJob
	case 1:
		job = jobs[0]
	default:
		return nil, errors.Errorf("model: [PopJob] SPOP 1 went wrong, got:%v", jobs)
	}

	dj, err := base64.StdEncoding.DecodeString(job)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [PopJob] base64 decode failed")
	}

	ret := new(Job)
	if err = json.Unmarshal(dj, ret); err != nil {
		return nil, errors.Wrapf(err, "model: [PopJob] json unmarshal failed")
	}

	return ret, nil
}

// JobType is a job's type
type JobType int

const (
	// ExternalKeyJob means using an external key to update.
	ExternalKeyJob = JobType(iota)
	// StoreIDJob means sync up a store's information and belong to it's products.
	StoreIDJob
	// BrandIDJob means sync up a brand's information.
	BrandIDJob
)

// Job is the structure of a sync up job in redis.
type Job struct {
	Type  JobType     `json:"type"`
	Value interface{} `json:"value"`
}

// GetValue will try to assign the same type as input dest variable from Job.Value
func (j *Job) GetValue(dest interface{}) (err error) {
	switch dest.(type) {
	case *string:
		if source, ok := j.Value.(string); ok {
			*(dest.(*string)) = source
		} else {
			err = errors.Errorf("model: [GetValue] job value cast to string failed")
		}
	case *int:
		if source, ok := j.Value.(int); ok {
			*(dest.(*int)) = source
		} else if source, ok := j.Value.(float64); ok {
			*(dest.(*int)) = int(source)
		} else {
			err = errors.Errorf("model: [GetValue] job value cast to int failed")
		}
	case *StoreIDJobValue:
		if source, ok := j.Value.(map[string]interface{}); ok {
			if id, ok := source["id"].(int); ok {
				dest.(*StoreIDJobValue).ID = id
			}
			if exkey, ok := source["externalKey"].(string); ok {
				dest.(*StoreIDJobValue).ExternalKey = exkey
			}
		} else if source, ok := j.Value.(*StoreIDJobValue); ok {
			dest.(*StoreIDJobValue).ID = source.ID
			dest.(*StoreIDJobValue).ExternalKey = source.ExternalKey
		} else {
			err = errors.Errorf("model: [GetValue] job value cast to map[string]interface{} failed")
		}
	default:
		err = errors.Errorf("model: [GetValue] unsupport cast type:%T", dest)
	}
	return err
}

// StoreIDJobValue is a struct for StoreIDJob value.
type StoreIDJobValue struct {
	ID          int    `json:"id"`
	ExternalKey string `json:"externalKey"`
}
