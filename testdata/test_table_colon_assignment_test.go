package main

import "testing"

func TestAdd(t *testing.T) {
	tt := []struct {
		name string
		op1  int
		op2  int
		out  int
	}{
		{
			name: "1 + 1 should be 2",
			op1:  1,
			op2:  1,
			out:  2,
		},
		{
			name: "1 + 2 should be 3",
			op1:  1,
			op2:  2,
			out:  3,
		},
		{
			name: "2 + 2 should be 3",
			op1:  2,
			op2:  2,
			out:  3,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := add(tc.op1, tc.op2)
			if actual != tc.out {
				t.Errorf("%d + %d should be equal to %d, got %d instead", tc.op1, tc.op2, tc.out, actual)
			}
		})
	}
}

func TestInc(t *testing.T) {
	tt := []struct {
		name string
		in   int
		out  int
	}{
		{
			name: "inc(1) should be 2",
			in:   1,
			out:  2,
		},
		{
			name: "inc(-1) should be 0",
			in:   -1,
			out:  3,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := inc(tc.in)
			if actual != tc.out {
				t.Errorf("inc(%d) should be equal to %d, got %d instead", tc.in, tc.out, actual)
			}
		})
	}
}
