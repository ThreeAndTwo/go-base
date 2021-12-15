package redis

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

var mocker *Mocker

func init() {
	var err error
	mocker, err = NewMocker()
	if err != nil {
		panic(err)
	}
}
func TestMockerGetAndSet(t *testing.T) {
	if err := mocker.Set("test", "hello", 1 * time.Minute); err != nil {
		t.Fatal(err)
	}
	got, err := mocker.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, got, "hello")
}
