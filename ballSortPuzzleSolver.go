
package main

import (
	"fmt"
	"sort"
	"hash/maphash"
	"encoding/binary"
)

const (
	TUBE_CAPACITY = 5
	RECURSION_LIMIT = 2000
	DEBUG = false
)

const (
	BALL_NONE int8 = 0
	BALL_RED = iota
	BALL_GREEN = iota
	BALL_DRKGREEN = iota
	BALL_BLUE = iota
	BALL_DRKBLUE = iota
	BALL_TEAL = iota
	BALL_PURPLE = iota
	BALL_PINK = iota
	BALL_ORANGE = iota
	BALL_YELLOW = iota
	BALL_TAN = iota
	BALL_BROWN = iota
	BALL_NAVY = iota
	BALL_WHITE = iota
	BALL_MAX = BALL_WHITE
)

var SET struct{}
var GameStateHashSeed maphash.Seed

type GameState [][]int8
func (state GameState) copy() GameState {
	cop := make(GameState, len(state))
	for i := 0; i < len(state); i++ {
		cop[i] = make([]int8, len(state[i]), TUBE_CAPACITY)
		copy(cop[i], state[i])

	}
	return cop
}

func (state GameState) hash() uint64 {
	if (BALL_MAX > 15) || (TUBE_CAPACITY > 8) {
		panic("change the hash function")
	}
	// < 16 colors, so 4 bits per color
	var h maphash.Hash
	h.SetSeed(GameStateHashSeed)

	// Sort the tubes so that ordering of tubes doesn't affect our check for if we've seen this state before or not.
	sorted := state.copy()
	sort.Slice(sorted, func(left, right int) bool {
		lTube := sorted[left]
		rTube := sorted[right]
		for k := 0; k < TUBE_CAPACITY; k++ {
			if len(rTube) <= k {
				return true
			} else if len(lTube) <= k {
				return false
			} else if lTube[k] != rTube[k] {
				return lTube[k] > rTube[k]
			}
		}
		return false
	})

	for i := 0; i < len(state); i++ {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, tubeHash(sorted[i]))
		h.Write(buf)
	}
	return h.Sum64()
}

func (state GameState) Print(prefix string) {
	for i := 0; i < len(state); i++ {
		fmt.Printf("%s%2d: ", prefix, i)
		for j := 0; j < len(state[i]); j++ {
			if j != 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%2d", state[i][j])
		}
		fmt.Println()
	}
}

// Counts the occurances of each color and makes sure you have 5 of each. Returns BALL_NONE if passed, or the problem color if failed.
func (state GameState) CheckValid() int8 {
	colorCounts := make(map[int8]int)
	for _, tube := range state {
		for _, color := range tube {
			colorCounts[color]++
		}
	}
	for color, count := range colorCounts {
		if count != TUBE_CAPACITY {
			return color
		}
	}
	return BALL_NONE
}

func tubeHash(tube []int8) (h uint32) {
	for i := 0; i < TUBE_CAPACITY; i++ {
		ball := BALL_NONE
		if i < len(tube) {
			ball = tube[i]
		}
		h |= uint32(ball & 0x0f) << (4 * i)
	}
	return
}

type GameMove struct {
	fromTube, toTube int
}
type ScoredGameMove struct {
	GameMove
	score int
}


func topBall(tube []int8) int8 {
	if len(tube) > 0 {
		return tube[len(tube)-1]
	} else {
		return BALL_NONE
	}
}

func availableMoves(state GameState) map[GameMove]struct{} {
	moves := make(map[GameMove]struct{})
	for i := 0; i < len(state); i++ {
		srcBall := topBall(state[i])
		if srcBall != BALL_NONE {
			for j := 0; j < len(state); j++ {
				if (j != i) && (len(state[j]) < TUBE_CAPACITY) {
					toBall := topBall(state[j])
					if (toBall == BALL_NONE) || (srcBall == toBall) {
						moves[GameMove{i, j}] = SET
					}
				}
			}
		}
	}
	return moves
}

