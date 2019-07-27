package main

import "testing"

func TestAddVarAssignment(t *testing.T) {
	var _, tt = 1, []struct {
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

func TestSimpleVarAssignment(t *testing.T) {
	var tt = []struct {
		name string
		op1  int
		op2  int
		out  int
	}{
		{
			name: "simple var 1 + 1 should be 2",
			op1:  1,
			op2:  1,
			out:  2,
		},
		{
			name: "simple var 1 + 2 should be 3",
			op1:  1,
			op2:  2,
			out:  3,
		},
		{
			name: "simple var 2 + 2 should be 3",
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
