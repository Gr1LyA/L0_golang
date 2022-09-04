package main

import (
	"fmt"
	"log"
	. "my_service/internal/database"
	. "my_service/internal/server"
	. "my_service/internal/subscriber"
	"os"
	"os/signal"
)

func main() {
	//создание бд и получение структуры
	log.Println("подключение к бд")
	db, err := NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	//подключение к nats и подписка на канал
	log.Println("подкючение к nats")
	closeStruct := SubscribeAndListen(db)

	log.Println("запуск http сервера")
	StartServer(db)

	waitAndClose(closeStruct)
}

func waitAndClose(closeStruct StructForClose) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
	closeStruct.Sub.Unsubscribe()
	closeStruct.Sc.Close()
	closeStruct.Nc.Close()
	closeStruct.Db.Close()
}
