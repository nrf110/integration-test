package integrationtest_test

import (
	"cloud.google.com/go/pubsub"
	"context"
	integrationtest "github.com/nrf110/integration-test/pkg/pubsub"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go/modules/gcloud"
)

var _ = Describe("pubsub.Dependency", func() {
	It("should publish and consume", func(ctx SpecContext) {
		topicName := "testtopic"
		subscriptionID := "testsubscription"

		pub := integrationtest.NewDependency(
			integrationtest.WithContainerOpts(
				gcloud.WithProjectID("test")))

		err := pub.Start(ctx)
		Expect(err).NotTo(HaveOccurred())

		client := pub.Client().(*pubsub.Client)
		topic, err := client.CreateTopic(ctx, topicName)
		Expect(err).NotTo(HaveOccurred())

		subscription, err := client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
			Topic: topic,
		})
		Expect(err).NotTo(HaveOccurred())

		topic.Publish(ctx, &pubsub.Message{Data: []byte("test")})

		subscription.Receive(ctx, func(ctx2 context.Context, m *pubsub.Message) {
			text := string(m.Data)
			Expect(text).To(Equal("test"))
			pub.Client().(*pubsub.Client).Close()
		})

		pub.Stop(ctx)
	})
})
