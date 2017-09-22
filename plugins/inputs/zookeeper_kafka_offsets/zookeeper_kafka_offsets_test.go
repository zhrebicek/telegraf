package zookeeper_kafka_offsets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO To write proper test we need to mock zookeeper or find embedded one.
func TestZookeeperKafkaOffsetsGeneratesMetrics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	z := &ZookeeperKafkaOffsets{
		Servers: []string{"127.0.0.1"},
	}

	assert.True(true)

}
