package main

import "testing"

func TestAStar(t *testing.T) {
	b := Board{
		Width:  3,
		Height: 3,
		Food:   []Coord{Coord{0, 0}},
		Snakes: []Battlesnake{Battlesnake{Head: Coord{2, 2}}},
	}

	dist, _ := runAStar(&b, 0, 0)
	if dist != 4 {
		t.Errorf("Expected: %d, Got: %d", dist, 4)
	}
}

var result int8

func BenchmarkAStar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := Board{
			Width:  11,
			Height: 11,
			Food:   []Coord{Coord{0, 0}},
			Snakes: []Battlesnake{
				Battlesnake{Head: Coord{2, 2}},
				Battlesnake{Head: Coord{3, 3}},
				Battlesnake{Head: Coord{10, 10}},
				Battlesnake{Head: Coord{10, 0}},
			},
		}
		result, _ = runAStar(&b, 2, 0)
	}
}
