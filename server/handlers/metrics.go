package handlers

import (
	"bufio"
	"net/http"
	"os"
	"regexp"

	"github.com/sirupsen/logrus"
)

var metrics_name_re = regexp.MustCompile(`^(# (HELP|TYPE) )?(?P<metric>[a-zA-Z_:][a-zA-Z0-9_:]*)`)

type MetricHandler struct {
	Endpoint string
	Logger   *logrus.Logger
	filter   func(string) bool
}

func (mh *MetricHandler) SetWhitelistFile(whitelist_file string) {
	whitelist := make(map[string]bool)

	file, err := os.Open(whitelist_file)
	if err != nil {
		mh.Logger.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) > 0 {
			whitelist[t] = true
		}
	}

	if err := scanner.Err(); err != nil {
		mh.Logger.Fatal(err)
	}

	mh.filter = func(metric string) bool {
		return whitelist[metric]
	}
}

func (mh *MetricHandler) filterLine(line string) bool {
	match := metrics_name_re.FindStringSubmatch(line)
	if len(match) < 1 {
		return false
	}
	for i, name := range metrics_name_re.SubexpNames() {
		if name == "metric" {
			return mh.filter(match[i])
		}
	}
	return false
}

func (mh *MetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(mh.Endpoint)
	if err != nil {
		mh.Logger.Infof("Error fetching metrics from upstream endpoint")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		t := scanner.Text()
		if mh.filterLine(t) {
			w.Write([]byte(t))
			w.Write([]byte("\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		mh.Logger.Infof("Error reading metrics from upstream endpoint")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
