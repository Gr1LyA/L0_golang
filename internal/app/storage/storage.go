package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/Gr1LyA/L0_golang/internal/app/model"
)

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
	if err := s.db.QueryRow("insert into orders (uid, data) values ($1, $2)", key, value).Err(); err != nil {
		return err
	}
	s.orders.Store(key, value)
	return nil
}

func (s *storage) Close() {
	s.db.Close()
}