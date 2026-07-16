package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

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
	g.drawEmptyMap()
	//g.draw()
}

func (g *Game) update() {
	if g.ball.x+g.ball.dx < 0 {
		if g.plat1.y-(g.plat1.widthFrac/2) < g.ball.y && g.ball.y < g.plat1.y+(g.plat1.widthFrac/2) {
			g.ball.dx = -g.ball.dx
		} else {
			g.scoreGoal(&g.plat1)
		}
	} else if g.ball.x+g.ball.dx > 1 {
		if g.plat1.y-(g.plat1.widthFrac/2) < g.ball.y && g.ball.y < g.plat1.y+(g.plat1.widthFrac/2) {
			g.ball.dx = -g.ball.dx
		} else {
			g.scoreGoal(&g.plat2)
		}
	}
	g.ball.x += g.ball.dx

	_, up_boundary_frac := cursor.posToFrac(0, 2)
	_, low_boundary_frac := cursor.posToFrac(0, cursor.termsize.Height-1)

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

func (g *Game) scoreGoal(plat *Platform) {
	(*plat).score += 1
	g.ball.dy = 0
	g.ball.dx = float32(rand.Intn(2)*2-1) * 0.01
	g.resetMap()

	if (*plat).score >= g.goal_limit {
		cursor.moveTo(cursor.termsize.Width/2-5, cursor.termsize.Height/2)
		if plat == &g.plat1 {
			cursor.write("PLAYER 1 WINS")
		} else {
			cursor.write("PLAYER 2 WINS")
		}
		time.Sleep(3 * time.Second)
		fmt.Fprint(os.Stdout, "\x1b[?25h")
		os.Exit(0)
	} else {
		time.Sleep(2 * time.Second)
	}
}
