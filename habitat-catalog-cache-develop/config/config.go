package config

import (
	"flag"
	"os"
)

var (
	// Version is the app version number at build time
	Version = "No Version Provided"
)

// SyncWorker is the sync up worker configuration.
type SyncWorker struct {
	Database                 *Database
	Cache                    *Cache
	Seeker                   *Seeker
	Procesor                 *Procesor
	TCPHealthCheckListenAddr string
	SyncCheckPeriodSec       int
}

func newSyncWorker() *SyncWorker {
	fset := flag.NewFlagSet("sync worker", flag.ContinueOnError)
	c := &SyncWorker{
		Database: new(Database),
		Seeker:   new(Seeker),
		Procesor: new(Procesor),
		Cache:    new(Cache),
	}
	fset.IntVar(&c.SyncCheckPeriodSec, "sync_check_period_sec", 30, "checking period of sync up job")
	fset.StringVar(&c.TCPHealthCheckListenAddr, "tcp_health_check_listen_addr", ":9090", "the TCP health check listening address")
	setFlagConfig(fset, c.Database, c.Procesor, c.Seeker, c.Cache)
	return c
}

// CatalogService is the catalog http api service configuration.
type CatalogService struct {
	Database  *Database
	Cache     *Cache
	HTTP      *HTTP
	AutoExec  *AutoExec
	BasicAuth *BasicAuth
}

func newCatalogService() *CatalogService {
	fset := flag.NewFlagSet("catalog service", flag.ContinueOnError)
	c := &CatalogService{
		Database:  new(Database),
		HTTP:      new(HTTP),
		Cache:     new(Cache),
		AutoExec:  new(AutoExec),
		BasicAuth: new(BasicAuth),
	}
	setFlagConfig(fset, c.Database, c.HTTP, c.Cache, c.AutoExec, c.BasicAuth)
	return c
}

type flagConfig interface {
	setFlags(*flag.FlagSet)
}

// AutoExec is the autoexec configuration.
type AutoExec struct {
	SyncupPeriodSec int
}

func (a *AutoExec) setFlags(flagSet *flag.FlagSet) {
	flagSet.IntVar(&a.SyncupPeriodSec, "autoexec_syncup_period_sec", 21600, "time period of a syncup job will push into job queue")
}

// HTTP is the server http configuration.
type HTTP struct {
	ReadTimeoutSec  int
	WriteTimeoutSec int
	IdleTimeoutSec  int
	ListenAddr      string
}

func (h *HTTP) setFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&h.ListenAddr, "http_listen_addr", ":8080", "the HTTP server listening address")
	flagSet.IntVar(&h.ReadTimeoutSec, "http_read_timeout_sec", 10, "the HTTP server read timeout seconds")
	flagSet.IntVar(&h.WriteTimeoutSec, "http_write_timeout_sec", 30, "the HTTP server write timeout seconds")
	flagSet.IntVar(&h.IdleTimeoutSec, "http_idle_timeout_sec", 360, "the HTTP server idle timeout seconds")
}

// BasicAuth is the basic authorization for catalog cache http request.
type BasicAuth struct {
	User string
	Pwd  string
}

func (b *BasicAuth) setFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&b.User, "basic_auth_user", "admin", "basic auth user")
	flagSet.StringVar(&b.Pwd, "basic_auth_pwd", "33456783345678", "basic auth password")
}

// Database is the postgres configuration.
type Database struct {
	Host                     string
	Port                     string
	User                     string
	Pwd                      string
	Name                     string
	ConnectTimeoutSec        int
	MaxIdle                  int
	MaxActive                int
	ReadTimeoutSec           int
	WriteTimeoutSec          int
	TransactionMaxTimeoutSec int
}

