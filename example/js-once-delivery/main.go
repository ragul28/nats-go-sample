package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Connect and get the JetStream context.
	nc, _ := nats.Connect(nats.DefaultURL)
	js, _ := nc.JetStream()

	// Create a test stream.
	_, err := js.AddStream(&nats.StreamConfig{
		Name:       "test",
		Storage:    nats.MemoryStorage,
		Subjects:   []string{"test.>"},
		Duplicates: time.Minute,
	})
	failOnErr(err)

	defer js.DeleteStream("test")

	// Publish some messages with duplicates.
	js.Publish("test.1", []byte("hello"), nats.MsgId("1"))
	js.Publish("test.2", []byte("world"), nats.MsgId("2"))
	js.Publish("test.1", []byte("hello"), nats.MsgId("1"))
	js.Publish("test.1", []byte("hello"), nats.MsgId("1"))
	js.Publish("test.2", []byte("world"), nats.MsgId("2"))
	js.Publish("test.2", []byte("world"), nats.MsgId("2"))

	// Create an explicit pull consumer on the stream.
	_, err = js.AddConsumer("test", &nats.ConsumerConfig{
		Durable:       "test",
		AckPolicy:     nats.AckExplicitPolicy,
		DeliverPolicy: nats.DeliverAllPolicy,
	})
	failOnErr(err)
	defer js.DeleteConsumer("test", "test")

	// Create a subscription on the pull consumer.
	// Subject can be empty since it defaults to all subjects bound to the stream.
	sub, err := js.PullSubscribe("", "test", nats.BindStream("test"))
	failOnErr(err)

	// Only two should be delivered.
	batch, _ := sub.Fetch(10)
	log.Printf("%d messages", len(batch))

	// AckSync both to ensure the server received the ack.
	batch[0].AckSync()
	batch[1].AckSync()

	// Should be zero.
	batch, _ = sub.Fetch(10, nats.MaxWait(time.Second))
	log.Printf("%d messages", len(batch))
}
