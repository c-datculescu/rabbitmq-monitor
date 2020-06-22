package main

import "github.com/prometheus/client_golang/prometheus"

type shovel struct {
	State       string `json:"state"`
	Name        string `json:"name"`
	Node        string `json:"node"`
	ClusterName string
}

func (s *shovel) update(localShovel *shovel) {
	s.State = localShovel.Name
	s.Name = localShovel.Name
	s.Node = localShovel.Node
}

func (s *shovel) updateMetrics() {
	if s.State == "running" {
		shovelGauges["running"].With(prometheus.Labels{
			"cluster_name": s.ClusterName,
			"name":         s.Name,
			"node":         s.Node,
			"environment":  environment,
		}).Set(1)
	} else {
		shovelGauges["running"].With(prometheus.Labels{
			"cluster_name": s.ClusterName,
			"name":         s.Name,
			"node":         s.Node,
			"environment":  environment,
		}).Set(0)
	}
}

var shovelGauges = map[string]*prometheus.GaugeVec{
	"running": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_shovel_running",
			Help: "Indicates if the current shovel is running",
		},
		[]string{"cluster_name", "name", "node", "environment"}),
}

func registerShovelMetrics() {
	for _, p := range shovelGauges {
		prometheus.MustRegister(p)
	}
}