func prioritizeMoves(state GameState, moveSet map[GameMove]struct{}) []GameMove {
	moves := make([]ScoredGameMove, len(moveSet))
	// If the tube only contains 1 color, the score is the number of balls in the tube.
	// Cache scores in tubeScores so you don't have to recompute them.
	tubeScores := make(map[int]int)
	scoreTube := func(tube int) int {
		score, alreadyScored := tubeScores[tube]
		if !alreadyScored {
			if len(state[tube]) > 0 {
				// If all the colors are the same, the score is the number of balls
				score = len(state[tube])
				// Reset score if they're not the same
				for i := 1; i < len(state[tube]); i++ {
					if state[tube][i] != state[tube][0] {
						score = 0
						break
					}
				}
			}
			tubeScores[tube] = score
		}
		return score
	}

	i := 0
	for move := range moveSet {
		score := scoreTube(move.toTube)
		if scoreTube(move.fromTube) > score {
			score = -(TUBE_CAPACITY + 1)
		}
		moves[i] = ScoredGameMove{move, score}
		i++
	}
	sort.Slice(moves, func(i, j int) bool {
		return moves[i].score > moves[j].score
	})

	result := make([]GameMove, len(moves), len(moves))
	for i := 0; i < len(moves); i++ {
		result[i].fromTube = moves[i].fromTube
		result[i].toTube = moves[i].toTube
	}
	return result
}

func isSolved(state GameState) bool {
	for i := 0; i < len(state); i++ {
		tube := state[i]
		if len(tube) == 0 {
			continue
		} else if len(tube) < TUBE_CAPACITY {
			return false
		} else {
			for i := 1; i < TUBE_CAPACITY; i++ {
				if tube[0] != tube[1] {
					return false
				}
			}
		}
	}
	return true
}


func removeTopBall(tube []int8) []int8 {
	return tube[:len(tube)-1]
}

func makeMove(oldState GameState, move GameMove) GameState {
	state := make(GameState, len(oldState))
	for i := 0; i < cap(state); i++ {
		length := len(oldState[i])
		if i == move.toTube {
			length++
		} else if i == move.fromTube {
			length--
		}
		state[i] = make([]int8, length, TUBE_CAPACITY)
		copy(state[i], oldState[i])
		if i == move.toTube {
			state[i][length-1] = topBall(oldState[move.fromTube])
		}

	}
	return state
}

func solve(state GameState, seenStates map[uint64]struct{}, recursionDepth int) ([]GameMove, int) {
	if isSolved(state) {
		fmt.Println("ERROR: Already solved!")
	}
	seenStates[state.hash()] = SET
	moves := availableMoves(state)
	if recursionDepth > RECURSION_LIMIT {
		fmt.Printf("FATAL: Recursed deeper than %d.\nState:", RECURSION_LIMIT)
		state.Print("  ")
		fmt.Println("Available moves:")
		for move := range moves {
			fmt.Printf("  %2d to %2d\n", move.fromTube, move.toTube)
		}
		panic("RECURSION_LIMIT exceeded.")
	}
	deepest := recursionDepth
	for _, move := range prioritizeMoves(state, moves) {
		newState := makeMove(state, move)
		if _, newStateIsOld := seenStates[newState.hash()]; !newStateIsOld {
			if isSolved(newState) {
				fmt.Println("Found valid solution.")
				return []GameMove{move}, recursionDepth
			} else {
				if DEBUG {
					fmt.Printf("%02d Move from %d to %d\n", recursionDepth, move.fromTube, move.toTube)
				}
				solution, childDepth := solve(newState, seenStates, recursionDepth + 1)
				if len(solution) > 0 {
					return append([]GameMove{move}, solution...), childDepth // [move] + solution
				}
				if childDepth > deepest {
					deepest = childDepth
				}
			}
		}
	}
	return make([]GameMove, 0), deepest
}

