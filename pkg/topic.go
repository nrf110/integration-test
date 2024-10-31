package integrationtest

import (
	"fmt"
	"slices"
)

type topicConfig struct {
	name       string
	partitions int32
	properties map[string]*string
}

type KafkaTopicOption func(config *topicConfig) error

func WithPartitions(value int32) KafkaTopicOption {
	return func(config *topicConfig) error {
		config.partitions = value
		return nil
	}
}

func WithProperties(pairs ...string) KafkaTopicOption {
	return func(config *topicConfig) error {
		if length := len(pairs); length < 2 || length%2 != 0 {
			return fmt.Errorf("properties must be pairs of keys and values")
		}

		for pair := range slices.Chunk(pairs, 2) {
			config.properties[pair[0]] = &pair[1]
		}

		return nil
	}
}
