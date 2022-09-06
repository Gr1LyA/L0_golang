package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"context"
	"log"

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

	srvRun := &http.Server{Addr: ":8080"}

	signalChan := make(chan os.Signal)
	cleanupDone := make(chan struct{})

	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		log.Printf("Closing connections...\n\n")
		srv.Close()
		log.Printf("Closed\n\n")
		srvRun.Shutdown(context.TODO())
		close(cleanupDone)
	}()

	if err := srvRun.ListenAndServe(); err != http.ErrServerClosed {
		signalChan <- os.Interrupt
		<- cleanupDone
		return err
	}
	<- cleanupDone
	return nil
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