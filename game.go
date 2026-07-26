package main

import (
	"fmt"
	"github.com/chewxy/math32"
	"golang.org/x/term"
	"math/rand"
	"os"
	"strings"
	"time"
)

var _, unit_frac = cursor.posToFrac(0, 1)
var _, up_boundary_frac = cursor.posToFrac(0, 2)
var _, low_boundary_frac = cursor.posToFrac(0, cursor.termsize.Height-1)

func sign(x float32) float32 {
	if x >= 0. {
		return 1.
	} else {
		return -1.
	}
}

type Ball struct {
	x, y, dx, dy float32
	char         string
}

type Platform struct {
	//widthFrac must be converted to a odd int later
	y, widthFrac float32
	score        int
	char         string
}

func (p *Platform) getIntWidth() int {
	int_width := int(p.widthFrac * float32(cursor.termsize.Height-3))
	if int_width%2 == 0 {
		int_width++
	}
	return int_width
}

func (p *Platform) move(move Movement) {
	if move == MoveUp && p.y-up_boundary_frac > p.widthFrac/2 {
		p.y -= p.widthFrac / 2
	} else if move == MoveDown && low_boundary_frac-p.y > p.widthFrac/2 {
		p.y += p.widthFrac / 2
	}
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

	//Hide cursor
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
	g.plat1.y = 0.5
	g.plat2.y = 0.5
	g.drawEmptyMap()
}

func (g *Game) bounceOnPlat(move Movement) {
	var offset float32 = 0.
	if move == MoveUp {
		offset = math32.Pi / 4
	} else if move == MoveDown {
		offset = -math32.Pi / 4
	}
	theta := math32.Atan(math32.Abs(g.ball.dy / g.ball.dx))
	new_theta := -theta + offset + (rand.Float32()-0.5)*math32.Pi*0.5
	speed := math32.Sqrt(g.ball.dx*g.ball.dx + g.ball.dy*g.ball.dy)
	g.ball.dx = -sign(g.ball.dx) * speed * math32.Cos(new_theta)
	g.ball.dy = -sign(g.ball.dy) * speed * math32.Sin(new_theta)
}

func (g *Game) update(p1_move Movement, p2_move Movement, restore func()) {

	// Update platform
	g.plat1.move(p1_move)
	g.plat2.move(p2_move)

	// Update ball
	if g.ball.x+g.ball.dx < 0 {
		if g.plat1.y-(g.plat1.widthFrac/2) < g.ball.y && g.ball.y < g.plat1.y+(g.plat1.widthFrac/2) {
			g.bounceOnPlat(p1_move)
		} else {
			g.scoreGoal(&g.plat2, restore)
		}
	} else if g.ball.x+g.ball.dx > 1 {
		if g.plat2.y-(g.plat2.widthFrac/2) < g.ball.y && g.ball.y < g.plat2.y+(g.plat2.widthFrac/2) {
			g.bounceOnPlat(p2_move)
		} else {
			g.scoreGoal(&g.plat1, restore)
		}
	}
	g.ball.x += g.ball.dx

	if g.ball.y+g.ball.dy < up_boundary_frac || g.ball.y+g.ball.dy > low_boundary_frac {
		g.ball.dy = -g.ball.dy
	}
	g.ball.y += g.ball.dy
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

	cursor.updateStatus(g)
}

func (g *Game) scoreGoal(plat *Platform, restore func()) {
	(*plat).score += 1

	g.ball.dx = float32(rand.Intn(2)*2-1) * math32.Sqrt(g.ball.dx*g.ball.dx+g.ball.dy*g.ball.dy)
	g.ball.dy = 0

	restore()
	g.resetMap()
	fd := int(os.Stdin.Fd())
	term.MakeRaw(fd)

	if (*plat).score >= g.goal_limit {
		cursor.moveTo(cursor.termsize.Width/2-5, cursor.termsize.Height/2)
		if plat == &g.plat1 {
			cursor.write("PLAYER 1 WINS")
		} else {
			cursor.write("PLAYER 2 WINS")
		}
		time.Sleep(3 * time.Second)
		restore()
		fmt.Fprint(os.Stdout, "\x1b[?25h")
		os.Exit(0)
	} else {
		time.Sleep(2 * time.Second)
	}
}
