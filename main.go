package main

import (
	"context"
	"flag"
	"log"
	"strings"
)

func main() {
	var brokers string
	flag.StringVar(&brokers, "brokers", "", "comma separated brokers, e.g. host1:port1,host2:port2")
	var sendTopic string
	flag.StringVar(&sendTopic, "sendTopic", "", "kafka send topic")
	var receiveTopic string
	flag.StringVar(&receiveTopic, "receiveTopic", "", "kafka receive topic")
	var testTopic string
	flag.StringVar(&testTopic, "testTopic", "", "kafka test topic")
	var group string
	flag.StringVar(&group, "group", "test-group", "kafka group name")
	var alg string
	flag.StringVar(&alg, "alg", "sha1", "sha algorithm (sha1, sha256, sha512) to use for anonimisation")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbosity level")
	flag.Parse()
	log.SetFlags(0)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// parse brokers
	var brks []string
	for _, b := range strings.Split(brokers, ",") {
		brks = append(brks, strings.Trim(b, " "))
	}

	// create a new context
	ctx := context.Background()

	// perform action
	if testTopic != "" {
		// produce messages in a new go routine, since
		// both the produce and consume functions are blocking
		go produce(ctx, brks, testTopic, group)
		consume(ctx, brks, testTopic, group)
	} else {
		process(ctx, brks, receiveTopic, sendTopic, group, alg, verbose)
	}
}
