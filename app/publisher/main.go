package main

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	//получение аргументов коммандной строки
	args := os.Args

	log.SetFlags(0)

	if len(args) != 3 {
		log.Fatal("wrong count arguments\nexpected 2 arg")
	}

	opts := []nats.Option{nats.Name("NATS Publisher")}

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL, opts...)
	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	subj, msg := args[1], []byte(args[2])

	if err := nc.Publish(subj, msg); err != nil {
		log.Fatal(err)
	}

	if err := nc.Flush(); err != nil {
		log.Fatal(err)
	}

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}
}
