tep size is perfectly lockpackage main

import (
	"fmt"
	"os"
	"time"
)

// ////////////////////////////////////////////
// TODO:
// implement controls
// capture SIGWINCH to recompute termsize and redraw
// capture ctrl-c to unhide cursor (or ESC key)
// make the ball movement smoother in diagonal.
//   to do it don't render the ball every frame but just when it reaches an integer value on both axis
//   could be convenient to quantize floats

// TODO bugs:
// - Setting a ball.y !=0.5 but the ball is still drawn at 0.5

// Further additions
// Create menu to set game parameters at runtime
// Implement online multiplayer

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
