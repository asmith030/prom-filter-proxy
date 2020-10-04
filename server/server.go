package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/asmith030/prom-filter-proxy/server/handlers"
)

func RunHTTPListener(endpoint string, port int, path string, filter string, logger *logrus.Logger) error {
	m := http.NewServeMux()

	hh := &handlers.HealthHandler{
		Healthy: false,
		Logger:  logger,
	}
	m.Handle("/healthz", hh)
	mh := &handlers.MetricHandler{
		Endpoint: endpoint,
		Logger:   logger,
	}
	mh.SetWhitelistFile(filter)

	m.Handle(path, mh)

	logger.Infof("Starting the webserver on port %v", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	hh.Healthy = true
	return http.Serve(lis, m)
}
