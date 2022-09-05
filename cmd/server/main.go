package main

import (
	"log"
	"github.com/Gr1LyA/L0_golang/internal/app/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}