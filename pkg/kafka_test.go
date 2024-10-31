package integrationtest_test

import (
	integrationtest "github.com/nrf110/integration-test/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/twmb/franz-go/pkg/kgo"
)

var _ = Describe("KafkaDependency", func() {
	It("should enable producing and consuming", func(ctx SpecContext) {
		kafka, err := integrationtest.NewKafkaDependency(
			integrationtest.WithKafkaClientOpts(
				kgo.ConsumerGroup("test-group"),
				kgo.ConsumeTopics("test-topic"),
			),
		)
		Expect(err).To(BeNil())
		Expect(kafka.Start(ctx)).To(BeNil())
		defer func() {
			Expect(kafka.Stop(ctx)).To(BeNil())
		}()

		Expect(kafka.Produce(ctx, &kgo.Record{
			Key:   []byte("Hello"),
			Value: []byte("Goodbye"),
		})).To(BeNil())

		kafka.ConsumeWhile(ctx, func(_ int, record *kgo.Record) (bool, error) {
			return string(record.Key) == "Hello", nil
		})
	})
})
