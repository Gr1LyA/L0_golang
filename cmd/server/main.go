package main

import (
	"log"
	"github.com/Gr1LyA/L0_golang/internal/app/server"
)

func main() {
	s := server.New()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}