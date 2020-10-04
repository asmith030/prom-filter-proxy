# prom-filter-proxy
Simple proxy to filter the list of metrics served by a prometheus endpoint

## Usage
```
Usage of prom-filter-proxy:
  -endpoint string
    	Address of the metrics endpoint to proxy (default "http://localhost/metrics")
  -filter string
    	File containing a whitelist of metrics to filter on, one per line (default "filter.txt")
  -listenPort int
    	TCP port that we should listen on (default 8080)
  -path string
    	Path on which to server filtered metrics (default "/metrics")
```
