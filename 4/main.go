package main

import (
	"context"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	ctx = context.Background()
)

func publish(cl *kgo.Client, topic string, value []byte) {
	if err := cl.ProduceSync(ctx, &kgo.Record{Topic: topic, Value: value}).FirstErr(); err != nil {
		fmt.Printf("record had a produce error while synchronously producing: %v\n", err)
	}
}

func read(cl *kgo.Client, topic string) ([]*kgo.Record, error) {
	fetches := cl.PollFetches(ctx)
	if errs := fetches.Errors(); len(errs) > 0 {
		return nil, fmt.Errorf("%v", fmt.Sprint(errs))
	}

	iter := fetches.RecordIter()
	records := []*kgo.Record{}
	for !iter.Done() {
		record := iter.Next()
		records = append(records, record)
	}
	if err := cl.CommitUncommittedOffsets(ctx); err != nil {
		return nil, err
	}
	return records, nil
}

func main() {

	seeds := []string{"localhost:9092"}
	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup("my-group-identifier"),
		kgo.ConsumeTopics("foo"),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		log.Fatal(err)
	}
	client2, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup("my-group-identifier"),
		kgo.ConsumeTopics("bar"),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	defer client2.Close()

	for i := 0; i < 5; i++ {
		publish(client, "foo", []byte(fmt.Sprint("bar ", i)))
	}
	for {
		records, err := read(client, "foo")
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range records {
			log.Println(string(r.Value))
			publish(client, "bar", r.Value)
		}
		records, err = read(client2, "bar")
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range records {
			log.Println(string(r.Value))
		}
	}
}
