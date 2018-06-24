package service

import "testing"

func TestGreet(t *testing.T) {
	want := "Hello world!"
	got := Greet()

	if got != want {
		t.Errorf("got: %#v, want: %#v", got, want)
	}
}
