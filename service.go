package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
)

const DefaultInfluxURL = "http://influx:8086"
const DefaultInfluxUsername = "admin"
const DefaultInfluxPassword = "admin"
const DefaultInfluxDatabase = "telemetry"
const DefaultInfluxPrecision = "ns"
const DefaultRunOnce = false

type Service struct {
	processor client.Client
	database  string
	precision string
}

type Point struct {
	Tags   map[string]string      `json:"tags"`
	Fields map[string]interface{} `json:"fields"`
	Point  string                 `json:"point"`
}

func (svc *Service) Init() error {
	var err error
	var influxUrl string
	var influxUsername string
	var influxPassword string
	var runOnce bool

	// Check for the INFLUX_URL
	if os.Getenv("INFLUX_URL") == "" {
		influxUrl = DefaultInfluxURL
	} else {
		influxUrl = os.Getenv("INFLUX_URL")
	}

	// Check for the INFLUX_USERNAME
	if os.Getenv("INFLUX_USERNAME") == "" {
		influxUsername = DefaultInfluxUsername
	} else {
		influxUsername = os.Getenv("INFLUX_USERNAME")
	}

	// Check for the INFLUX_PASSWORD
	if os.Getenv("INFLUX_PASSWORD") == "" {
		influxPassword = DefaultInfluxPassword
	} else {
		influxPassword = os.Getenv("INFLUX_PASSWORD")
	}

	// Check for the INFLUX_DATABASE
	if os.Getenv("INFLUX_DATABASE") == "" {
		svc.database = DefaultInfluxDatabase
	} else {
		svc.database = os.Getenv("INFLUX_DATABASE")
	}

	// Check for the INFLUX_PRECISION
	if os.Getenv("INFLUX_PRECISION") == "" {
		svc.precision = DefaultInfluxPrecision
	} else {
		svc.precision = os.Getenv("INFLUX_PRECISION")
	}

	// Check for the RUN_ONCE
	if os.Getenv("RUN_ONCE") == "" {
		runOnce = DefaultRunOnce
	} else {
		runOnce, err = strconv.ParseBool(os.Getenv("RUN_ONCE"))
		if err != nil {
			return err
		}
	}

	// Setup the influx client.
	processor, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxUrl,
		Username: influxUsername,
		Password: influxPassword,
	})
	if err != nil {
		return err
	}
	defer processor.Close()

	// Setup the service processor.
	svc.processor = processor

	// Setup http server.
	mux := http.NewServeMux()

	// Setup the http handle for incoming telemetry.
	mux.HandleFunc("/", svc.handler)
	mux.HandleFunc("/ping", svc.ping)

	// Run once mode for unit tests.
	if runOnce {
		return nil
	}

	return http.ListenAndServe(":3000", mux)
}

// Private

func (svc *Service) ping(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
}

func (svc *Service) handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)

		var point Point
		err := decoder.Decode(&point)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  svc.database,
			Precision: svc.precision,
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pt, err := client.NewPoint(point.Point, point.Tags, point.Fields, time.Now())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		bp.AddPoint(pt)

		if err := svc.processor.Write(bp); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := svc.processor.Close(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
