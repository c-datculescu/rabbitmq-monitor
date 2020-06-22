# Monitoring solution for RabbitMQ clusters

## Healthchecking metrics

- the healthcheck of the cluster should be checked more frequently than other things
- it should check the following metrics
  - api is reachable - this usually indicates an issue with the api
    - prometheus metrics
      - ww_rmq_api_reachable{cluster}
      - ww_rmq_api_latency{cluster}
  - attempt to connect to the node - this usually indicates a issue with the rmq itself
    - prometheus metrics
      - ww_rmq_core_reachable{cluster}
      - ww_rmq_core_latency{cluster}
  - node monitoring - this can indicate multiple issues with the nodes
    - prometheus metrics:
      - ww_rmq_node_fd_max{cluster, node}
      - ww_rmq_node_fd_current{cluster, node}
      - ww_rmq_node_sock_max{cluster, node}
      - ww_rmq_node_sock_current{cluster, node}
      - ww_rmq_node_proc_max{cluster, node}
      - ww_rmq_node_proc_current{cluster, node}
      - ww_rmq_mem_max{cluster, node}
      - ww_rmq_mem_current{cluster, node}
      - ww_rmq_disk_min{cluster, node}
      - ww_rmq_disk_current{cluster, node}
      - ww_rmq_mem_alarm{cluster, node}
      - ww_rmq_disk_alarm{cluster, node}
      - ww_rmq_context_switch{cluster, node}

## Additional metrics

The additional metrics expose things like vhosts, and are less granular than the healthchecking metrics

- vhost metrics
  - ww_rmq_vhost_messages_ready{vhost, cluster}
  - ww_rmq_vhost_messages_unacknowledged{vhost, cluster}
  - ww_rmq_vhost_messages_ready{vhost, cluster}
- queue metrics
  - ww_rmq_queue_autodelete{vhost, queue, cluster}
  - ww_rmq_queue_consumer_utilization{vhost, queue, cluster}
  - ww_rmq_queue_consumer_count{vhost, queue, cluster}
  - ww_rmq_queue_durable{vhost, queue, cluster}
  - ww_rmq_queue_exclusive{vhost, queue, cluster}
  - ww_rmq_queue_memory{vhost, queue, cluster}
  - ww_rmq_queue_messages_bytes{vhost, queue, cluster}
  - ww_rmq_queue_messages_bytes_ram{vhost, queue, cluster}
  - ww_rmq_queue_messages{vhost, queue, cluster}
  - ww_rmq_queue_messages_ram{vhost, queue, cluster}

## Alerting

Alerting will be handled via alertmanager and falcon