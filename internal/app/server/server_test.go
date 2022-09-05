package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gr1LyA/L0_golang/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

var (
	uid  = `b563feb7b2b84b6test`
	data = `{"order_uid":"b563feb7b2b84b6test","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"b563feb7b2b84b6test","request_id":"","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test","delivery_service":"meest","shardkey":"9","sm_id":99,"date_created":"2021-11-26T06:22:19Z","oof_shard":"1"}`
)

type testStorage struct {
	orders map[string]string
}

func (s *testStorage) Open(dbUrl string) error {
	return nil
}

func (s *testStorage) Load(key string) (string, bool) {
	val, ok := s.orders[key]
	return val, ok
}

func (s *testStorage) Store(key string, value string) error {
	s.orders[key] = value
	return nil
}

func (s *testStorage) Close() {

}

func NewTestStorage() storage.ServerStorage {
	return &testStorage{
		orders: make(map[string]string),
	}
}

func TestServer(t *testing.T) {
	s := NewServer()
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.midHandle("../../../static/index.html").ServeHTTP(rec, req)

	dataPage, err := ioutil.ReadFile("../../../static/index.html")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rec.Body.Bytes(), dataPage)

	//test POST
	s.store = NewTestStorage()
	s.store.Store(uid, data)

	rec = httptest.NewRecorder()

	r := strings.NewReader(uid)

	req, err = http.NewRequest(http.MethodPost, "/", r)
	if err != nil {
		t.Fatal(err)
	}

	s.midHandle("../../../static/index.html").ServeHTTP(rec, req)

	assert.Equal(t, rec.Body.String(), data)
}
