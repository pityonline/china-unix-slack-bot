package service

import "testing"

func TestPing(t *testing.T) {
	want := "pong"
	got := Ping()

	if got != want {
		t.Errorf("got: %#v, want: %#v", got, want)
	}
}
