# spurf

[![Build Status](https://travis-ci.org/ViBiOh/spurf.svg?branch=master)](https://travis-ci.org/ViBiOh/spurf)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/spurf)](https://goreportcard.com/report/github.com/ViBiOh/spurf)

## Usage

```bash
Usage of spurf:
  -c    Healthcheck (check and exit)
  -dbHost string
        [db] Host {SPURF_DB_HOST}
  -dbName string
        [db] Name {SPURF_DB_NAME}
  -dbPass string
        [db] Pass {SPURF_DB_PASS}
  -dbPort string
        [db] Port {SPURF_DB_PORT} (default "5432")
  -dbUser string
        [db] User {SPURF_DB_USER}
  -enedisEmail string
        [enedis] Email {SPURF_ENEDIS_EMAIL}
  -enedisPassword string
        [enedis] Password {SPURF_ENEDIS_PASSWORD}
  -enedisTimezone string
        [enedis] Timezone {SPURF_ENEDIS_TIMEZONE} (default "Europe/Paris")
  -schedulerHour int
        [scheduler] Hour of running {SPURF_SCHEDULER_HOUR} (default 8)
  -schedulerInterval string
        [scheduler] Duration between two runs {SPURF_SCHEDULER_INTERVAL} (default "24h")
  -schedulerMaxRetry int
        [scheduler] Max retry {SPURF_SCHEDULER_MAX_RETRY} (default 10)
  -schedulerMinute int
        [scheduler] Minute of running {SPURF_SCHEDULER_MINUTE}
  -schedulerOnStart
        [scheduler] Start scheduler on start {SPURF_SCHEDULER_ON_START}
  -schedulerRetry string
        [scheduler] Duration between two retries {SPURF_SCHEDULER_RETRY} (default "10m")
  -schedulerTimezone string
        [scheduler] Timezone of running {SPURF_SCHEDULER_TIMEZONE} (default "Europe/Paris")
  -tracingAgent string
        [tracing] Jaeger Agent (e.g. host:port) {SPURF_TRACING_AGENT} (default "jaeger:6831")
  -tracingName string
        [tracing] Service name {SPURF_TRACING_NAME}
```
