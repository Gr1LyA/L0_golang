package stan

import (
	"github.com/nats-io/stan.go"
	"github.com/Gr1LyA/L0_golang/internal/app/storage"
	"encoding/json"
	"fmt"
)

type StanStruct struct {
	sub stan.Subscription
	sc  stan.Conn
}

func (st *StanStruct) ConnectAndSubscribe(store storage.ServerStorage) error {
	sc, err := stan.Connect("test-cluster", "stan-sub")
	if err != nil {
		return err
	}

	mcb := func(msg *stan.Msg) {
		parserMsg(msg, store)
	}

	sub, err := sc.Subscribe("json-receive", mcb)
	if err != nil {
		return err
	}

	st.sub = sub
	st.sc = sc
	return nil
}

func parserMsg(m *stan.Msg, store storage.ServerStorage) {
	var jsn struct {
		OrderUID string `json:"order_uid"`
	}
	
	if err := json.Unmarshal(m.Data, &jsn); err != nil {
		fmt.Println(err)
		return
	}

	if err := store.Store(jsn.OrderUID, string(m.Data)); err != nil {
		fmt.Println(err)
	}
}

func (st *StanStruct) Close() {
	st.sub.Close()
	st.sc.Close()
}