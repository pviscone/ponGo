package main

import (
	"fmt"
	"github.com/kopoli/go-terminal-size"
	"math/rand"
	"os"
	"strings"
	"time"
)

func moveCursor(x int, y int) {
	fmt.Fprintf(os.Stdout, "\x1b[%d;%dH", y, x)
}

func moveCursorTo(x int, y int) {
	fmt.Fprint(os.Stdout, "\x1b[H")
	moveCursor(x, y)
}

type Ball struct {
	x, y, dx, dy float32
	char         string
}

type Platform struct {
	//widthFrac must be converted to a odd int later
	y, dy, widthFrac float32
	score            int
	char             string
}

func (p Platform) getIntWidth(term_height int) {
	//TODO return the closest odd integer to widthFrac*term_height
}

type Game struct {
	size       tsize.Size
	ball       Ball
	plat1      Platform
	plat2      Platform
	goal_limit int
}

func (g Game) drawEmptyMap() {
	// Clear + cursor home
	fmt.Fprint(os.Stdout, "\x1b[H\x1b[2J")
	// Hide cursor
	fmt.Fprint(os.Stdout, "\x1b[?25l")

	score1Str := fmt.Sprintf("%d", g.plat1.score)
	score2Str := fmt.Sprintf("%d", g.plat2.score)
	middleSpaces := g.size.Width - 4 - len(score1Str) - len(score2Str)

	scoreString := strings.Repeat(" ", 2) + score1Str + strings.Repeat(" ", middleSpaces) + score2Str + strings.Repeat(" ", 2)
	fmt.Fprint(os.Stdout, scoreString)
	fmt.Fprint(os.Stdout, "\n")

	fmt.Fprint(os.Stdout, strings.Repeat("-", g.size.Width))
	fmt.Fprint(os.Stdout, "\n")

	for h := 0; h < g.size.Height-4; h++ {
		fmt.Fprint(os.Stdout, strings.Repeat(" ", g.size.Width))
		fmt.Fprint(os.Stdout, "\n")
	}

	fmt.Fprint(os.Stdout, strings.Repeat("-", g.size.Width))
	fmt.Fprint(os.Stdout, "\n")
}

func (g Game) resetMap() {
	g.ball.x = float32(g.size.Width / 2)
	g.ball.y = float32(1 + (g.size.Height-2)/2)
	g.drawEmptyMap()
	g.draw()
}

func (g Game) update() {
	//Move ball logic here
	//Move ball by dx dy
	//if (dy>0 && y==(size.Height-2) || dy<0 && y==1 ) then dy=-dy and move
	//if (dx>)
}

func (g Game) draw() {
	//Draw ball and platform
}

func (g Game) scoreGoal(plat *Platform) {
	(*plat).score += 1
	g.ball.dy = 0
	g.ball.dx = float32(rand.Intn(2)*2 - 1)
	g.resetMap()
	time.Sleep(2 * time.Second)
}

//////////////////////////////////////////////

func main() {

	size, _ := tsize.GetSize()
	game := Game{
		plat1:      Platform{y: 0, dy: 0, char: "#", widthFrac: 0.1, score: 0},
		plat2:      Platform{y: 0, dy: 0, char: "#", widthFrac: 0.1, score: 0},
		ball:       Ball{x: float32(size.Width / 2), dx: -1, y: float32(1 + (size.Height-2)/2), dy: 0, char: "O"},
		size:       size,
		goal_limit: 10,
	}

	game.resetMap()
	defer func() {
		// Show cursor again on exit
		fmt.Fprint(os.Stdout, "\x1b[?25h")
	}()

	const time_per_frame = 100 * time.Millisecond
	for {
		start_frame_time := time.Now()

		size, _ = tsize.GetSize()

		elapsed_frame_time := time.Since(start_frame_time)
		if elapsed_frame_time < time_per_frame {
			time.Sleep(time_per_frame - elapsed_frame_time)
		}
	}
}
