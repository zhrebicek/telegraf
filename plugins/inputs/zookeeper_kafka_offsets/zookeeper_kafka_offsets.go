package zookeeper_kafka_offsets

import (
	"strconv"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

// Zookeeper configuration
type ZookeeperKafkaOffsets struct {
	Servers []string
}

var sampleConfig = `
	## An array of zookeeper address to gather offsets from.
	## Specify an ip or hostname with port. ie localhost, 10.0.0.1, etc.
	## If no servers are specified, then localhost is used as the address.
	servers = ["127.0.0.1"]
`

var defaultTimeout = time.Second * time.Duration(5)

// SampleConfig returns sample configuration message
func (z *ZookeeperKafkaOffsets) SampleConfig() string {
	return sampleConfig
}

// Description returns description of Zookeeper plugin
func (z *ZookeeperKafkaOffsets) Description() string {
	return `Reads kafka offsets from zookeeper. The path for reading the offsets is:
# /consumers/<consumer>/offsets/<topic>/<partitions>`
}

// Gather reads stats from all configured servers accumulates stats
func (z *ZookeeperKafkaOffsets) Gather(acc telegraf.Accumulator) error {
	if len(z.Servers) == 0 {
		z.Servers = []string{"127.0.0.1"}
	}

	for _, serverAddress := range z.Servers {
		acc.AddError(z.gatherServer(serverAddress, acc))
	}
	return nil
}

func (z *ZookeeperKafkaOffsets) gatherServer(address string, acc telegraf.Accumulator) error {

	c, _, err := zk.Connect([]string{address}, defaultTimeout)
	defer c.Close()

	// We are getting part of the path that is unknown to us and iterate over it
	services, _, _, err := c.ChildrenW("/consumers")
	if err != nil {
		acc.AddError(err)
	}

	for _, service := range services {
		path := strings.Join([]string{"/consumers", service, "offsets"},"/")
		exists, _, _ := c.Exists(path)
		serviceTopics := []string{}
		if exists {
			var err error
			serviceTopics, _, err = c.Children(path)
			if err != nil {
				acc.AddError(err)
			}
		} 
		for _, topic := range serviceTopics {
			pathToTopic := strings.Join([]string{path, topic},"/")
			partitions, _, err := c.Children(pathToTopic)
			if err != nil {
				acc.AddError(err)
			}
			for _, partition := range partitions {
				pathToPartition := strings.Join([]string{pathToTopic, partition},"/")
				value, _, err := c.Get(pathToPartition)

				if err != nil{
					acc.AddError(err)
				} else {
				result := make(map[string]interface{})
				offsets := "offsets"
				result[offsets], _ = strconv.Atoi(string(value))
				tags := map[string]string{
					"service": service,
					"topic": topic,
					"partition": partition,
				}
				acc.AddFields("zookeeper_kafka", result, tags)
				}
			}

		}
	}

	return nil
}

func init() {
	inputs.Add("zookeeper_kafka_offsets", func() telegraf.Input {
		return &ZookeeperKafkaOffsets{}
	})
}
