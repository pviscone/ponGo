package main

import (
	"fmt"
	"github.com/kopoli/go-terminal-size"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Cursor struct {
	x, y, ballX, ballY, plat1Y, plat2Y int
	termsize                           tsize.Size
}

func (c *Cursor) fracToPos(fracX float32, fracY float32) (int, int) {
	return int(fracX * float32(c.termsize.Width)), int(fracY * float32(c.termsize.Height-1))
}

func (c *Cursor) clean() {
	fmt.Fprint(os.Stdout, "\x1b[2J")
}

func (c *Cursor) goHome() {
	c.x = 0
	c.y = 0
	fmt.Fprint(os.Stdout, "\x1b[H")
}

func (c *Cursor) writeAndMove(str string) {
	// \x1b[s saves the cursor position,
	// writes the string,
	// and \x1b[u restores the cursor back to where it started.
	fmt.Fprint(os.Stdout, str)
}

func (c *Cursor) write(str string) {
	// \x1b[s saves the cursor position,
	// writes the string,
	// and \x1b[u restores the cursor back to where it started.
	fmt.Fprintf(os.Stdout, "\x1b[s%s\x1b[u", str)
}

func (c *Cursor) move(x int, y int) {
	c.x = c.x + x
	c.y = c.y + y

	var xCode, yCode string

	// Horizontal movement
	if x > 0 {
		xCode = fmt.Sprintf("\x1b[%dC", x) //Right
	} else if x < 0 {
		xCode = fmt.Sprintf("\x1b[%dD", -x) //Left
	}

	// Vertical movement
	if y > 0 {
		yCode = fmt.Sprintf("\x1b[%dB", y) //Down
	} else if y < 0 {
		yCode = fmt.Sprintf("\x1b[%dA", -y) //Up
	}
	fmt.Fprint(os.Stdout, xCode+yCode)
}

func (c *Cursor) moveTo(x int, y int) {
	c.goHome()
	c.move(x, y)
}

func (c *Cursor) moveFrac(fracX float32, fracY float32) {
	c.move(c.fracToPos(fracX, fracY))
}

func (c *Cursor) moveFracTo(fracX float32, fracY float32) {
	c.moveTo(c.fracToPos(fracX, fracY))
}

var size, _ = tsize.GetSize()
var cursor = Cursor{
	x: 0, y: 0,
	ballX: 0, ballY: 0,
	plat1Y: 0, plat2Y: 0,
	termsize: size,
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

func (p *Platform) getIntWidth() int {
	int_width := int(p.widthFrac * float32(cursor.termsize.Height-3))
	if int_width%2 == 0 {
		int_width++
	}
	return int_width
}

type Game struct {
	ball       Ball
	plat1      Platform
	plat2      Platform
	goal_limit int
}

func (g *Game) drawEmptyMap() {
	// Clean + cursor home
	cursor.clean()
	cursor.goHome()

	// Hide cursor
	fmt.Fprint(os.Stdout, "\x1b[?25l")

	score1Str := fmt.Sprintf("%d", g.plat1.score)
	score2Str := fmt.Sprintf("%d", g.plat2.score)
	middleSpaces := cursor.termsize.Width - 4 - len(score1Str) - len(score2Str)

	scoreString := strings.Repeat(" ", 2) + score1Str + strings.Repeat(" ", middleSpaces) + score2Str + strings.Repeat(" ", 2)
	cursor.writeAndMove(scoreString)
	cursor.writeAndMove("\n")

	cursor.writeAndMove(strings.Repeat("-", cursor.termsize.Width))
	cursor.writeAndMove("\n")

	for h := 0; h < cursor.termsize.Height-4; h++ {
		cursor.writeAndMove(strings.Repeat(" ", cursor.termsize.Width))
		cursor.writeAndMove("\n")
	}

	cursor.writeAndMove(strings.Repeat("-", cursor.termsize.Width))
	cursor.writeAndMove("\n")

	cursor.goHome()
}

func (g *Game) resetMap() {
	g.ball.x = 0.5
	g.ball.y = 0.5
	g.drawEmptyMap()
	//g.draw()
}

func (g *Game) update() {
	//Move ball logic here
	//Move ball by dx dy
	//if (dy>0 && y==(size.Height-2) || dy<0 && y==1 ) then dy=-dy and move
	//if (dx>)
	//Call draw here if the int of the ball or of the plat change
	//call score here (recheck implementation)
	cursor.ballX = int(g.ball.x)
	cursor.ballY = int(g.ball.y) - 1
	cursor.plat1Y = int(g.plat1.y) - 1
	cursor.plat2Y = int(g.plat2.y) - 1
}

func (g *Game) draw() {
	//Move ball
	cursor.moveTo(cursor.ballX, cursor.ballY)
	cursor.write(" ")
	cursor.moveFracTo(g.ball.x, g.ball.y)
	cursor.write(g.ball.char)

	//Move plat1
	var plat1_width int = g.plat1.getIntWidth()
	cursor.moveTo(0, cursor.plat1Y)
	cursor.move(0, -((plat1_width - 1) / 2))
	for i := 0; i < plat1_width; i++ {
		cursor.write(" ")
		cursor.move(0, 1)
	}

	cursor.moveFracTo(0., g.plat1.y)
	cursor.move(0, -((plat1_width - 1) / 2))
	for i := 0; i < plat1_width; i++ {
		cursor.write(g.plat1.char)
		cursor.move(0, 1)
	}

	//Move plat2
	var plat2_width int = g.plat2.getIntWidth()
	cursor.moveTo(cursor.termsize.Width, cursor.plat2Y)
	cursor.move(0, -((plat2_width - 1) / 2))
	for i := 0; i < plat2_width; i++ {
		cursor.write(" ")
		cursor.move(0, 1)
	}

	cursor.moveFracTo(1., g.plat2.y)
	cursor.move(0, -((plat2_width - 1) / 2))
	for i := 0; i < plat2_width; i++ {
		cursor.write(g.plat2.char)
		cursor.move(0, 1)
	}
}

func (g *Game) scoreGoal(plat *Platform) {
	(*plat).score += 1
	g.ball.dy = 0
	g.ball.dx = float32(rand.Intn(2)*2-1) * 0.05
	g.resetMap()
	time.Sleep(2 * time.Second)
}

//////////////////////////////////////////////

func main() {

	game := Game{
		plat1:      Platform{y: 0.5, dy: 0, char: "%", widthFrac: 0.15, score: 0},
		plat2:      Platform{y: 0.5, dy: 0, char: "%", widthFrac: 0.15, score: 0},
		ball:       Ball{x: 0.5, dx: -0.05, y: 0.5, dy: 0, char: "O"},
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
		game.update()
		game.draw()

		elapsed_frame_time := time.Since(start_frame_time)
		if elapsed_frame_time < time_per_frame {
			time.Sleep(time_per_frame - elapsed_frame_time)
		}
	}
}
