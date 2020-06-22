package main

import "github.com/prometheus/client_golang/prometheus"

type vhost struct {
	Messages    int    `json:"messages"`
	Name        string `json:"name"`
	ClusterName string
}

func (v *vhost) update(localVhost *vhost) {
	v.Messages = localVhost.Messages
	v.Name = localVhost.Name
}

func (v *vhost) updateMetrics() {
	vhostGauges["messages"].With(prometheus.Labels{"cluster_name": v.ClusterName, "environment": environment, "vhost": v.Name}).Set(float64(v.Messages))
}

var vhostGauges = map[string]*prometheus.GaugeVec{
	"messages": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_vhost_messages",
			Help: "Current number of messages in the vhost",
		},
		[]string{"cluster_name", "vhost", "environment"}),
}

func registerVhostMetrics() {
	for _, p := range vhostGauges {
		prometheus.MustRegister(p)
	}
}
