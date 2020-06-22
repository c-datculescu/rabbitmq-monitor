package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type node struct {
	FdMax         int    `json:"fd_total"`
	FdCurrent     int    `json:"fd_used"`
	SockMax       int    `json:"sockets_total"`
	SockCurrent   int    `json:"sockets_used"`
	ProcMax       int    `json:"proc_total"`
	ProcCurrent   int    `json:"proc_used"`
	MemMax        int    `json:"mem_limit"`
	MemCurrent    int    `json:"mem_used"`
	DiskMin       int    `json:"disk_free_limit"`
	DiskCurrent   int    `json:"disk_free"`
	MemAlarm      bool   `json:"mem_alarm"`
	DiskAlarm     bool   `json:"disk_free_alarm"`
	ContextSwitch int    `json:"context_switches"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	ClusterName   string
}

func (n *node) update(localNode *node) {
	n.FdMax = localNode.FdMax
	n.FdCurrent = localNode.FdCurrent
	n.SockMax = localNode.SockMax
	n.SockCurrent = localNode.SockCurrent
	n.ProcMax = localNode.ProcMax
	n.ProcCurrent = localNode.ProcCurrent
	n.MemCurrent = localNode.MemCurrent
	n.DiskMin = localNode.DiskMin
	n.DiskCurrent = localNode.DiskCurrent
	n.MemAlarm = localNode.MemAlarm
	n.DiskAlarm = localNode.DiskAlarm
	n.ContextSwitch = localNode.ContextSwitch
	n.Name = localNode.Name
	n.Type = localNode.Type
}

func (n *node) updateMetrics() {
	nodeGauges["fd_max"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.FdMax))
	nodeGauges["fd_current"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.FdCurrent))
	nodeGauges["sock_max"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.SockMax))
	nodeGauges["sock_current"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.SockCurrent))
	nodeGauges["proc_max"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.ProcMax))
	nodeGauges["proc_current"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.ProcCurrent))
	nodeGauges["mem_max"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.MemMax))
	nodeGauges["mem_current"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.MemCurrent))
	nodeGauges["disk_min"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.DiskMin))
	nodeGauges["disk_current"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.DiskCurrent))
	if n.MemAlarm {
		nodeGauges["mem_alarm"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(1)
	} else {
		nodeGauges["mem_alarm"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(0)
	}
	if n.DiskAlarm {
		nodeGauges["disk_alarm"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(1)
	} else {
		nodeGauges["disk_alarm"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(0)
	}
	nodeGauges["context_switch"].With(prometheus.Labels{"cluster_name": n.ClusterName, "environment": environment, "node": n.Name}).Set(float64(n.ContextSwitch))
}

var nodeGauges = map[string]*prometheus.GaugeVec{
	"fd_max": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_fd_max",
			Help: "Maximum allowed number of file descriptors",
		},
		[]string{"cluster_name", "node", "environment"}),
	"fd_current": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_fd_current",
			Help: "Current amount of file descriptors",
		},
		[]string{"cluster_name", "node", "environment"}),
	"sock_max": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_sock_max",
			Help: "Maximum allowed number of sockets",
		},
		[]string{"cluster_name", "node", "environment"}),
	"sock_current": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_sock_current",
			Help: "Current amount of sockets",
		},
		[]string{"cluster_name", "node", "environment"}),
	"proc_max": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_proc_max",
			Help: "Maximum allowed number of erlang processes ",
		},
		[]string{"cluster_name", "node", "environment"}),
	"proc_current": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_proc_current",
			Help: "Current amount of erlang processes ",
		},
		[]string{"cluster_name", "node", "environment"}),
	"mem_max": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_mem_max",
			Help: "Maximum allowed memory consumption",
		},
		[]string{"cluster_name", "node", "environment"}),
	"mem_current": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_mem_current",
			Help: "Current amount of memory",
		},
		[]string{"cluster_name", "node", "environment"}),
	"disk_min": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_disk_min",
			Help: "Minimum amount of disk after which disk alarm triggers",
		},
		[]string{"cluster_name", "node", "environment"}),
	"disk_current": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_disk_current",
			Help: "Current amount of free disk",
		},
		[]string{"cluster_name", "node", "environment"}),
	"mem_alarm": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_mem_alarm",
			Help: "Indicates whether the memory alarm is active",
		},
		[]string{"cluster_name", "node", "environment"}),
	"disk_alarm": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_disk_alarm",
			Help: "Indicates whether the disk alarm is active",
		},
		[]string{"cluster_name", "node", "environment"}),
	"context_switch": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_node_context_switch",
			Help: "The current amount of context switches for the current node",
		},
		[]string{"cluster_name", "node", "environment"}),
}

func registerNodeMetrics() {
	for _, p := range nodeGauges {
		prometheus.MustRegister(p)
	}
}
