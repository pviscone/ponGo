package main

import (
	"time"
)

// ////////////////////////////////////////////
// Further additions:
// Create menu to set game parameters at runtime
// Implement online multiplayer
// Implement CPU player
//    // Level 10 = perfect prediction, < 10 random smearing of the prediction
// make the ball movement smoother in diagonal.
//    // to do it don't render the ball every frame but just when it reaches an integer value on both axis
//    // could be convenient to quantize floats
//    // checl DDA vs Bresenham Line Drawing Algorithm
//    // !!could be convevnient to switch from float aritmetic to int. Storing generic positions with 0-100 or 0-500

func main() {

	game := Game{
		plat1:      Platform{y: 0.5, char: "%", widthFrac: 0.2, score: 0},
		plat2:      Platform{y: 0.5, char: "%", widthFrac: 0.2, score: 0},
		ball:       Ball{x: 0.5, dx: -0.01, y: 0.5, dy: 0., char: "O"},
		goal_limit: 5,
	}
	const refershInterval = 16 * time.Millisecond // ~60Hz

	game.drawEmptyMap()
	cursor.updateStatus(&game)

	restore := startInputReader()
	defer restore()

	ticker := time.NewTicker(refershInterval)
	defer ticker.Stop()

	for range ticker.C {
		game.draw()
		p1_movement := getMovement(&p1_tstamp)
		p2_movement := getMovement(&p2_tstamp)
		game.update(p1_movement, p2_movement, restore)
	}
}
