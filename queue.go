package main

import "github.com/prometheus/client_golang/prometheus"

type queue struct {
	Consumers       int    `json:"consumers"`
	Memory          int    `json:"memory"`
	MessageBytes    int    `json:"message_bytes"`
	MessageBytesRAM int    `json:"message_bytes_ram"`
	Messages        int    `json:"messages"`
	MessagesRAM     int    `json:"messages_ram"`
	Name            string `json:"name"`
	Node            string `json:"node"`
	State           string `json:"state"`
	Vhost           string `json:"vhost"`
	ClusterName     string
}

func (q *queue) update(localQueue *queue) {
	q.Consumers = localQueue.Consumers
	q.Memory = localQueue.Memory
	q.MessageBytes = localQueue.MessageBytes
	q.MessageBytesRAM = localQueue.MessageBytesRAM
	q.Messages = localQueue.Messages
	q.MessagesRAM = localQueue.MessagesRAM
	q.Name = localQueue.Name
	q.Node = localQueue.Node
	q.State = localQueue.State
	q.Vhost = localQueue.Vhost
}

func (q *queue) updateMetrics() {
	queueGauges["consumers"].With(prometheus.Labels{
		"cluster_name": q.ClusterName,
		"environment":  environment,
		"vhost":        q.Name,
		"node":         q.Node,
		"queue":        q.Name,
	}).Set(float64(q.Consumers))
	queueGauges["memory"].With(prometheus.Labels{
		"cluster_name": q.ClusterName,
		"environment":  environment,
		"vhost":        q.Name,
		"node":         q.Node,
		"queue":        q.Name,
	}).Set(float64(q.Memory))
	queueGauges["message_bytes"].With(prometheus.Labels{
		"cluster_name": q.ClusterName,
		"environment":  environment,
		"vhost":        q.Name,
		"node":         q.Node,
		"queue":        q.Name,
	}).Set(float64(q.MessageBytes))
	queueGauges["message_bytes_ram"].With(prometheus.Labels{
		"cluster_name": q.ClusterName,
		"environment":  environment,
		"vhost":        q.Name,
		"node":         q.Node,
		"queue":        q.Name,
	}).Set(float64(q.MessageBytesRAM))
	queueGauges["messages"].With(prometheus.Labels{
		"cluster_name": q.ClusterName,
		"environment":  environment,
		"vhost":        q.Name,
		"node":         q.Node,
		"queue":        q.Name,
	}).Set(float64(q.Messages))
	queueGauges["messages_ram"].With(prometheus.Labels{
		"cluster_name": q.ClusterName,
		"environment":  environment,
		"vhost":        q.Name,
		"node":         q.Node,
		"queue":        q.Name,
	}).Set(float64(q.MessagesRAM))
	if q.State == "running" {
		queueGauges["running"].With(prometheus.Labels{
			"cluster_name": q.ClusterName,
			"environment":  environment,
			"vhost":        q.Name,
			"node":         q.Node,
			"queue":        q.Name,
		}).Set(1)
	} else {
		queueGauges["running"].With(prometheus.Labels{
			"cluster_name": q.ClusterName,
			"environment":  environment,
			"vhost":        q.Name,
			"node":         q.Node,
			"queue":        q.Name,
		}).Set(0)
	}
}

var queueGauges = map[string]*prometheus.GaugeVec{
	"consumers": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_queue_consumers",
			Help: "Current number of consumers for the queue",
		},
		[]string{"cluster_name", "vhost", "node", "queue", "environment"}),
	"memory": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_queue_memory",
			Help: "Current memory consumed by the queue in bytes",
		},
		[]string{"cluster_name", "vhost", "node", "queue", "environment"}),
	"message_bytes": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_queue_message_bytes",
			Help: "Current size of messages in the queue in bytes",
		},
		[]string{"cluster_name", "vhost", "node", "queue", "environment"}),
	"message_bytes_ram": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_queue_message_bytes_ram",
			Help: "Current size of messages in the queue in RAM in bytes",
		},
		[]string{"cluster_name", "vhost", "node", "queue", "environment"}),
	"messages": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_queue_messages",
			Help: "Total number of messages in the queue",
		},
		[]string{"cluster_name", "vhost", "node", "queue", "environment"}),
	"messages_ram": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_queue_messages_ram",
			Help: "Total number of messages in RAM in the queue",
		},
		[]string{"cluster_name", "vhost", "node", "queue", "environment"}),
	"running": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_queue_running",
			Help: "Indicates if the current queue is running",
		},
		[]string{"cluster_name", "vhost", "node", "queue", "environment"}),
}

func registerQueueMetrics() {
	for _, p := range queueGauges {
		prometheus.MustRegister(p)
	}
}