func test() {
	a := GameState{[]int8{3, 10, 12, 5}, []int8{3, 11, 9, 1}, []int8{6, 7, 5, 2, 2}, []int8{5, 10, 1, 10}, []int8{3, 4, 7, 1, 1}, []int8{9, 8, 6, 11, 3}, []int8{4, 10, 4, 4}, []int8{9, 6, 1, 8, 2}, []int8{12, 11, 13, 8, 8}, []int8{11, 5, 2, 2, 3}, []int8{9, 11, 10, 7, 6}, []int8{7, 8, 12, 12, 12}, []int8{6, 4, 13, 5, 7}, []int8{9}, []int8{13, 13, 13}}

	b := GameState{[]int8{3, 10, 12, 5}, []int8{3, 11, 9, 1, 1}, []int8{6, 7, 5, 2, 2}, []int8{5, 10, 1, 10}, []int8{3, 4, 7, 1}, []int8{9, 8, 6, 11, 3}, []int8{4, 10, 4, 4}, []int8{9, 6, 1, 8, 2}, []int8{12, 11, 13, 8, 8}, []int8{11, 5, 2, 2, 3}, []int8{9, 11, 10, 7, 6}, []int8{7, 8, 12, 12, 12}, []int8{6, 4, 13, 5, 7}, []int8{9}, []int8{13, 13, 13}}

	c := GameState{[]int8{3, 10, 12, 5}, []int8{3, 11, 9, 1}, []int8{6, 7, 5, 2, 2}, []int8{5, 10, 1, 10}, []int8{3, 4, 7, 1, 1}, []int8{9, 8, 6, 11, 3}, []int8{4, 10, 4, 4}, []int8{9, 6, 1, 8, 2}, []int8{12, 11, 13, 8, 8}, []int8{11, 5, 2, 2, 3}, []int8{9, 11, 10, 7, 6}, []int8{7, 8, 12, 12, 12}, []int8{6, 4, 13, 5, 7}, []int8{9}, []int8{13, 13, 13}}

	d := GameState{[]int8{3, 10, 12, 5}, []int8{3, 11, 9, 1, 1}, []int8{6, 7, 5, 2, 2}, []int8{5, 10, 1, 10}, []int8{3, 4, 7, 1}, []int8{9, 8, 6, 11, 3}, []int8{4, 10, 4, 4}, []int8{9, 6, 1, 8, 2}, []int8{12, 11, 13, 8, 8}, []int8{11, 5, 2, 2, 3}, []int8{9, 11, 10, 7, 6}, []int8{7, 8, 12, 12, 12}, []int8{6, 4, 13, 5, 7}, []int8{9}, []int8{13, 13, 13}}



	solved := GameState{
		[]int8{BALL_RED, BALL_RED, BALL_RED, BALL_RED, BALL_RED},
		[]int8{BALL_GREEN, BALL_GREEN, BALL_GREEN, BALL_GREEN, BALL_GREEN},
		[]int8{BALL_DRKGREEN, BALL_DRKGREEN, BALL_DRKGREEN, BALL_DRKGREEN, BALL_DRKGREEN},
		[]int8{BALL_BLUE, BALL_BLUE, BALL_BLUE, BALL_BLUE, BALL_BLUE},
		[]int8{BALL_DRKBLUE, BALL_DRKBLUE, BALL_DRKBLUE, BALL_DRKBLUE, BALL_DRKBLUE},
		[]int8{BALL_TEAL, BALL_TEAL, BALL_TEAL, BALL_TEAL, BALL_TEAL},
		[]int8{BALL_PURPLE, BALL_PURPLE, BALL_PURPLE, BALL_PURPLE, BALL_PURPLE},
		[]int8{BALL_PINK, BALL_PINK, BALL_PINK, BALL_PINK, BALL_PINK},
		[]int8{BALL_ORANGE, BALL_ORANGE, BALL_ORANGE, BALL_ORANGE, BALL_ORANGE},
		[]int8{BALL_YELLOW, BALL_YELLOW, BALL_YELLOW, BALL_YELLOW, BALL_YELLOW},
		[]int8{BALL_TAN, BALL_TAN, BALL_TAN, BALL_TAN, BALL_TAN},
		[]int8{BALL_BROWN, BALL_BROWN, BALL_BROWN, BALL_BROWN, BALL_BROWN},
		[]int8{BALL_WHITE, BALL_WHITE, BALL_WHITE, BALL_WHITE, BALL_WHITE},
		[]int8{},
		[]int8{},
	}
	if solved.CheckValid() != BALL_NONE {
		panic("CheckValid() is bad")
	}
	if !isSolved(solved) {
		panic("isSolved() is bad")
	}

	//fmt.Print("Hashing A: ")
	ah := a.hash()
	//fmt.Print("Hashing B: ")
	bh := b.hash()
	//fmt.Print("Hashing C: ")
	ch := c.hash()
	//fmt.Print("Hashing D: ")
	dh := d.hash()
	if !((ah == ch) && (bh == dh) && (ah != bh)) {
		panic(fmt.Sprintf("Tests failed! %d %d %d %d\n", ah, bh, ch, dh))
	}
}

