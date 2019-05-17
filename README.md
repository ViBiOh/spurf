# spurf

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
  -schedulerMinute int
        [scheduler] Minute of running
  -schedulerRetry string
        [scheduler] Duration between two retries (default "10m")
  -schedulerTimezone string
        [scheduler] Timezone of running (default "Europe/Paris")
```