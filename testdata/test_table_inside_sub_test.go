package main

import (
	"fmt"
	"testing"
)

func TestFunction(t *testing.T) {
	t.Run("Test Error Handling", func(t *testing.T) {
		var tt = []struct {
			nameNested string
			in         int
			out        int
		}{
			{
				nameNested: "inc(3) should be 4",
				in:         3,
				out:        4,
			},
		}

		for _, tc := range tt {
			t.Run(tc.nameNested, func(t *testing.T) {
				fmt.Println(tc.nameNested)
			})
		}
	})
}
