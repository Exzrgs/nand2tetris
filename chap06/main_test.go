package main

import "testing"

func TestMain(m *testing.M) {
	symbolAddress = make(map[string]int, 0)

	m.Run()
}
