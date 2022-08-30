package subscriber

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	. "my_service/internal/database"
	"os"
	"os/signal"
	"time"
)

var usageStr = `
Usage: stan-sub [options] <subject>

Options:
	-s,  --server   <url>            NATS Streaming server URL(s)
	-c,  --cluster  <cluster name>   NATS Streaming cluster name
	-id, --clientid <client ID>      NATS Streaming client ID
`

func usage() {
	log.Fatalf(usageStr)
}

func SubscribeAndListen(db *DBStruct) {
	var (
		clusterID, clientID string
		URL                 string
		unsubscribe         bool
		durable             string
	)
	flag.StringVar(&URL, "s", stan.DefaultNatsURL, "The nats server URLs (separated by comma)")
	flag.StringVar(&URL, "server", stan.DefaultNatsURL, "The nats server URLs (separated by comma)")
	flag.StringVar(&clusterID, "c", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clusterID, "cluster", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-sub", "The NATS Streaming client ID to connect with")
	flag.StringVar(&clientID, "clientid", "stan-sub", "The NATS Streaming client ID to connect with")
	flag.BoolVar(&unsubscribe, "unsub", false, "Unsubscribe the durable on exit")
	flag.BoolVar(&unsubscribe, "unsubscribe", false, "Unsubscribe the durable on exit")
	flag.StringVar(&durable, "durable", "", "Durable subscriber name")

	log.SetFlags(0)
	flag.Parse()

	//получение аргументов коммандной строки
	args := flag.Args()
	if len(args) < 1 {
		log.Printf("Error: A subject must be specified.")
		usage()
	}

	//настройка опций подключения
	opts := []nats.Option{nats.Name("NATS Publisher")} // Connect to NATS
	opts = setupConnOptions(opts)

	//подключение к серверу
	nc, err := nats.Connect(URL, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect(clusterID, clientID, stan.NatsOptions(opts...))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", URL, clusterID, clientID)

	subj, i := args[0], 0

	mcb := func(msg *stan.Msg) {
		i += 1
		parserMsg(msg, i, db)
	}

	//Подписка на канал
	sub, err := sc.Subscribe(subj, mcb)
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s], durable=[%s]\n", subj, clientID, durable)

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			if durable == "" || unsubscribe {
				sub.Unsubscribe()
			}
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func parserMsg(m *stan.Msg, i int, db *DBStruct) {
	//log.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
	var jsn modelJson

	//проверка json на валидность
	if !json.Valid(m.Data) {
		log.Println("invalid json")
		return
	}

	//распаковка json в структуру
	err := json.Unmarshal(m.Data, &jsn)
	if err != nil {
		log.Fatal(err)
	}

	if jsn.OrderUID == "" {
		log.Println("invalid uid")
		return
	}
	db.Arr.Store(jsn.OrderUID, string(m.Data))
	fmt.Println(jsn.OrderUID)

	//fmt.Println(db.Arr.Load(jsn.OrderUID))

	//Запись в бд
	db.Db.QueryRow("insert into orders (uid, data) values ($1, $2)", jsn.OrderUID, string(m.Data))
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
