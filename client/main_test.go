package main

import (
	"sync"
	"testing"
)

type A struct {
	a *sync.Map
}

func TestA(t *testing.T) {
	a := A{}
	a.a = &sync.Map{}
	a.a.Store("a", 1)
	a.a = &sync.Map{}
	a.a.Range(func(key, value any) bool {
		t.Log(key, value)
		return true
	})
}
