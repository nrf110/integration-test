package integrationtest

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	integrationtest "github.com/nrf110/integration-test/pkg/pubsub"
	"github.com/stretchr/testify/assert"
	tcpubsub "github.com/testcontainers/testcontainers-go/modules/gcloud/pubsub"
)

func TestPubSubDependency(t *testing.T) {
	t.Run("can publish and consume", func(t *testing.T) {
		topicName := "testtopic"
		subscriptionID := "testsubscription"

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		t.Cleanup(cancel)

		pub := integrationtest.NewDependency(
			integrationtest.WithContainerOpts(
				tcpubsub.WithProjectID("test")))

		err := pub.Start(ctx)
		assert.NoError(t, err)

		t.Cleanup(func() {
			assert.NoError(t, pub.Stop(ctx))
		})

		client := pub.Client().(*pubsub.Client)
		topic, err := client.CreateTopic(ctx, topicName)
		assert.NoError(t, err)

		subscription, err := client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
			Topic: topic,
		})
		assert.NoError(t, err)

		topic.Publish(ctx, &pubsub.Message{Data: []byte("test")})

		subscription.Receive(ctx, func(ctx2 context.Context, m *pubsub.Message) {
			text := string(m.Data)
			assert.Equal(t, "test", text)
			assert.NoError(t, pub.Client().(*pubsub.Client).Close())
		})
	})
}