func (d *Database) setFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&d.Name, "db_name", "testing", "the sync up database name")
	flagSet.StringVar(&d.Pwd, "db_pwd", "", "the sync up database password")
	flagSet.StringVar(&d.Host, "db_host", "localhost", "the sync up database host")
	flagSet.StringVar(&d.Port, "db_port", "5432", "the sync up database port")
	flagSet.StringVar(&d.User, "db_user", "root", "the sync up database user")
	flagSet.IntVar(&d.MaxIdle, "db_max_idle", 500, "the sync up database max idle")
	flagSet.IntVar(&d.MaxActive, "db_max_active", 1000, "the sync up database max active")
	flagSet.IntVar(&d.ConnectTimeoutSec, "db_connect_timeout_sec", 5, "the sync up database connect timeout second")
	flagSet.IntVar(&d.ReadTimeoutSec, "db_read_timeout_sec", 10, "the sync up database read timeout second")
	flagSet.IntVar(&d.WriteTimeoutSec, "db_write_timeout_sec", 15, "the sync up database write timeout second")
	flagSet.IntVar(&d.TransactionMaxTimeoutSec, "db_transaction_max_timeout_sec", 60, "the sync up databasse transaction max timeout second")
}

// Cache is the redis configuration.
type Cache struct {
	MaxIdle                 int
	MaxActive               int
	IdleTimeoutSec          int
	Wait                    bool
	ConnectTimeoutSec       int
	ReadTimeoutSec          int
	WriteTimeoutSec         int
	GetConnectionTimeoutSec int
	Host                    string
	Port                    string
	Pwd                     string
	Index                   int
}

func (c *Cache) setFlags(flagSet *flag.FlagSet) {
	flagSet.IntVar(&c.MaxIdle, "cache_max_idle", 500, "cache max idle")
	flagSet.IntVar(&c.MaxActive, "cache_max_active", 1000, "cache max active")
	flagSet.IntVar(&c.GetConnectionTimeoutSec, "cache_get_connection_timeout_sec", 3, "redis pool get connection timeout seconds")
	flagSet.IntVar(&c.IdleTimeoutSec, "cache_idle_timeout_sec", 1200, "close connections after remaining idle for this duration")
	flagSet.BoolVar(&c.Wait, "cache_wait", false, "if true and the pool is at the MaxActive limit then Get() waits for a connection to be returned to the pool before returning")
	flagSet.IntVar(&c.ConnectTimeoutSec, "cache_connect_timeout_sec", 5, "cache connect timeout second")
	flagSet.IntVar(&c.ReadTimeoutSec, "cache_read_timeout_sec", 10, "cache read timeout second")
	flagSet.IntVar(&c.WriteTimeoutSec, "cache_write_timeout_sec", 15, "cache write timeout second")
	flagSet.StringVar(&c.Host, "cache_host", "127.0.0.1", "cache host")
	flagSet.StringVar(&c.Port, "cache_port", "6379", "cache port")
	flagSet.StringVar(&c.Pwd, "cache_pwd", "", "cache password")
	flagSet.IntVar(&c.Index, "cache_db_index", 1, "cache db index")
}

// Seeker is the seeker package configuration.
type Seeker struct {
	FetchDomain    string
	TimeoutSec     int
	RetryTimes     int
	RetryPeriodSec int
}

func (s *Seeker) setFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&s.FetchDomain, "seeker_base_url", "localhost", "where to seek the information")
	flagSet.IntVar(&s.TimeoutSec, "seeker_timeout_sec", 10, "http client timeout seconds")
	flagSet.IntVar(&s.RetryTimes, "seeker_retry_times", 3, "http failed retry times")
	flagSet.IntVar(&s.RetryPeriodSec, "seeker_retry_period_sec", 5, "http retry period seconds")
}

// Procesor is the procesor package configuration.
type Procesor struct {
	WorkerNum int
	PoolSize  int
}

func (p *Procesor) setFlags(flagSet *flag.FlagSet) {
	flagSet.IntVar(&p.PoolSize, "procesor_pool_size", 50, "procesor worker pool size")
	flagSet.IntVar(&p.WorkerNum, "procesor_worker_num", 10, "procesor the numbers of worker")
}

func setFlagConfig(flagSet *flag.FlagSet, configs ...flagConfig) {
	for _, config := range configs {
		config.setFlags(flagSet)
	}
	flagSet.Parse(os.Args[1:])
}

// NewSyncWorker returns a SyncWorker configuration.
func NewSyncWorker() *SyncWorker {
	return newSyncWorker()
}

// NewCatalogService returns a CatalogService configuration.
func NewCatalogService() *CatalogService {
	return newCatalogService()
}
