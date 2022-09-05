package server

import (
	"io"
	"log"
	"net/http"
	"fmt"

	"github.com/spf13/viper"
	"github.com/Gr1LyA/L0_golang/internal/app/storage"
)

type Server struct {
	store *storage.Storage
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	log.Println("start")

	if err := s.configureStore(); err != nil {
		return err
	}

	http.HandleFunc("/", s.midHandle())
	return http.ListenAndServe("localhost:8080", nil)
}

func (s *Server) midHandle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "alo")
		// http.ServeFile(w, r, "static/index.html")
	}
}

func (s *Server) configureStore() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
		viper.Get("postgres.host"), viper.Get("postgres.port"), viper.Get("postgres.username"), 
		viper.Get("postgres.password"), viper.Get("postgres.dbname"), viper.Get("postgres.sslmode"))

	st := storage.New()

	if err := st.Open(dbUrl); err != nil {
		return err
	}

	s.store = st

	return nil
}