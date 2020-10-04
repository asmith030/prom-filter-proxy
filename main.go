package main

import (
	"flag"
	"os"
	"time"

	"github.com/asmith030/prom-filter-proxy/server"
	"github.com/sirupsen/logrus"
)

var (
	endpoint   = flag.String("endpoint", "http://localhost/metrics", "Address of the metrics endpoint to proxy")
	listenPort = flag.Int("listenPort", 8080, "TCP port that we should listen on")
	path       = flag.String("path", "/metrics", "Path on which to server filtered metrics")
	filter     = flag.String("filter", "filter.txt", "File containing a whitelist of metrics to filter on, one per line")
)

func newLogger() *logrus.Logger {
	var logger = logrus.New()
	logger.Out = os.Stderr
	jsonFormatter := new(logrus.JSONFormatter)
	jsonFormatter.TimestampFormat = time.RFC3339Nano
	logger.Formatter = jsonFormatter
	logger.Level = logrus.InfoLevel
	return logger
}

func main() {
	flag.Parse()
	logger := newLogger()

	server.RunHTTPListener(
		*endpoint, *listenPort, *path, *filter, logger,
	)
}
