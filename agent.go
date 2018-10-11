package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/mdlayher/apcupsd"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	hostname, _ = os.Hostname()
)

type BaseConfig struct {
	ApcupsdAddr string
	ApcupsdNetwork string
	HealthHttpPort string
	InfluxDatabase string
	InfluxAddr string
	InfluxUser string
	InfluxPass string
}

func HealthHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OK"))
}

func readUpsData(bc *BaseConfig) {
	for {
		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr: bc.InfluxAddr,
			Username: bc.InfluxUser,
			Password: bc.InfluxPass,
		})

		if err != nil {
			log.Fatal(err)
		}

		defer c.Close()

		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database: bc.InfluxDatabase,
			Precision: "s",
		})

		if err != nil {
			log.Fatal(err)
		}

		apcClient, err := apcupsd.Dial(bc.ApcupsdNetwork, bc.ApcupsdAddr)

		if err != nil {
			log.Fatal(err)
		}
		status, err := apcClient.Status()

		if err != nil {
			log.Fatal(err)
		}

		tags := map[string]string{
			"ups": bc.ApcupsdAddr,
			"host": hostname,
		}

		fields := map[string]interface{}{
			"lineVoltage": status.LineVoltage,
			"loadPercent": status.LoadPercent,
			"batteryChargePercent": status.BatteryChargePercent,
			"timeLeft": int64(status.TimeLeft / time.Second),
			"maximumTime": int64(status.MaximumTime / time.Second),
			"lowTransferVoltage": status.LowTransferVoltage,
			"highTransferVoltage": status.HighTransferVoltage,
			"batteryVoltage": status.BatteryVoltage,
			"numberTransfers": status.NumberTransfers,
			"timeOnBattery": int64(status.TimeOnBattery / time.Second),
			"cumulativeTimeOnBattery": int64(status.CumulativeTimeOnBattery / time.Second),
			"nominalInputVoltage": status.NominalInputVoltage,
			"nominalBatteryVoltage": status.NominalBatteryVoltage,
			"nominalPower": status.NominalPower,
		}

		pt, err := client.NewPoint("apcupsd_readings", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		bp.AddPoint(pt)
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}

		// Print temperature and humidity
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}

		log.Print("Points submitted to influxdb...")
		time.Sleep(10 * time.Second)
	}
}

func main() {
	bc := new(BaseConfig)
	flag.StringVar(&bc.ApcupsdAddr,"apcupsd-addr", ":3551", "address of apcupsd Network Information Server (NIS)")
	flag.StringVar(&bc.ApcupsdNetwork,"apcupsd-network", "tcp", `network of apcupsd Network Information Server (NIS): typically "tcp", "tcp4", or "tcp6"`)
	flag.StringVar(&bc.HealthHttpPort,"http-port", "8084", "port for the http server to listen on for health checks")
	flag.StringVar(&bc.InfluxDatabase,"influxdb-database", "homelab_custom", "influxdb database to store datapoints")
	flag.StringVar(&bc.InfluxAddr,"influxdb-addr", "http://127.0.0.1:8086", "address of influxdb endpoint, ex: http://127.0.0.1:8086")
	flag.StringVar(&bc.InfluxUser,"influxdb-user", "admin", "username for influxdb access")
	flag.StringVar(&bc.InfluxPass,"influxdb-pass", "admin", "password for influxdb access")

	flag.Parse()

	go readUpsData(bc)

	r := mux.NewRouter()
	r.HandleFunc("/healthz", HealthHandler)
	log.Fatal(http.ListenAndServe(":8084", r))
	log.Print("Listening on :8084")
}
