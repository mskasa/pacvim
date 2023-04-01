package main

import (
	"reflect"
	"testing"
)

func Test_squaresAround(t *testing.T) {
	expected := [][]int{
		{0, 0, 0}, {-1, 0, 1}, {1, 0, 1}, {0, 1, 1}, {0, -1, 1},
		{-1, -1, 2}, {0, 2, 2}, {2, 0, 2}, {-1, 1, 2}, {1, 1, 2}, {0, -2, 2}, {1, -1, 2}, {-2, 0, 2},
		{-2, 1, 3}, {1, -2, 3}, {-2, -1, 3}, {-1, 2, 3}, {1, 2, 3}, {2, -1, 3}, {-1, -2, 3}, {2, 1, 3},
		{-2, -2, 4}, {-2, 2, 4}, {2, -2, 4}, {2, 2, 4}}
	result := squaresAround(2)

	if !reflect.DeepEqual(expected, result) {
		t.Error("expected:", expected, "result:", result)
	}
}
