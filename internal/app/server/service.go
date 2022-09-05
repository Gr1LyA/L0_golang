package server

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func Start() error {
	if err := parseConfig(); err != nil {
		return err
	}

	srv := NewServer()
	if err := srv.startSrv(databaseURL()); err != nil {
		return err
	}
	
	return http.ListenAndServe(":8080", nil)
}

func parseConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func databaseURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
		viper.Get("postgres.host"), viper.Get("postgres.port"), viper.Get("postgres.username"), 
		viper.Get("postgres.password"), viper.Get("postgres.dbname"), viper.Get("postgres.sslmode"),
	)
}