package main

import (
	"log"
	"my_service/internal/database"
	"my_service/internal/subscriber"
	"runtime"
)

//http.HandleFunc("/", handler) // each request calls handler
//log.Fatal(http.ListenAndServe("localhost:8000", nil))
//
//func handler(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
//}

func main() {

	//s := server.New()
	//if err := s.Start(); err != nil {
	//	log.Fatal(err)
	//}

	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	//str := ""
	//if err := db.Db.QueryRow("").Scan(str); err != nil {
	//	log.Fatal(err)
	//}
	go subscriber.SubscribeAndListen(db)

	runtime.Goexit()
	//fmt.Println("hello")
	//for {
	//	fmt.Println(db.Arr.Load("1"))
	//}
}
