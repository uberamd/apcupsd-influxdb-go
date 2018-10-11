apcupsd-influxdb-go
================

What does it do?
----------------

Reads apcupsd statistics and sends them to an InfluxDB instance in 10 second intervals. 

Usage
-----

Available flags:

```cgo
$ ./apcupsd-influxdb-go
Usage of ./apcupsd-influxdb-go:
  -apcupsd-addr string
    	address of apcupsd Network Information Server (NIS) (default ":3551")
  -apcupsd-network string
    	network of apcupsd Network Information Server (NIS): typically "tcp", "tcp4", or "tcp6" (default "tcp")
  -http-port string
    	port for the http server to listen on for health checks (default "8084")
  -influxdb-addr string
    	address of influxdb endpoint, ex: http://127.0.0.1:8086 (default "http://127.0.0.1:8086")
  -influxdb-database string
    	influxdb database to store datapoints (default "homelab_custom")
  -influxdb-pass string
    	password for influxdb access (default "admin")
  -influxdb-user string
    	username for influxdb access (default "admin")

```

To Do
------

- Add real error handling
- Make the interval customizable