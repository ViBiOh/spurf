# spurf

[![Build](https://github.com/ViBiOh/spurf/workflows/Build/badge.svg)](https://github.com/ViBiOh/spurf/actions)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_spurf&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_spurf)

## Usage

```bash
Usage of spurf:
  -dbHost string
        [db] Host {SPURF_DB_HOST}
  -dbMaxConn uint
        [db] Max Open Connections {SPURF_DB_MAX_CONN} (default 5)
  -dbName string
        [db] Name {SPURF_DB_NAME}
  -dbPass string
        [db] Pass {SPURF_DB_PASS}
  -dbPort uint
        [db] Port {SPURF_DB_PORT} (default 5432)
  -dbSslmode string
        [db] SSL Mode {SPURF_DB_SSLMODE} (default "disable")
  -dbTimeout uint
        [db] Connect timeout {SPURF_DB_TIMEOUT} (default 10)
  -dbUser string
        [db] User {SPURF_DB_USER}
  -enedisFile string
        [enedis] CSV export to load {SPURF_ENEDIS_FILE}
  -enedisName string
        [enedis] Name {SPURF_ENEDIS_NAME} (default "home")
  -loggerJson
        [logger] Log format as JSON {SPURF_LOGGER_JSON}
  -loggerLevel string
        [logger] Logger level {SPURF_LOGGER_LEVEL} (default "INFO")
  -loggerLevelKey string
        [logger] Key for level in JSON {SPURF_LOGGER_LEVEL_KEY} (default "level")
  -loggerMessageKey string
        [logger] Key for message in JSON {SPURF_LOGGER_MESSAGE_KEY} (default "message")
  -loggerTimeKey string
        [logger] Key for timestamp in JSON {SPURF_LOGGER_TIME_KEY} (default "time")
```
