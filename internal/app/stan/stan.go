package stan

import (
	"github.com/nats-io/stan.go"
	"github.com/Gr1LyA/L0_golang/internal/app/storage"
	"github.com/Gr1LyA/L0_golang/internal/app/model"
	"encoding/json"
	"fmt"
)

type stanStruct struct {
	sub stan.Subscription
	sc  stan.Conn
}

func New() *stanStruct {
	return &stanStruct{}
}

func (st *stanStruct) connectAndSubscribe(store storage.ServerStorage) error {
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
	var jsn model.ModelJson

	//проверка json на валидность
	if !json.Valid(m.Data) {
		fmt.Println("invalid json")
		return
	}

	//распаковка json в структуру
	err := json.Unmarshal(m.Data, &jsn)
	if err != nil {
		fmt.Println(err)
		return
	}

	if jsn.OrderUID == "" {
		fmt.Println("invalid uid")
		return
	}

	err = store.Store(jsn.OrderUID, string(m.Data))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("add: ", jsn.OrderUID)
}

func (st *stanStruct) close(){
	st.sub.Close()
	st.sc.Close()
}