## Telegraf Plugin: Zookeeper Kafka Offsets

#### Description

The zookeeper kafka offsets plugin collects kafka offsets with help of
Native Go Zookeeper Client Library https://github.com/samuel/go-zookeeper.

This plugin is supposed to collect metrics from same node on which the telegraf
and your zookeeper is running, because of the collection speed. But it is still
left configurable to collect from more nodes but note that is was not tested.

## Configuration

```
# Reads kafka offsets from zookeeper. The path for reading the offsets is:
# /consumers/<consumer>/offsets/<topic>/<partitions>
[[inputs.zookeeper_kafka_offsets]]
	## An array of zookeeper address to gather offsets from.
	## Specify an ip. ie localhost, 10.0.0.1, etc.
	## If no servers are specified, then localhost is used as the address.
	servers = ["127.0.0.1"]
```

### Metrics:

- zookeeper_kafka_offsets
  - tags:
    - partition (number of topic partition)
    - service (name of the consumer aka service)
    - topic (name of the kafka topic)
  - fields:
    - numeric value of the measurement


### Example Output:

```
zookeeper_kafka_offsets{host="hostname",partition="0",service="servicename",topic="topicName"} 46465465
zookeeper_kafka_offsets{host="hostname",partition="1",service="servicename",topic="topicName"} 46466500
zookeeper_kafka_offsets{host="hostname",partition="2",service="servicename",topic="topicName"} 462365465
zookeeper_kafka_offsets{host="hostname",partition="3",service="servicename",topic="topicName"} 41265465
```
