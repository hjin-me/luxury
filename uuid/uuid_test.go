package uuid

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	id := New()
	t.Log(id)
	if id.Time().Sub(time.Now()) > time.Second {
		t.Log(id.Time())
		t.Error("uuid time is not current time")
	}
}

func TestSetTime(t *testing.T) {

	id := New()
	fiveMinutesLater := time.Now().Add(5 * time.Minute)
	id.SetTime(fiveMinutesLater)
	if id.Time().Unix() != fiveMinutesLater.Unix() {
		t.Log("id time", id.Time())
		t.Log("ta time", fiveMinutesLater)
		t.Error("time not equal")
	}
}
