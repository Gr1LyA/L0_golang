package server

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"

	"github.com/Gr1LyA/L0_golang/internal/app/storage"
)

type server struct {
	store storage.ServerStorage
}

func NewServer() *server {
	return &server{}
}

func (s *server) startSrv(dbUrl string) error {
	log.Println("start")

	if err := s.configureStore(dbUrl); err != nil {
		return err
	}

	http.HandleFunc("/", s.midHandle("static/index.html"))
	return nil
}

func (s *server) midHandle(pagePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				http.ServeFile(w, r, pagePath)
			case "POST":
				if err := r.ParseForm(); err != nil {
					fmt.Fprintln(w, err)
				}
				if v, ok := s.store.Load(r.FormValue("uid")); ok {
					fmt.Fprint(w, v)
				} else if b, err := ioutil.ReadAll(r.Body) ; err == nil{ 
					if v, ok := s.store.Load(string(b)); ok {
						fmt.Fprint(w, v)
					} else {
						fmt.Fprint(w, "sorry, uid not found!")
					}
				} 
			default:
				fmt.Fprintln(w, "only get and post requests")
		}
	}
}

func (s *server) configureStore(dbUrl string) error {
	st := storage.New()

	if err := st.Open(dbUrl); err != nil {
		return err
	}

	s.store = st

	return nil
}