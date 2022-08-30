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
	SSLMode  string
}

type DBStruct struct {
	Arr *mapMutex
	Db  *sql.DB
}

func NewPostgresDB() (*DBStruct, error) {
	infoConnection := infoDB{
		Host:     "localhost",
		Port:     "5432",
		Username: "ilya",
		Password: "4774",
		DBName:   "wb",
		SSLMode:  "disable",
	}

	dbase, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", infoConnection.Host, infoConnection.Port, infoConnection.Username, infoConnection.Password, infoConnection.DBName, infoConnection.SSLMode))

	//dbase.SetConnMaxLifetime(time.Minute)

	if err != nil {
		return nil, err
	}

	err = dbase.Ping()
	if err != nil {
		return nil, err
	}

	//создание таблицы если она не создана
	dbase.QueryRow("CREATE TABLE IF NOT EXISTS orders(" +
		"uid text," +
		"data json);")

	return &DBStruct{
		Arr: NewCounters(),
		Db:  dbase,
	}, nil
}
