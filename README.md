# spurf

[![Build](https://github.com/ViBiOh/spurf/workflows/Build/badge.svg)](https://github.com/ViBiOh/spurf/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/spurf)](https://goreportcard.com/report/github.com/ViBiOh/spurf)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_spurf&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_spurf)

## Usage

```bash
Usage of spurf:
  -datahubAccessToken string
        [datahub] Access Token {SPURF_DATAHUB_ACCESS_TOKEN}
  -datahubClientID string
        [datahub] Client ID {SPURF_DATAHUB_CLIENT_ID}
  -datahubClientSecret string
        [datahub] Client Secret {SPURF_DATAHUB_CLIENT_SECRET}
  -datahubRedirectUri string
        [datahub] Redirect URI {SPURF_DATAHUB_REDIRECT_URI} (default "https://api.vibioh.fr/dump/")
  -datahubRefreshToken string
        [datahub] Refresh Token {SPURF_DATAHUB_REFRESH_TOKEN}
  -datahubSandbox
        [datahub] Sandbox mode {SPURF_DATAHUB_SANDBOX}
  -datahubUsagePointId string
        [datahub] Identifiant du point de livraison {SPURF_DATAHUB_USAGE_POINT_ID}
  -dbHost string
        [db] Host {SPURF_DB_HOST}
  -dbName string
        [db] Name {SPURF_DB_NAME}
  -dbPass string
        [db] Pass {SPURF_DB_PASS}
  -dbPort uint
        [db] Port {SPURF_DB_PORT} (default 5432)
  -dbSslmode string
        [db] SSL Mode {SPURF_DB_SSLMODE} (default "disable")
  -dbUser string
        [db] User {SPURF_DB_USER}
  -enedisName string
        [enedis] Name {SPURF_ENEDIS_NAME} (default "home")
  -enedisTimezone string
        [enedis] Timezone {SPURF_ENEDIS_TIMEZONE} (default "Europe/Paris")
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
