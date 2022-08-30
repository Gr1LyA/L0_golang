package subscriber

import (
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"

	"github.com/nats-io/stan.go"
	. "my_service/internal/database"
)

func SubscribeAndListen(db DBStruct) {
	var urls = flag.String("s", stan.DefaultNatsURL, "The nats server URLs")

	log.SetFlags(0)
	flag.Parse()

	opts := []nats.Option{nats.Name("NATS Publisher")} // Connect to NATS
	opts = setupConnOptions(opts)

	//подключение к серверу
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	//получение аргументов коммандной строки
	args := flag.Args()

	subj, i := args[0], 0

	//Подписка на канал
	nc.Subscribe(subj, func(msg *nats.Msg) {
		i += 1
		parserMsg(msg, i, db)
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", subj)

}

func parserMsg(m *nats.Msg, i int, db DBStruct) {
	log.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
	db.Arr.Store("1", string(m.Data))
	fmt.Println(db.Arr.Load("1"))
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}
