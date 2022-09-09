package main

import (
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "stan-pub")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	args := os.Args

	if len(args) != 2 {
		panic("expected only message")
	}

	subj := "json-receive"
	msg := []byte(args[1])

	err = sc.Publish(subj, msg)
	if err != nil {
		log.Fatalf("Error during publish: %v\n", err)
	}
	log.Printf("Published [%s] : '%s'\n", subj, msg)
}
