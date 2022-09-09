package stan

import (
	"testing"

	"github.com/Gr1LyA/L0_golang/internal/app/model"
)

func TestStan(t *testing.T) {
	var s StanStruct

	if err := s.ConnectAndSubscribe(model.NewTestStorage()); err != nil {
		t.Fatal(err)
	}

	s.Close()
}
