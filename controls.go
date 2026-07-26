package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"sync/atomic"
	"time"
)

const (
	CONTROL_P1U   = 'q'
	CONTROL_P1D   = 'a'
	CONTROL_P2U   = 'o'
	CONTROL_P2D   = 'l'
	HOLD_INTERVAL = 130 * time.Millisecond
)

type Movement int

const (
	Still Movement = iota
	MoveUp
	MoveDown
)

type TStamp struct {
	lastU, lastD atomic.Int64
}

var p1_tstamp TStamp
var p2_tstamp TStamp

func touch(target *atomic.Int64) {
	target.Store(time.Now().UnixNano())
}

func getMovement(target *TStamp) Movement {
	lastU := (target.lastU).Load()
	lastD := (target.lastD).Load()
	if lastU == 0 && lastD == 0 {
		return Still
	}

	var time_since_up time.Duration = time.Duration(time.Now().UnixNano() - lastU)
	var time_since_down time.Duration = time.Duration(time.Now().UnixNano() - lastD)

	var fired_up bool = time_since_up < HOLD_INTERVAL
	var fired_down bool = time_since_down < HOLD_INTERVAL

	target.lastU.Store(0)
	target.lastD.Store(0)

	if fired_up && (time_since_up < time_since_down) {
		return MoveUp
	} else if fired_down {
		return MoveDown
	}
	return Still
}

func startInputReader() (restore func()) {
	fd := int(os.Stdin.Fd())
	oldState, _ := term.MakeRaw(fd)

	restore = func() { _ = term.Restore(fd, oldState) }

	go func() {
		r := bufio.NewReader(os.Stdin)
		buf := make([]byte, 1)
		for {
			if _, err := r.Read(buf); err != nil {
				panic("Error in reading stdin buffer")
			}
			switch buf[0] {
			case 3: // Ctrl+C
				restore()
				fmt.Fprint(os.Stdout, "\x1b[?25h")
				os.Exit(0)
			case CONTROL_P1U:
				touch(&(p1_tstamp.lastU))
			case CONTROL_P1D:
				touch(&(p1_tstamp.lastD))
			case CONTROL_P2U:
				touch(&(p2_tstamp.lastU))
			case CONTROL_P2D:
				touch(&(p2_tstamp.lastD))
			}
		}
	}()
	return restore
}
