package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)
import "fmt"

type infoDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type DBStruct struct {
	Arr *mapMutex
	Db  *sql.DB
}

func NewPostgresDB() (*DBStruct, error) {
	infoConnection := infoDB{
		Host:     "localhost",
		Port:     "5437",
		Username: "ilya",
		Password: "4774",
		DBName:   "wb",
	}
	//dbase, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", infoConnection.Host, infoConnection.Port, infoConnection.Username, infoConnection.Password, infoConnection.DBName, "disable"))

	dbase, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", infoConnection.Host, infoConnection.Username, infoConnection.Password, infoConnection.DBName, "disable"))

	if err != nil {
		return nil, err
	}

	err = dbase.Ping()
	if err != nil {
		return nil, err
	}

	return &DBStruct{
		Arr: NewCounters(),
		Db:  dbase,
	}, nil
}
