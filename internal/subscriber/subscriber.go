package subscriber

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	. "my_service/internal/database"
)

var usageStr = `
Usage: stan-sub [options] <subject>

Options:
	-s,  --server   <url>            NATS Streaming server URL(s)
	-c,  --cluster  <cluster name>   NATS Streaming cluster name
	-id, --clientid <client ID>      NATS Streaming client ID
`

type modelJson struct {
	OrderUID string `json:"order_uid"`
}

type StructForClose struct {
	Sub stan.Subscription
	Sc  stan.Conn
	Nc  *nats.Conn
	Db  *sql.DB
}

func SubscribeAndListen(db *DBStruct) StructForClose {
	var (
		clusterID, clientID string
		URL                 string
		durable             string
	)
	flag.StringVar(&URL, "s", stan.DefaultNatsURL, "The nats server URLs (separated by comma)")
	flag.StringVar(&URL, "server", stan.DefaultNatsURL, "The nats server URLs (separated by comma)")
	flag.StringVar(&clusterID, "c", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clusterID, "cluster", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-sub", "The NATS Streaming client ID to connect with")
	flag.StringVar(&clientID, "clientid", "stan-sub", "The NATS Streaming client ID to connect with")

	//log.SetFlags(0)
	flag.Parse()

	//получение аргументов коммандной строки
	args := flag.Args()
	if len(args) < 1 {
		db.Db.Close()
		log.Fatal(usageStr)
	}

	opts := []nats.Option{nats.Name("NATS Publisher")}

	//подключение к серверу
	nc, err := nats.Connect(URL, opts...)
	if err != nil {
		log.Fatal(err)
	}

	sc, err := stan.Connect(clusterID, clientID, stan.NatsOptions(opts...))
	if err != nil {
		log.Fatal(err)
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
		nc.Close()
		sc.Close()
		db.Db.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s], durable=[%s]\n", subj, clientID, durable)

	return StructForClose{
		Sub: sub,
		Sc:  sc,
		Nc:  nc,
		Db:  db.Db,
	}
}

func parserMsg(m *stan.Msg, i int, db *DBStruct) {
	var jsn modelJson

	//проверка json на валидность
	if !json.Valid(m.Data) {
		log.Println("invalid json")
		return
	}

	//распаковка json в структуру
	err := json.Unmarshal(m.Data, &jsn)
	if err != nil {
		log.Println(err)
		return
	}

	if jsn.OrderUID == "" {
		log.Println("invalid uid")
		return
	}
	db.Arr.Store(jsn.OrderUID, string(m.Data))
	fmt.Println("add: ", jsn.OrderUID)

	//Запись в бд
	db.Db.QueryRow("insert into orders (uid, data) values ($1, $2)", jsn.OrderUID, string(m.Data))
}
