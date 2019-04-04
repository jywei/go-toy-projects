# Habitat Catalog Cache service

habitat catalog cache service for worker app and PLC controller to query catalog informations.

[doc](https://honestbee.atlassian.net/wiki/spaces/EN/pages/534381673/Catalog+caching+service+design)

## Preparations

### Install Golang
please follow the instruction here for [mac](https://golang.org/doc/install)
and install the latest version.

### Checking Golang Version
```bash
go version >= go1.10 (go test multi package converage issue, fixed on go1.10)
```

### Install Dependencies
using [dep](https://github.com/golang/dep), after installation, just type
```bash
dep ensure
```

### Setup Config Variables
| variable name                    | default value    | description                                                                                                                  |
| -------------------------------- | ---------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| check_perioud_sec                | 30               | checking period of sync up job                                                                                               |
| http_listen_addr                 | "localhost:8080" | the HTTP server listening address                                                                                            |
| http_read_timeout_sec            | 10               | the HTTP server read timeout seconds                                                                                         |
| http_write_timeout_sec           | 30               | the HTTP server write timeout seconds                                                                                        |
| http_idle_timeout_sec            | 360              | the HTTP server idle timeout seconds                                                                                         |
| db_name                          | "testing"        | the sync up database name                                                                                                    |
| db_pwd                           | ""               | the sync up database password                                                                                                |
| db_host                          | "localhost"      | the sync up database host                                                                                                    |
| db_port                          | "5432"           | the sync up database port                                                                                                    |
| db_user                          | "root"           | the sync up database user                                                                                                    |
| db_max_idle                      | 500              | the sync up database max idle                                                                                                |
| db_max_active                    | 1000             | the sync up database max active                                                                                              |
| db_connect_timeout_sec           | 5                | the sync up database connect timeout second                                                                                  |
| db_read_timeout_sec              | 10               | the sync up database read timeout second                                                                                     |
| db_write_timeout_sec             | 15               | the sync up database write timeout second                                                                                    |
| cache_max_idle                   | 500              | cache max idle                                                                                                               |
| cache_max_active                 | 1000             | cache max active                                                                                                             |
| cache_get_connection_timeout_sec | 3                | redis pool get connection timeout seconds                                                                                    |
| cache_idle_timeout_sec           | 1200             | close connections after remaining idle for this duration                                                                     |
| cache_wait                       | false            | if true and the pool is at the MaxActive limit then Get() waits for a connection to be returned to the pool before returning |
| cache_connect_timeout_sec        | 5                | cache connect timeout second                                                                                                 |
| cache_read_timeout_sec           | 10               | cache read timeout second                                                                                                    |  |
| cache_write_timeout_sec          | 15               | cache write timeout second                                                                                                   |
| cache_host                       | "127.0.0.1"      | cache host                                                                                                                   |
| cache_port                       | "6379"           | cache port                                                                                                                   |
| cache_password                   | ""               | cache password                                                                                                               |
| cache_db_index                   | 1                | cache db index                                                                                                               |
| seeker_base_url                  | "localhost"      | where to seek the information                                                                                                |
| seeker_timeout_sec               | 10               | http client timeout seconds                                                                                                  |
| seeker_retry_times               | 3                | http failed retry times                                                                                                      |  |
| seeker_retry_period_sec          | 5                | http retry period seconds                                                                                                    |
| procesor_pool_size               | 50               | procesor worker pool size                                                                                                    |
| procesor_worker_num              | 10               | procesor the numbers of worker                                                                                               |

### Install Database
using Postgres database, please follow [the instruction](https://www.postgresql.org/download/)

### Database Migration
using [goose](https://bitbucket.org/liamstask/goose)

#### Install Goose
```bash
go get bitbucket.org/liamstask/goose/cmd/goose
```

#### Setup Database Env
setup ENV for goose to connect to the database

develop usage:

using database name **zen**

```bash
export DB_USER={the user name}
export DB_PASSWORD={the user password}
```

ci usage:

```bash
export DB_USER=
export DB_PASSWORD=
export CATALOG_CACHE_DATABASE_URI=
export CATALOG_CACHE_DATABASE_NAME=
```

#### Goose version (Print the current version of the database)
```bash
goose dbversion

goose: dbversion 20180301103424
```

#### Goose up (Apply all available migrations)
```bash
goose -env=ci up

goose: migrating db environment 'development', current version: 0, target: 20180301103754
OK    20180301100347_addCategories.sql
OK    20180301103424_addSections.sql
OK    20180301103754_addArticles.sql
```

#### Goose down (Roll back a single migration from the current version)
```bash
goose down

goose: migrating db environment 'development', current version: 20180301103754, target: 20180301103424
OK    20180301103754_addArticles.sql
...
2018/03/02 10:43:08 no previous version found
```

#### Goose status (Dump the migration status for the current DB)
```bash
goose status

goose: status for environment 'development'
    Applied At                  Migration
    =======================================
    Fri Mar  2 02:44:31 2018 -- 20180301100347_addCategories.sql
    Fri Mar  2 02:44:31 2018 -- 20180301103424_addSections.sql
    Fri Mar  2 02:44:31 2018 -- 20180301103754_addArticles.sql
```

### Testing
using TDD testing, so simply type in project root folder
```bash
go test ./... -race -v -cover
```

## DevOps

### Build Time Setup App Version
```bash
go build -ldflags "-X github.com/honestbee/habitat-catalog-cache/config.Version={$APP_VERSION}"
```

### Check Server Status
```bash
curl host:port/api/v1/status
{"go-version":"go1.10.3","app-version":"No Version Provided","server-time":"2018-06-20 07:42:37.146032785 +0000 UTC"}
```
