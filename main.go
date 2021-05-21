package main

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"flag"
	"hash"
	"log"
	"strconv"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func produce(ctx context.Context) {
	// initialize a counter
	i := 0

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: Config.Brokers,
		Topic:   Config.SendTopic,
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
	// initialize a new reader with the brokers and receive topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  Config.Brokers,
		Topic:    Config.ReceiveTopic,
		GroupID:  Config.Group,
		MinBytes: 5,
		MaxBytes: 1e6,
		// wait for at most 3 seconds before receiving new data
		MaxWait: 3 * time.Second,
		// this will start consuming messages from the earliest available
		StartOffset: kafka.FirstOffset,
		// if you set it to `kafka.LastOffset` it will only consume new messages
	})
	log.Println("kafak reader", r)

	// intialize the writer with the broker addresses, and the send topic
	var w *kafka.Writer
	if Config.SendTopic != "" {
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers: Config.Brokers,
			Topic:   Config.SendTopic,
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
	}
	// initialize a counter
	i := 0

	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("could not read message " + err.Error())
			continue
		}
		log.Println("received: ", string(msg.Value))
		// after receiving the message, log its value
		if Config.ReceiveTopic == "cmssw_pop_raw_metric" {
			// convert message to our struct
			var rec CMSSWRecord
			err = json.Unmarshal(msg.Value, &rec)
			if err != nil {
				log.Println("unable to convert message to struct", err)
			}
			rec.Data.UserDN = hashFunc(rec.Data.UserDN)
			data, err := json.Marshal(rec)
			if err != nil {
				log.Println("unable to marshal the record", err)
			}
			log.Println("new record: ", string(data))
			if Config.SendTopic != "" {
				key := []byte(strconv.Itoa(i))
				err := w.WriteMessages(ctx, kafka.Message{Key: key, Value: data})
				if err != nil {
					log.Println("could not write message " + err.Error())
				} else {
					log.Println("send message", i)
					i++ // increment counter if we send successful message
				}
			}
		}
	}
}

// helper function t convert given value to a hash one
func hashFunc(val string) string {
	var h hash.Hash
	sha := strings.ToLower(Config.SHA)
	if sha == "sha256" {
		h = sha256.New()
	} else if sha == "sha512" {
		h = sha512.New()
	} else {
		h = sha1.New()
	}
	h.Write([]byte(val))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	var config string
	flag.StringVar(&config, "config", "config.json", "config file")
	flag.Parse()
	log.SetFlags(0)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := ParseConfig(config)
	if err != nil {
		log.Fatalf("fail to parse config file %v", err)
	}

	log.Println("config", Config.String())
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
