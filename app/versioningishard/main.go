package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/inlined/versioningishard/lib/cloudevents"
)

var (
	dataLoc = flag.String("data", "", "Path of event to diagnose")
)

func main() {
	flag.Parse()

	if *dataLoc == "" {
		log.Fatal("Must specify the path to a JSON event file with --data")
	}

	b, err := ioutil.ReadFile(*dataLoc)
	if err != nil {
		log.Fatalf("Failed to read file %q: %s", *dataLoc, err)
	}

	var event cloudevents.CloudEvent
	if err := json.Unmarshal(b, &event); err != nil {
		log.Fatalf("Failed to parse event at file %q: %s", *dataLoc, err)
	}

	log.Println("Got event", event.EventID)
	sampledRateRaw := event.Extensions["sampledRate"]
	// (Our custom library for inline attributes only supports strings)
	if sampledRate, ok := sampledRateRaw.(string); ok {
		log.Printf("(It was sampled at a rate of 1 in %s)\n", sampledRate)
	} else {
		log.Println("(It was not sampled)")
	}
}
