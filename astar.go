package main

import (
	"container/heap"
	"fmt"
	"sort"
)

type AStarEntry struct {
	dist       int8
	cost       int8
	pos        Coord
	first_move MoveType
}

func NewAStarEntry(cur, tar Coord, cost int8, first_move MoveType) *AStarEntry {
	return &AStarEntry{
		dist:       ManhattanDistance(cur, tar),
		cost:       cost,
		pos:        cur,
		first_move: first_move,
	}
}

type AStarHeap []AStarEntry

func (h AStarHeap) Len() int {
	return len(h)
}

func (h AStarHeap) Less(i, j int) bool {
	return h[i].dist+h[i].cost < h[j].dist+h[j].cost
}

func (h AStarHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *AStarHeap) Push(x interface{}) {
	*h = append(*h, x.(AStarEntry))
}

func (h *AStarHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func makeGrid(board *Board) []bool {
	grid := make([]bool, board.Height*board.Width)
	var s Battlesnake
	for i := range board.Snakes {
		s = board.Snakes[i]
		grid[s.Head.Y*board.Width+s.Head.X] = true
		for j := range s.Body {
			grid[s.Body[j].Y*board.Width+s.Body[j].X] = true
		}
	}
	return grid
}

func extendGrid(board *Board, grid []bool, you *Battlesnake) []bool {
	moves := func() map[Coord]struct{} {
		pos := you.Head
		moves := make(map[Coord]struct{})
		if pos.X > 0 {
			moves[Coord{X: pos.X - 1, Y: pos.Y}] = struct{}{}
		}
		if pos.X < board.Width-1 {
			moves[Coord{X: pos.X + 1, Y: pos.Y}] = struct{}{}
		}
		if pos.Y > 0 {
			moves[Coord{X: pos.X, Y: pos.Y - 1}] = struct{}{}
		}
		if pos.Y < board.Height-1 {
			moves[Coord{X: pos.X, Y: pos.Y + 1}] = struct{}{}
		}
		return moves
	}()
	handleMove := func(pos Coord) {
		_, ok := moves[pos]
		if ok {
			fmt.Printf("Avoided loosing head-to-head\n")
			grid[board.Width*pos.Y+pos.X] = true
		}
	}
	for i := range board.Snakes {
		if board.Snakes[i].ID == you.ID {
			continue
		}
		if board.Snakes[i].Length >= you.Length {
			pos := board.Snakes[i].Head
			if pos.X > 0 {
				handleMove(pos.Left())
			}
			if pos.X < board.Width-1 {
				handleMove(pos.Right())
			}
			if pos.Y > 0 {
				handleMove(pos.Down())
			}
			if pos.Y < board.Height-1 {
				handleMove(pos.Up())
			}
		}
	}
	return grid
}

func makeHeap(board *Board) AStarHeap {
	return make(AStarHeap, 0, board.Height*board.Width)
}

func possibleMoves(g, visited []bool, h, w int8, cur AStarEntry, tar Coord) [4]*AStarEntry {
	var moves [4]*AStarEntry
	// Down
	if cur.pos.Y > 0 {
		k := w*(cur.pos.Y-1) + cur.pos.X
		if !visited[k] && !g[k] {
			visited[k] = true
			moves[0] = NewAStarEntry(
				Coord{X: cur.pos.X, Y: cur.pos.Y - 1},
				tar,
				cur.cost+1,
				Down,
			)
		}
	}
	// Up
	if cur.pos.Y < h-1 {
		k := w*(cur.pos.Y+1) + cur.pos.X
		if !visited[k] && !g[k] {
			visited[k] = true
			moves[1] = NewAStarEntry(
				Coord{X: cur.pos.X, Y: cur.pos.Y + 1},
				tar,
				cur.cost+1,
				Up,
			)
		}
	}
	// Left
	if cur.pos.X > 0 {
		k := w*cur.pos.Y + cur.pos.X - 1
		if !visited[k] && !g[k] {
			visited[k] = true
			moves[2] = NewAStarEntry(
				Coord{X: cur.pos.X - 1, Y: cur.pos.Y},
				tar,
				cur.cost+1,
				Left,
			)
		}
	}
	// Right
	if cur.pos.X < w-1 {
		k := w*cur.pos.Y + cur.pos.X + 1
		if !visited[k] && !g[k] {
			visited[k] = true
			moves[3] = NewAStarEntry(
				Coord{X: cur.pos.X + 1, Y: cur.pos.Y},
				tar,
				cur.cost+1,
				Right,
			)
		}
	}
	return moves
}

func runAStarWithGrid(board *Board, snake, food int, g []bool) (int8, MoveType) {
	h := board.Height
	w := board.Width
	q := makeHeap(board)
  visited := make([]bool, h*w)
	tar := board.Food[food]
	start := board.Snakes[snake].Head
	visited[w*start.Y+start.X] = true
	cur := *NewAStarEntry(start, tar, 0, Up)
	moves := possibleMoves(g, visited, h, w, cur, tar)
	for i := range moves {
		if moves[i] != nil {
			heap.Push(&q, *moves[i])
		}
	}
	for len(q) > 0 {
		cur = q.Pop().(AStarEntry)
		// fmt.Printf("(%d, %d), %d\n", cur.pos.X, cur.pos.Y, cur.dist)
		if cur.pos.X == tar.X && cur.pos.Y == tar.Y {
			return cur.cost, cur.first_move
		}
		moves = possibleMoves(g, visited, h, w, cur, tar)
		for i := range moves {
			if moves[i] != nil {
				moves[i].first_move = cur.first_move
				heap.Push(&q, *moves[i])
			}
		}
	}
	for i := range moves {
		if moves[i] != nil {
			fmt.Printf("Did not find good move.\n")
			return moves[i].dist, moves[i].first_move
		}
	}
	fmt.Printf("Did not find any move.\n")
	return 127, Up
}

func runAStar(board *Board, snake, food int) (int8, MoveType) {
	g := extendGrid(board, makeGrid(board), &board.Snakes[snake])
  return runAStarWithGrid(board, snake, food, g)
}

type SnakeComp struct {
	dist  int8
	index int
	move  MoveType
}

// For each food location, return an ordered list of indices
// of all the snakes by the distance they would have to travel to
// get there.
func rankSnakesByDistanceToFood(board *Board) [][]SnakeComp {
	res := make([][]SnakeComp, len(board.Food))
	for i := range res {
		res[i] = make([]SnakeComp, len(board.Snakes))
		for j := range res[i] {
			d, m := runAStar(board, j, i)
			res[i][j] = SnakeComp{
				dist:  d,
				index: j,
				move:  m,
			}
		}
		sort.Slice(res[i], func(a, b int) bool {
			return res[i][a].dist < res[i][b].dist
		})
	}
	return res
}

func ChooseAStarMove(board *Board, you *Battlesnake) MoveType {
	id := you.ID
	var index int
	for i := range board.Snakes {
		if board.Snakes[i].ID == id {
			index = i
			break
		}
	}
	rankings := rankSnakesByDistanceToFood(board)
	relative_scores := make([]int, 0, len(rankings))
	for i := range rankings {
		fmt.Printf("rankings for food %d: %v\n", i, rankings[i])
		for j := range rankings[i] {
			if rankings[i][j].index == index {
				relative_scores = append(relative_scores, j)
			}
		}
	}
	betterThan := func(i, j int) bool {
		if relative_scores[i] < relative_scores[j] {
			return true
		}
		if relative_scores[i] == relative_scores[j] {
			return rankings[i][relative_scores[i]].dist-rankings[i][0].dist <
				rankings[j][relative_scores[j]].dist-rankings[j][0].dist
		}
		return false
	}
	betterThan2 := func(i, j int) bool {
		return rankings[i][relative_scores[i]].dist < rankings[j][relative_scores[j]].dist
	}
	_ = betterThan
	_ = betterThan2
	best := 0
	for i := 1; i < len(rankings); i++ {
		if betterThan(i, best) {
			best = i
		}
	}
	return rankings[best][relative_scores[best]].move
}

func ChooseWeightedMove(board *Board, you *Battlesnake) MoveType {
	weights := make(map[MoveType]int)
	weights[Up] = 0
	weights[Down] = 0
	weights[Right] = 0
	weights[Left] = 0
	g := extendGrid(board, makeGrid(board), you)
	if you.Head.X == 0 {
		weights[Left] = -64
	} else if you.Head.X == board.Width-1 {
		weights[Right] = -64
	}
	if you.Head.Y == 0 {
		weights[Down] = -64
	} else if you.Head.Y == board.Height-1 {
		weights[Up] = -64
	}
  _ = g
  return Up
}
