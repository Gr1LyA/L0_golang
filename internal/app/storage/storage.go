package storage

import (
	"database/sql"
	"log"
	"encoding/json"

	_ "github.com/lib/pq"
	"github.com/Gr1LyA/L0_golang/internal/app/model"
	"github.com/go-playground/validator/v10"
)

type ServerStorage interface {
	Open(string) error
	Load(string) (string, bool)
	Store(string, string) error
	Close()
}

type storage struct {
	db *sql.DB
	orders *model.MapMutex
}

func New() *storage {
	return &storage{}
}

func (s *storage) Open(dbUrl string) error {

	
	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		return err
	}

	
	if err = db.Ping(); err != nil {
		return err
	}

	s.db = db

	if err = s.loadCache(); err != nil {
		return err
	}

	return nil
}

func (s *storage) loadCache() error {
	err := s.db.QueryRow(
		"CREATE TABLE IF NOT EXISTS orders(" +
		"uid text unique," +
		"data json);").Err()
	if err != nil {
		return err
	}

	rows, err := s.db.Query("select * from orders")
	if err != nil {
		return err
	}

	orders := model.NewRWMap()

	for rows.Next() {
		var uid, data string

		err := rows.Scan(&uid, &data)
		if err != nil {
			return err
		}
		orders.Store(uid, data)
	}

	rows.Close()

	s.orders = orders

	return nil
}

func (s *storage) Load(key string) (string, bool){
	val, ok := s.orders.Load(key)
	return val, ok
}

func (s *storage) Store(key string, value string) error {
	if !validOrders(value) {
		return nil
	}

	if err := s.db.QueryRow("insert into orders (uid, data) values ($1, $2)", key, value).Err(); err != nil {
		return err
	}

	s.orders.Store(key, value)
	log.Println("add: ", key)

	return nil
}

func validOrders(value string) bool {
	var jsonData model.ModelOrder
	
	if !json.Valid([]byte(value)) {
		log.Println("invalid json")
		return false
	}

	if err := json.Unmarshal([]byte (value), &jsonData); err != nil {
		log.Println(err)
		return false
	}

	validate := validator.New()
	err := validate.Struct(jsonData)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (s *storage) Close() {
	s.db.Close()
}