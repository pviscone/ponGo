package main

import (
	"fmt"
	"os"
	"time"
)

// ////////////////////////////////////////////
// TODO:
// implement scattering on platform logic
//    // Reflection + gaussian noise (dependent on dy)
// implement controls
// 	  // Use atomic ints to store last pressure of qapl keys
//    // Check in the render loop that Delta time with pressure < frametime/3 (e.g.)
// capture SIGWINCH to recompute termsize and redraw
// capture ctrl-c to unhide cursor (or ESC key)

// Refactor:
// change render loop using ticks

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
		plat1:      Platform{y: 0.5, dy: 0, char: "%", widthFrac: 0.15, score: 0},
		plat2:      Platform{y: 0.5, dy: 0, char: "%", widthFrac: 0.15, score: 0},
		ball:       Ball{x: 0.5, dx: -0.01, y: 0.5, dy: 0., char: "O"},
		goal_limit: 5,
	}
	game.drawEmptyMap()
	cursor.updateStatus(&game)
	defer func() {
		// Show cursor again on exit
		fmt.Fprint(os.Stdout, "\x1b[?25h")
	}()

	const time_per_frame = 10 * time.Millisecond

	for {
		start_frame_time := time.Now()
		game.draw()
		game.update()

		elapsed_frame_time := time.Since(start_frame_time)
		if elapsed_frame_time < time_per_frame {
			time.Sleep(time_per_frame - elapsed_frame_time)
		}
	}
}
