# Telemetry

Simple lightweight API that accepts a POST request with a JSON payload that is
written directly into InfluxDB.

## Usage

This task is configured by the following environment variables:

```bash
INFLUX_URL       # InfluxDB URL
INFLUX_USERNAME  # username to connect to InfluxDB
INFLUX_PASSWORD  # password for given user
INFLUX_DATABASE  # Influx database to store telemetry
INFLUX_PRECISION # precision of data point
RUN_ONCE         # used for unit testing to not start the http server
```

## Example

```
curl -XPOST http://127.0.0.1:3000/ -d '{"tags": {"cpu": "cpu-total"}, "fields": {"idle": 10.1,"system": 53.3,"user": 46.6}, "point": "cpu_usage"}}'
```

![alt text](/images/screenshot.png)

## Contributing

* `make run` - runs the api in a docker container
* `make build` - builds your telemetry docker container
* `make vet` - go fmt and vet code
* `make test` - run unit tests

Before you submit a pull request please update the semantic version inside of
`main.go` with what you feel is appropriate and then edit the `CHANGELOG.md` with
your changes and follow a similar structure to what is there.
