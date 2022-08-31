package main

import (
	"fmt"
	"log"
	"my_service/internal/database"
	"my_service/internal/subscriber"
	"net/http"
	"runtime"
)

func main() {

	//s := server.New()
	//if err := s.Start(); err != nil {
	//	log.Fatal(err)
	//}

	//создание бд и получение структуры
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	//подключение к nats и подписка на канал
	go subscriber.SubscribeAndListen(db)

	//сервер
	//handler для сервера
	mainPage := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "static/index.html")
		case "POST":
			if err := r.ParseForm(); err != nil {
				fmt.Fprintln(w, err)
			}

			//fmt.Fprintf(w, "Post from website r.postform = %v\n", r.PostForm)
			if v, ok := db.Arr.Load(r.FormValue("uid")); ok {
				fmt.Fprintln(w, v)
			} else {
				fmt.Fprintln(w, "sorry, uid not found!")
			}
		default:
			fmt.Fprintln(w, "only get and post requests")
		}
	}
	http.HandleFunc("/", mainPage)

	//запуск сервера
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}
