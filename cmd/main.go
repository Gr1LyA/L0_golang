package main

import (
	"fmt"
	"log"
	"my_service/internal/database"
	"my_service/internal/server"
	"runtime"
)

func main() {
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	if val, ok := db.Arr.Load("1"); ok {
		fmt.Println(val)
	}

	s := server.New()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	//subscriber.SubscribeAndListen(db)

	runtime.Goexit()
}
