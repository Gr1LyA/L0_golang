package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)
import "fmt"

type DBStruct struct {
	Arr *mapMutex
	Db  *sql.DB
}

func PostgresDB() (*DBStruct, error) {
	//чтение конфига для подключени к БД
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbase, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", viper.Get("postgres.host"), viper.Get("postgres.port"), viper.Get("postgres.username"), viper.Get("postgres.password"), viper.Get("postgres.dbname"), viper.Get("postgres.sslmode")))

	if err != nil {
		return nil, err
	}

	err = dbase.Ping()
	if err != nil {
		return nil, err
	}

	//создание таблицы если она не создана
	dbase.QueryRow("CREATE TABLE IF NOT EXISTS orders(" +
		"uid text unique," +
		"data json);")

	//выгрузка из дб в кеш
	rows, err := dbase.Query("select * from orders")
	if err != nil {
		log.Println(err)
	}

	arr := NewCounters()

	for rows.Next() {
		var uid, data string

		err := rows.Scan(&uid, &data)
		if err != nil {
			log.Println(err)
		}
		arr.Store(uid, data)
	}

	rows.Close()

	return &DBStruct{
		Arr: arr,
		Db:  dbase,
	}, nil
}
