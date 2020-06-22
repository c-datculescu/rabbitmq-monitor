package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/streadway/amqp"
)

type cluster struct {
	Address  string
	Username string
	Password string

	Nodes           map[string]*node
	Vhosts          map[string]*vhost
	Queues          map[string]*queue
	Shovels         map[string]*shovel
	FederationLinks map[string]*federationLink

	ClusterName string // the current name of the cluster

	apiReachable       int                  // indicates if the api is reachable
	apiReachableGauge  *prometheus.GaugeVec // the gauge for the api reachable
	apiLatency         int                  // indicates what is the observed api latency
	apiLatencyGauge    *prometheus.GaugeVec // the gauge for the latency for connecting to the api
	coreReachable      int                  // indicates if rmq core is reachable
	coreReachableGauge *prometheus.GaugeVec // the core reachability gauge
	coreLatency        int                  // indicates what is the observed latency for connecting to core
	coreLatencyGauge   *prometheus.GaugeVec // the core latency gauge
}

// runs all the scans for the api and the core
func (c *cluster) scan() {
	log.Println("Scanning...")
	// reset all fields when we enter the scan so we dont leave behind stuff that is not valid anymore
	c.apiReachable = 0
	c.apiLatency = 0
	c.coreReachable = 0
	c.coreLatency = 0

	if c.Nodes == nil {
		c.Nodes = map[string]*node{}
	}
	if c.Vhosts == nil {
		c.Vhosts = map[string]*vhost{}
	}
	if c.Queues == nil {
		c.Queues = map[string]*queue{}
	}
	if c.Shovels == nil {
		c.Shovels = map[string]*shovel{}
	}
	if c.FederationLinks == nil {
		c.FederationLinks = map[string]*federationLink{}
	}

	// populate the fields
	c.connect()
	c.apiConnect()
	c.nodes()
	c.vhosts()
	c.queues()
	c.shovels()
	c.federationLinks()

	// attach all the metrics to the prometheus instance
	c.updateMetrics()
}

func (c *cluster) updateMetrics() {
	clusterGauges["api_reachable"].With(prometheus.Labels{"cluster_name": c.ClusterName, "environment": environment}).Set(float64(c.apiReachable))
	clusterGauges["api_latency"].With(prometheus.Labels{"cluster_name": c.ClusterName, "environment": environment}).Set(float64(c.apiLatency))
	clusterGauges["core_reachable"].With(prometheus.Labels{"cluster_name": c.ClusterName, "environment": environment}).Set(float64(c.coreReachable))
	clusterGauges["core_latency"].With(prometheus.Labels{"cluster_name": c.ClusterName, "environment": environment}).Set(float64(c.coreLatency))

	for _, node := range c.Nodes {
		node.updateMetrics()
	}

	for _, vhost := range c.Vhosts {
		vhost.updateMetrics()
	}

	for _, queue := range c.Queues {
		queue.updateMetrics()
	}

	for _, shovel := range c.Shovels {
		shovel.updateMetrics()
	}

	for _, fl := range c.FederationLinks {
		fl.updateMetrics()
	}
}

// connects to the cluster using the rabbitmq connector
func (c *cluster) connect() {
	beforeConn := time.Now().UnixNano()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:5672/%%2F", c.Username, c.Password, c.Address))
	afterConn := time.Now().UnixNano()
	if err != nil {
		return
	}
	defer conn.Close()
	c.coreLatency = int(afterConn - beforeConn)
	c.coreReachable = 1
	return
}

// connects to the cluster api using web calls
func (c *cluster) apiConnect() {
	beforeConn := time.Now().UnixNano()
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:15672/api/overview", c.Address), nil)
	if err != nil {
		return
	}
	request.SetBasicAuth(c.Username, c.Password)
	response, err := client.Do(request)
	afterConn := time.Now().UnixNano()
	if err != nil {
		return
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()

	c.apiLatency = int(afterConn - beforeConn)
	c.apiReachable = 1
}

// retrieves node information using the node api
func (c *cluster) nodes() {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:15672/api/nodes", c.Address), nil)
	request.SetBasicAuth(c.Username, c.Password)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	nodes := []*node{}
	err = json.Unmarshal(contents, &nodes)
	if err != nil {
		return
	}

	for _, node := range nodes {
		node.ClusterName = c.ClusterName
		if _, exists := c.Nodes[node.Name]; exists == true {
			c.Nodes[node.Name].update(node)
		} else {
			c.Nodes[node.Name] = node
		}
	}
}

// retrieves all the current vhosts using the vhosts api
func (c *cluster) vhosts() {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:15672/api/vhosts", c.Address), nil)
	request.SetBasicAuth(c.Username, c.Password)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	vhosts := []*vhost{}
	err = json.Unmarshal(contents, &vhosts)
	if err != nil {
		return
	}
	for _, vhost := range vhosts {
		vhost.ClusterName = c.ClusterName
		if _, exists := c.Vhosts[vhost.Name]; exists == true {
			c.Vhosts[vhost.Name].update(vhost)
		} else {
			c.Vhosts[vhost.Name] = vhost
		}
	}
}

func (c *cluster) queues() {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:15672/api/queues", c.Address), nil)
	request.SetBasicAuth(c.Username, c.Password)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	queues := []*queue{}
	err = json.Unmarshal(contents, &queues)
	if err != nil {
		return
	}

	for _, queue := range queues {
		queue.ClusterName = c.ClusterName
		if _, exists := c.Queues[queue.Name]; exists == true {
			c.Queues[queue.Name].update(queue)
		} else {
			c.Queues[queue.Name] = queue
		}
	}
}

func (c *cluster) shovels() {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:15672/api/shovels", c.Address), nil)
	request.SetBasicAuth(c.Username, c.Password)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	shovels := []*shovel{}
	err = json.Unmarshal(contents, &shovels)
	if err != nil {
		return
	}

	for _, shovel := range shovels {
		shovel.ClusterName = c.ClusterName
		if _, exists := c.Shovels[shovel.Name]; exists == true {
			c.Shovels[shovel.Name].update(shovel)
		} else {
			c.Shovels[shovel.Name] = shovel
		}
	}
}

func (c *cluster) federationLinks() {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%s:15672/api/federation-links", c.Address), nil)
	request.SetBasicAuth(c.Username, c.Password)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	defer client.CloseIdleConnections()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	allLinks := []*federationLink{}
	err = json.Unmarshal(contents, &allLinks)
	if err != nil {
		return
	}

	for _, link := range allLinks {
		link.ClusterName = c.ClusterName
		if _, exists := c.FederationLinks[link.Name]; exists == true {
			c.FederationLinks[link.Name].update(link)
		} else {
			c.FederationLinks[link.Name] = link
		}
	}
}

var clusterGauges = map[string]*prometheus.GaugeVec{
	"api_reachable": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_api_reachable",
			Help: "The api is reachable over local interface",
		},
		[]string{"cluster_name", "environment"}),
	"api_latency": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_api_latency",
			Help: "The api latency over the local interface",
		},
		[]string{"cluster_name", "environment"}),
	"core_reachable": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_core_reachable",
			Help: "The connectivity over the amqp protocol to the local instance",
		},
		[]string{"cluster_name", "environment"}),
	"core_latency": prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rmq_core_latency",
			Help: "The latency of connecting to the local server using the amqp protocol",
		},
		[]string{"cluster_name", "environment"}),
}

func registerClusterMetrics() {
	for _, p := range clusterGauges {
		prometheus.MustRegister(p)
	}
}
