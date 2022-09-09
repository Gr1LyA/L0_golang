package main

import (
	"github.com/Gr1LyA/L0_golang/internal/app/server"
	"log"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
