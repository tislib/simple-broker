package main

import "testing"

func TestAdd(t *testing.T) {
	actual := Add(3, 5)
	expected := 8

	if actual != expected {
		t.Fail()
	}
}
