package main

import "testing"

func TestNested(t *testing.T) {
	t.Run("base test", func(t *testing.T) {
		t.Run("nested1", func(t *testing.T) {
			t.Log("doing nested work1")
		})
		t.Run("nested2", func(t *testing.T) {
			t.Log("doing nested work2")
		})
	})
}
