package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type federationLink struct {
	Status       string                  `json:"status"`
	Vhost        string                  `json:"vhost"`
	Name         string                  `json:"upstream"`
	Node         string                  `json:"node"`
	LocalChannel *federationLocalChannel `json:"local_channel"`
	ClusterName  string
}

func (fl *federationLink) update(localFl *federationLink) {
	fl.Status = localFl.Status
	fl.Vhost = localFl.Vhost
	fl.Name = localFl.Name
	fl.Node = localFl.Name
	fl.LocalChannel = localFl.LocalChannel
}

func (fl *federationLink) updateMetrics() {
	if fl.Status == "running" {
		federationLinksGauges["running"].With(prometheus.Labels{
			"cluster_name": fl.ClusterName,
			"vhost":        fl.Vhost,
			"name":         fl.Name,
			"node":         fl.Node,
			"environment":  environment,
		}).Set(1)
	} else {
		federationLinksGauges["running"].With(prometheus.Labels{
			"cluster_name": fl.ClusterName,
			"vhost":        fl.Vhost,
			"name":         fl.Name,
			"node":         fl.Node,
			"environment":  environment,
		}).Set(0)
	}

	if fl.LocalChannel.State == "running" {
		federationLinksGauges["channel_running"].With(prometheus.Labels{
			"cluster_name": fl.ClusterName,
			"vhost":        fl.Vhost,
			"name":         fl.Name,
			"node":         fl.Node,
			"environment":  environment,
		}).Set(1)
	} else {
		federationLinksGauges["channel_running"].With(prometheus.Labels{
			"cluster_name": fl.ClusterName,
			"vhost":        fl.Vhost,
			"name":         fl.Name,
			"node":         fl.Node,
			"environment":  environment,
		}).Set(0)
	}
	federationLinksGauges["messages_unacknowledged"].With(prometheus.Labels{
		"cluster_name": fl.ClusterName,
		"vhost":        fl.Vhost,
		"name":         fl.Name,
		"node":         fl.Node,
		"environment":  environment,
	}).Set(float64(fl.LocalChannel.MessagesUnacknowledged))

	federationLinksGauges["messages_uncommited"].With(prometheus.Labels{
		"cluster_name": fl.ClusterName,
		"vhost":        fl.Vhost,
		"name":         fl.Name,
		"node":         fl.Node,
		"environment":  environment,
	}).Set(float64(fl.LocalChannel.MessagesUncommited))

	federationLinksGauges["messages_unconfirmed"].With(prometheus.Labels{
		"cluster_name": fl.ClusterName,
		"vhost":        fl.Vhost,
		"name":         fl.Name,
		"node":         fl.Node,
		"environment":  environment,
	}).Set(float64(fl.LocalChannel.MessagesUnconfirmed))
}

type federationLocalChannel struct {
	MessagesUnacknowledged int    `json:"messages_unacknowledged"`
	MessagesUncommited     int    `json:"messages_uncommited"`
	MessagesUnconfirmed    int    `json:"messages_unconfirmed"`
	State                  string `json:"state"`
}

var federationLinksGauges = map[string]*prometheus.GaugeVec{
	"running": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_federation_link_running",
			Help: "Indicates if the current federation link is in a running state",
		},
		[]string{"cluster_name", "vhost", "name", "node", "environment"}),
	"messages_unacknowledged": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_federation_link_messages_unacknowledged",
			Help: "The current amount of unacknowledged messages",
		},
		[]string{"cluster_name", "vhost", "name", "node", "environment"}),
	"messages_uncommited": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_federation_link_messages_uncommited",
			Help: "The current amount of uncommited messages",
		},
		[]string{"cluster_name", "vhost", "name", "node", "environment"}),
	"messages_unconfirmed": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_federation_link_messages_unconfirmed",
			Help: "The current amount of unconfirmed messages",
		},
		[]string{"cluster_name", "vhost", "name", "node", "environment"}),
	"channel_running": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_federation_link_channel_running",
			Help: "Indicates whether the underlying channel is running",
		},
		[]string{"cluster_name", "vhost", "name", "node", "environment"}),
}

func registerFederationLinksMetrics() {
	for _, p := range federationLinksGauges {
		prometheus.MustRegister(p)
	}
}
