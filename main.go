package main

import (
	"context"
	"flag"
	"log"
	"strconv"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

// the topic and broker address are initialized as constants
// const (
//     topic          = "test"
//     broker0Address = "localhost:9094"
//     broker1Address = "localhost:9093"
//     broker2Address = "localhost:9094"
//     broker3Address = "localhost:9095"
// )

func produce(ctx context.Context) {
	// initialize a counter
	i := 0

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		//         Brokers: []string{broker1Address, broker2Address, broker3Address},
		Brokers: Config.Brokers,
		Topic:   Config.Topic,
		// wait until we get 10 messages before writing
		// if we want to send messages immediately set it to 1
		BatchSize: 10,
		// no matter what happens, write all pending messages
		// every 2 seconds
		BatchTimeout: 2 * time.Second,
		// can be set to -1, 0, or 1
		// 1 is a good default for most non-transactional data
		RequiredAcks: 1,
	})

	for {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			// create an arbitrary message payload for the value
			Value: []byte("this is message" + strconv.Itoa(i)),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		log.Println("writes:", i)
		i++
		// sleep for a second
		time.Sleep(time.Second)
	}
}

func consume(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		//         Brokers: []string{broker1Address, broker2Address, broker3Address},
		Brokers:  Config.Brokers,
		Topic:    Config.Topic,
		GroupID:  Config.Group,
		MinBytes: 5,
		MaxBytes: 1e6,
		// wait for at most 3 seconds before receiving new data
		MaxWait: 3 * time.Second,
		// this will start consuming messages from the earliest available
		StartOffset: kafka.FirstOffset,
		// if you set it to `kafka.LastOffset` it will only consume new messages
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		log.Println("received: ", string(msg.Value))
	}
}

func main() {
	var config string
	flag.StringVar(&config, "config", "config.json", "config file")
	flag.Parse()

	log.SetFlags(0)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// create a new context
	ctx := context.Background()
	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking
	if Config.Produce {
		go produce(ctx)
	}
	consume(ctx)
}
