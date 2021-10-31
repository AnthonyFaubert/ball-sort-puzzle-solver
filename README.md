# ball-sort-puzzle-solver
Golang solver algorithm for the "Ball Sort Puzzle" game.

I made this because level 578 appeared impossible to solve without using the "+1 tube" button. Turns out level 578 is indeed impossible to solve without the extra tube. I also tested level 656, but that appears to be solvable in as little as 108 moves without needing an extra tube.

# Understanding the GameState representation:
Open up the .go file and look at the following lines from that file: (I added a few comments)

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
		make([]int8, 0, TUBE_CAPACITY),                                         // tube 13
		make([]int8, 0, TUBE_CAPACITY),                                         // tube 14
		make([]int8, 0, TUBE_CAPACITY),                                         // tube 15 (the extra tube from tapping the "+1 tube" button)
    }

That's the GameState representation of this initial game board:

![Reference screenshot for level 578 of the ball sort puzzle game](https://github.com/AnthonyFaubert/ball-sort-puzzle-solver/blob/master/Level578Example.jpg?raw=true)

The bottom ball in a tube is the 0th ball in the corresponding tube list. The tubes are numbered as shown in the image and comments.

# Adding a new level and solving it:
1. Copy any of the existing GameState variable declarations (`level578` or `level656` or `level732` or `level764`) and change all the balls/tubes to match what you've got in the game right now.
Any partially-full tube lists need to be padded out with BALL_NONE values where there's an empty space.
2. Change the `solution, deepestRecursion := solve(level656, seenStates, 0)` line to use your new GameState variable.
3. Run it with `go run ballSortPuzzleSolver.go`

You can add 1 new color in-between the "BALL_NONE" and "BALL_WHITE" definitions. Only 1 because my `tubeHash()` function can only handle up to 16 colors including BALL_NONE (BALL_MAX isn't a separate color).

The order of the tubes within the GameState doesn't matter, so if you need to add more tubes, add them wherever makes sense.

# GUI?
If somebody wants to write some sort of front end for this, that'd be super cool. I'd love to host a little website that uses this.
