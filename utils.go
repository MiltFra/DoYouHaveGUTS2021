package main

import "math"

func ManhattanDistance(a, b Coord) int8 {
	res := int8(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
	return res
}
