// Copyright 2016-2019 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"log"

	"github.com/nats-io/stan.go"
)


func main() {
	sc, err := stan.Connect("test-cluster", "stan-sub")
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
