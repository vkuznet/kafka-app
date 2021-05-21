package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

// helper function to yield (produce) messages to kafka
func produce(ctx context.Context, brokers []string, topic, group string) {
	// initialize a counter
	i := 0

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
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

// helper function to consume messages from kafka
func consume(ctx context.Context, brokers []string, topic, group string) {
	// initialize a new reader with the brokers and receive topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  group,
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
			log.Println("could not read message " + err.Error())
			continue
		}
		log.Println("received: ", string(msg.Value))
	}
}

// helper function to process messages from kafka
func process(ctx context.Context, brokers []string, receiveTopic, sendTopic, group, alg string, verbose int) {
	// initialize a new reader with the brokers and receive topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    receiveTopic,
		GroupID:  group,
		MinBytes: 5,
		MaxBytes: 1e6,
		// wait for at most 3 seconds before receiving new data
		MaxWait: 3 * time.Second,
		// this will start consuming messages from the earliest available
		StartOffset: kafka.FirstOffset,
		// if you set it to `kafka.LastOffset` it will only consume new messages
	})
	log.Println("kafka reader", r)

	// intialize the writer with the broker addresses, and the send topic
	var w *kafka.Writer
	if sendTopic != "" {
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   sendTopic,
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
		log.Println("kafka writer", w)
	}

	// main loop which receive message from kafka stream,
	// then process is to anonimise the fields and send it back
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("could not read message " + err.Error())
			continue
		}
		if verbose > 0 {
			log.Println("received: ", string(msg.Value))
		}

		// anonimise our data
		data, err := anonimise(alg, sendTopic, msg.Value)
		if err != nil {
			log.Println("fail to anonimise the data", err)
			continue
		}
		// if we have a writer we'll send our data
		if w != nil {
			ts := time.Now().UnixNano()
			key := []byte(fmt.Sprintf("%d", ts))
			err := w.WriteMessages(ctx, kafka.Message{Key: key, Value: data})
			if err != nil {
				log.Println("could not write message " + err.Error())
			} else {
				log.Println("send message", ts)
			}
		}

	}
}

// helper function to anonimise our data record
func anonimise(alg, topic string, val []byte) ([]byte, error) {
	var data []byte
	var err error
	// after receiving the message, log its value
	if topic == "cmssw_pop_raw_metric" {
		var rec CMSSWRecord
		err = json.Unmarshal(val, &rec)
		if err != nil {
			return data, err
		}
		rec.Data.UserDN = hashFunc(rec.Data.UserDN, alg)
		data, err = json.Marshal(rec)
	}
	return data, err
}
