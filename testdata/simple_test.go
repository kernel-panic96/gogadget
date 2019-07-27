package main

import (
	"testing"
	aliasedTesting "testing"
)

import aliasedTwo "testing"

func TestSimple(t *aliasedTesting.T) {
	t.Run("simple test", func(t *testing.T) {
		t.Log("doing work")
	})
	t.Run("second test", func(t *testing.T) {
		t.Log("Doing work twice")
	})
	_ = aliasedTwo.Cover{}
}
