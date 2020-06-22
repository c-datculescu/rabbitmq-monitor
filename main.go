package main

import (
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	environment   = os.Getenv("ENV")
	clusterName   = os.Getenv("CLUSTER_NAME")
	username      = os.Getenv("USERNAME")
	password      = os.Getenv("PASSWORD")
	sleepInterval = os.Getenv("SLEEP_INTERVAL")
	host          = os.Getenv("HOST")
	sleepDuration time.Duration
)

func init() {
	registerClusterMetrics()
	registerNodeMetrics()
	registerVhostMetrics()
	registerQueueMetrics()
	registerShovelMetrics()
	registerFederationLinksMetrics()

	// process duration
	sleep, err := time.ParseDuration(sleepInterval)
	if err != nil {
		os.Exit(1)
	}
	sleepDuration = sleep
}

func main() {
	localCluster := &cluster{
		Address:     host,
		Username:    username,
		Password:    password,
		ClusterName: clusterName,
	}

	go func() {
		for {
			localCluster.scan()
			time.Sleep(sleepDuration)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":17762", nil)
}
