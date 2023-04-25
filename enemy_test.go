package main

import "testing"

func TestRandom(t *testing.T) {
	min := 1
	max := 5
	m := map[int]int{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
	}
	result := make(map[int]int, len(m))
	for len(m) == 0 {
		key := random(min, max)
		if v, ok := m[key]; ok {
			result[key] = v
			delete(m, key)
		} else {
			if _, ok := result[key]; !ok {
				t.Error()
			}
		}
	}
}
