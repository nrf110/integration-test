package integrationtest

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"strings"
)

const defaultKafkaImage = "confluentinc/confluent-local:7.5.0"

type KafkaDependency struct {
	Dependency

	image         string
	containerOpts []testcontainers.ContainerCustomizer
	container     *kafka.KafkaContainer
	clientOpts    []kgo.Opt
	client        *kgo.Client
	adminClient   *kadm.Client
	env           map[string]string
	topics        []*topicConfig
}

func NewKafkaDependency(opts ...KafkaDependencyOpt) (*KafkaDependency, error) {
	dep := &KafkaDependency{image: defaultKafkaImage}
	for _, opt := range opts {
		if err := opt(dep); err != nil {
			return nil, err
		}
	}
	return dep, nil
}

type KafkaDependencyOpt func(*KafkaDependency) error

func WithKafkaImage(image string) KafkaDependencyOpt {
	return func(kd *KafkaDependency) error {
		kd.image = image
		return nil
	}
}

func WithKafkaClientOpts(opts ...kgo.Opt) KafkaDependencyOpt {
	return func(d *KafkaDependency) error {
		d.clientOpts = append(d.clientOpts, opts...)
		return nil
	}
}

func WithKafkaContainerOpts(opts ...testcontainers.ContainerCustomizer) KafkaDependencyOpt {
	return func(d *KafkaDependency) error {
		d.containerOpts = append([]testcontainers.ContainerCustomizer{
			kafka.WithClusterID("test-cluster"),
		}, opts...)

		return nil
	}
}

func WithKafkaTopic(name string, opts ...KafkaTopicOption) KafkaDependencyOpt {
	return func(d *KafkaDependency) error {
		topic := &topicConfig{
			name:       name,
			partitions: 1,
		}

		for _, opt := range opts {
			if err := opt(topic); err != nil {
				return err
			}
		}

		return nil
	}
}

func (k *KafkaDependency) Start(ctx context.Context) error {
	container, err := kafka.Run(ctx, k.image, k.containerOpts...)
	if err != nil {
		return err
	}
	k.container = container

	seeds, err := container.Brokers(ctx)
	if err != nil {
		return err
	}

	clientOpts := append([]kgo.Opt{
		kgo.SeedBrokers(seeds...),
	}, k.clientOpts...)

	client, err := kgo.NewClient(clientOpts...)

	k.adminClient = kadm.NewClient(client)

	k.client = client
	k.env = map[string]string{
		"KAFKA_BROKERS": strings.Join(seeds, ","),
	}

	for _, topic := range k.topics {
		if _, err = k.adminClient.CreateTopics(
			ctx,
			topic.partitions,
			1,
			topic.properties,
			topic.name,
		); err != nil {
			return err
		}
	}

	return nil
}

func (k *KafkaDependency) Env() map[string]string {
	return k.env
}

func (k *KafkaDependency) Client() any {
	return k.client
}

func (k *KafkaDependency) Stop(ctx context.Context) error {
	if k.container == nil {
		return k.container.Terminate(ctx)
	}
	return nil
}

func (k *KafkaDependency) Produce(ctx context.Context, records ...*kgo.Record) error {
	results := k.client.ProduceSync(ctx, records...)
	return results.FirstErr()
}

func (k *KafkaDependency) ConsumeWhile(ctx context.Context, predicate func(int, *kgo.Record) (bool, error)) error {
	var (
		more bool
		err  error
	)

	for more {
		fetches := k.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			return errs[0].Err
		}

		for idx, record := range fetches.Records() {
			if more, err = predicate(idx, record); err != nil {
				return err
			}
		}
	}
	return nil
}
