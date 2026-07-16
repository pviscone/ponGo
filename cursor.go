package main

import (
	"fmt"
	"os"

	tsize "github.com/kopoli/go-terminal-size"
)

type Cursor struct {
	x, y, ballX, ballY, plat1Y, plat2Y int
	termsize                           tsize.Size
}

func (c *Cursor) fracToPos(fracX float32, fracY float32) (int, int) {
	return int(fracX * float32(c.termsize.Width)), int(fracY * float32(c.termsize.Height-1))
}

func (c *Cursor) posToFrac(x int, y int) (float32, float32) {
	return float32(x) / float32(c.termsize.Width), float32(y) / float32(c.termsize.Height)
}

func (c *Cursor) updateStatus(game *Game) {
	c.ballX, c.ballY = c.fracToPos(game.ball.x, game.ball.y)
	_, c.plat1Y = c.fracToPos(0., game.plat1.y)
	_, c.plat2Y = c.fracToPos(1., game.plat2.y)
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
