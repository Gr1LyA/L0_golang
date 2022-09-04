package server

import (
	"fmt"
	"log"
	. "my_service/internal/database"
	"net/http"
)

func StartServer(db *DBStruct) {
	//handler для сервера

	mainPage := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "static/index.html")
		case "POST":
			if err := r.ParseForm(); err != nil {
				fmt.Fprintln(w, err)
			}

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
	go func() {
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
