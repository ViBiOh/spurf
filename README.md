# spurf

[![Build Status](https://travis-ci.org/ViBiOh/spurf.svg?branch=master)](https://travis-ci.org/ViBiOh/spurf)
[![Go Report Card](https://goreportcard.com/badge/github.com/ViBiOh/spurf)](https://goreportcard.com/report/github.com/ViBiOh/spurf)

## Usage

```bash
Usage of spurf:
  -c    Healthcheck (check and exit)
  -dbHost string
        [db] Host
  -dbName string
        [db] Name
  -dbPass string
        [db] Pass
  -dbPort string
        [db] Port (default "5432")
  -dbUser string
        [db] User
  -enedisEmail string
        [enedis]  Email
  -enedisPassword string
        [enedis]  Password
  -enedisTimezone string
        [enedis] Timezone (default "Europe/Paris")
  -schedulerHour int
        [scheduler] Hour of running (default 8)
  -schedulerInterval string
        [scheduler] Duration between two runs (default "24h")
  -schedulerMaxRetry int
        [scheduler] Max retry (default 10)
  -schedulerMinute int
        [scheduler] Minute of running
  -schedulerRetry string
        [scheduler] Duration between two retries (default "10m")
  -schedulerTimezone string
        [scheduler] Timezone of running (default "Europe/Paris")
  -tracingAgent string
        [tracing] Jaeger Agent (e.g. host:port) (default "jaeger:6831")
  -tracingName string
        [tracing] Service name
```