func main() {
	GameStateHashSeed = maphash.MakeSeed()
	test()

	level578 := GameState{
		[]int8{BALL_DRKGREEN, BALL_YELLOW, BALL_BROWN, BALL_DRKBLUE, BALL_RED}, // tube 0
		[]int8{BALL_DRKGREEN, BALL_TAN, BALL_ORANGE, BALL_RED, BALL_GREEN},     // tube 1
		[]int8{BALL_TEAL, BALL_PURPLE, BALL_DRKBLUE, BALL_GREEN, BALL_WHITE},   // tube 2
		[]int8{BALL_DRKBLUE, BALL_YELLOW, BALL_RED, BALL_YELLOW, BALL_BLUE},    // ...
		[]int8{BALL_DRKGREEN, BALL_BLUE, BALL_PURPLE, BALL_RED, BALL_ORANGE},
		[]int8{BALL_ORANGE, BALL_PINK, BALL_TEAL, BALL_TAN, BALL_DRKGREEN},
		[]int8{BALL_BLUE, BALL_YELLOW, BALL_BLUE, BALL_PINK, BALL_WHITE},
		[]int8{BALL_ORANGE, BALL_TEAL, BALL_RED, BALL_PINK, BALL_GREEN},        // tube 7

		[]int8{BALL_BROWN, BALL_TAN, BALL_WHITE, BALL_PINK, BALL_BROWN},        // tube 8
		[]int8{BALL_TAN, BALL_DRKBLUE, BALL_GREEN, BALL_GREEN, BALL_DRKGREEN},  // ...
		[]int8{BALL_ORANGE, BALL_TAN, BALL_YELLOW, BALL_PURPLE, BALL_TEAL},
		[]int8{BALL_PURPLE, BALL_PINK, BALL_BROWN, BALL_BROWN, BALL_WHITE},
		[]int8{BALL_TEAL, BALL_BLUE, BALL_WHITE, BALL_DRKBLUE, BALL_PURPLE},    // tube 12
		make([]int8, 0, TUBE_CAPACITY),					        // tube 13
		make([]int8, 0, TUBE_CAPACITY),					        // tube 14
		make([]int8, 0, TUBE_CAPACITY),                                         // tube 15 (the extra tube from tapping the "+1 tube" button)
		// This level has no solution without the additional tube (tube 15).
	}
	if level578.CheckValid() != BALL_NONE { panic("level578 invalid") }

	level656 := GameState{
		[]int8{BALL_RED, BALL_BROWN, BALL_RED, BALL_DRKBLUE, BALL_DRKBLUE},
		[]int8{BALL_TAN, BALL_TEAL, BALL_PINK, BALL_RED, BALL_TAN},
		[]int8{BALL_YELLOW, BALL_GREEN, BALL_GREEN, BALL_YELLOW, BALL_PINK},
		[]int8{BALL_DRKBLUE, BALL_WHITE, BALL_PURPLE, BALL_NAVY, BALL_TEAL},
		[]int8{BALL_WHITE, BALL_BLUE, BALL_TAN, BALL_BLUE, BALL_PINK},
		[]int8{BALL_NAVY, BALL_ORANGE, BALL_BROWN, BALL_TAN, BALL_BLUE},
		[]int8{BALL_BROWN, BALL_WHITE, BALL_DRKGREEN, BALL_BROWN, BALL_PINK},
		[]int8{BALL_YELLOW, BALL_GREEN, BALL_NAVY, BALL_ORANGE, BALL_BLUE},

		[]int8{BALL_YELLOW, BALL_WHITE, BALL_NAVY, BALL_BLUE, BALL_ORANGE},
		[]int8{BALL_DRKBLUE, BALL_ORANGE, BALL_DRKGREEN, BALL_BROWN, BALL_TEAL},
		[]int8{BALL_DRKGREEN, BALL_TEAL, BALL_WHITE, BALL_RED, BALL_DRKGREEN},
		[]int8{BALL_TEAL, BALL_ORANGE, BALL_PINK, BALL_DRKBLUE, BALL_GREEN},
		[]int8{BALL_PURPLE, BALL_GREEN, BALL_PURPLE, BALL_TAN, BALL_RED},
		[]int8{BALL_NAVY, BALL_DRKGREEN, BALL_PURPLE, BALL_YELLOW, BALL_PURPLE},
		make([]int8, 0, TUBE_CAPACITY),
		make([]int8, 0, TUBE_CAPACITY),
		// No additional tube needed to solve.
	}
	if level656.CheckValid() != BALL_NONE { panic("level656 invalid") }

	level732 := GameState{
		[]int8{BALL_PINK,    BALL_RED,    BALL_BROWN,  BALL_BLUE,    BALL_GREEN},
		[]int8{BALL_PURPLE,  BALL_ORANGE, BALL_WHITE,  BALL_ORANGE,  BALL_GREEN},
		[]int8{BALL_RED,     BALL_BROWN,  BALL_PURPLE, BALL_DRKGREEN,BALL_BLUE},
		[]int8{BALL_DRKGREEN,BALL_PINK,   BALL_YELLOW, BALL_BROWN,   BALL_ORANGE},
		[]int8{BALL_NAVY,    BALL_YELLOW, BALL_RED,    BALL_BROWN,   BALL_PURPLE},
		[]int8{BALL_BLUE,    BALL_DRKBLUE,BALL_TAN,    BALL_DRKBLUE, BALL_TEAL},
		[]int8{BALL_BLUE,    BALL_GREEN,  BALL_NAVY,   BALL_TAN,     BALL_WHITE},
		[]int8{BALL_WHITE,   BALL_NAVY,   BALL_TEAL,   BALL_DRKGREEN,BALL_DRKGREEN},

		[]int8{BALL_DRKGREEN,BALL_NAVY,   BALL_DRKBLUE,BALL_DRKBLUE, BALL_RED},
		[]int8{BALL_NAVY,    BALL_TEAL,   BALL_GREEN,  BALL_PINK,    BALL_PURPLE},
		[]int8{BALL_PURPLE,  BALL_TEAL,   BALL_ORANGE, BALL_TAN,     BALL_ORANGE},
		[]int8{BALL_GREEN,   BALL_YELLOW, BALL_BLUE,   BALL_RED,     BALL_PINK},
		[]int8{BALL_YELLOW,  BALL_YELLOW, BALL_BROWN,  BALL_TEAL,    BALL_WHITE},
		[]int8{BALL_TAN,     BALL_DRKBLUE,BALL_WHITE,  BALL_TAN,     BALL_PINK},
		make([]int8, 0, TUBE_CAPACITY),
		make([]int8, 0, TUBE_CAPACITY),
	}
	if level732.CheckValid() != BALL_NONE { panic("level732 invalid") }

	seenStates := make(map[uint64]struct{})
	fmt.Println("Solving...")
	_ = level578
	solution, deepestRecursion := solve(level656, seenStates, 0)
	if len(solution) == 0 {
		fmt.Println("No solution!")
	} else {
		fmt.Printf("Solved! Number of moves: %d.\n", len(solution))
	}
	fmt.Printf("Deepest recursion: %d. Number of unique game states: %d.\n", deepestRecursion, len(seenStates))
	for _, move := range solution {
		fmt.Printf("Move from %d to %d.\n", move.fromTube, move.toTube)
	}
}
