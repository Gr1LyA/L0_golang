package storage

import (
	"testing"
	"github.com/spf13/viper"
	"fmt"
)

var (
	uid = `b563ftest`
	data = `{"order_uid":"b563feb7b2b84b6test","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"b563feb7b2b84b6test","request_id":"","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test","delivery_service":"meest","shardkey":"9","sm_id":99,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}`
)

func TestStore(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}

	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
		viper.Get("postgres.host"), viper.Get("postgres.port"), viper.Get("postgres.username"), 
		viper.Get("postgres.password"), viper.Get("postgres.dbname"), viper.Get("postgres.sslmode"))

	s := New()

	if err := s.Open(dbUrl); err != nil {
		t.Fatal(err)
	}

	if err = s.Store(uid, data); err != nil {
		t.Fatal(err)
	}

	s.db.Query("delete from orders where uid='b563feb7b2b84b6test';")

	s.Close()
}