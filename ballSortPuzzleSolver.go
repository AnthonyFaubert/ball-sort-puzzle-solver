
package main

import (
	"fmt"
	"hash/maphash"
	"encoding/binary"
)

const (
	TUBE_CAPACITY = 5
	NUM_TUBES = 15
	RECURSION_LIMIT = 2000
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
	BALL_WHITE = iota
	BALL_MAX = BALL_WHITE
)

var SET struct{}
var GameStateHashSeed maphash.Seed

type GameState [][]int8
func (a GameState) hash() uint64 {
	if (BALL_MAX > 15) || (TUBE_CAPACITY > 8) {
		panic("change the hash function")
	}
	// < 16 colors, so 4 bits per color

	var h maphash.Hash
	h.SetSeed(GameStateHashSeed)
	for i := 0; i < NUM_TUBES; i++ {
		buf := make([]byte, 4)
		th := tubeHash(a[i])
		//fmt.Printf("%05x ", th)
		binary.LittleEndian.PutUint32(buf, th)
		h.Write(buf)
	}
	s := h.Sum64()
	//fmt.Printf("%016x\n", s)
	return s
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
	state := make(GameState, NUM_TUBES)
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
	seenStates[state.hash()] = SET
	moves := availableMoves(state)
	if recursionDepth > RECURSION_LIMIT {
		panic(fmt.Sprintf("Recursed deeper than %d.", RECURSION_LIMIT))
	}
	deepest := recursionDepth
	for move := range moves {
		newState := makeMove(state, move)
		if _, newStateIsOld := seenStates[newState.hash()]; !newStateIsOld {
			if isSolved(newState) {
				fmt.Println("Found valid solution.")
				return []GameMove{move}, recursionDepth
			} else {
				fmt.Printf("%02d (%016x %%#v) Move from %d to %d\n", recursionDepth, newState.hash(), /*newState,*/ move.fromTube, move.toTube)
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
		[]int8{BALL_DRKGREEN, BALL_YELLOW, BALL_BROWN, BALL_DRKBLUE, BALL_RED},
		[]int8{BALL_DRKGREEN, BALL_TAN, BALL_ORANGE, BALL_RED, BALL_GREEN},
		[]int8{BALL_TEAL, BALL_PURPLE, BALL_DRKBLUE, BALL_GREEN, BALL_WHITE},
		[]int8{BALL_DRKBLUE, BALL_YELLOW, BALL_RED, BALL_YELLOW, BALL_BLUE},
		[]int8{BALL_DRKGREEN, BALL_BLUE, BALL_PURPLE, BALL_RED, BALL_ORANGE},
		[]int8{BALL_ORANGE, BALL_PINK, BALL_TEAL, BALL_TAN, BALL_DRKGREEN},
		[]int8{BALL_BLUE, BALL_YELLOW, BALL_BLUE, BALL_PINK, BALL_WHITE},
		[]int8{BALL_ORANGE, BALL_TEAL, BALL_RED, BALL_PINK, BALL_GREEN},
		[]int8{BALL_BROWN, BALL_TAN, BALL_WHITE, BALL_PINK, BALL_BROWN},
		[]int8{BALL_TAN, BALL_DRKBLUE, BALL_GREEN, BALL_GREEN, BALL_DRKGREEN},
		[]int8{BALL_ORANGE, BALL_TAN, BALL_YELLOW, BALL_PURPLE, BALL_TEAL},
		[]int8{BALL_PURPLE, BALL_PINK, BALL_BROWN, BALL_BROWN, BALL_WHITE},
		[]int8{BALL_TEAL, BALL_BLUE, BALL_WHITE, BALL_DRKBLUE, BALL_PURPLE},
		make([]int8, 0, TUBE_CAPACITY),
		make([]int8, 0, TUBE_CAPACITY),
		//make([]int8, 0, TUBE_CAPACITY), // Joci solved it with this additional tube.
	}

	seenStates := make(map[uint64]struct{})
	fmt.Println("Solving...")
	solution, deepestRecursion := solve(level578, seenStates, 0)
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
