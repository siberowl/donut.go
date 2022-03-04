package main

import (
	"fmt"
	"math"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
)

const (
	screen_width  = 60
	screen_height = 60
	theta_spacing = 0.07
	phi_spacing   = 0.02
	delay         = 16
	A_spacing     = delay * 0.07 / 40
	B_spacing     = delay * 0.03 / 40
	R1            = 1
	R2            = 2
	K2            = 20
	K1            = 100
	x_offset      = 30
	y_offset      = 20
)

func donut(t float64, p float64, A float64, B float64) ([3]float64, float64) {
	ct := math.Cos(t)
	st := math.Sin(t)
	cp := math.Cos(p)
	sp := math.Sin(p)
	cA := math.Cos(A)
	sA := math.Sin(A)
	cB := math.Cos(B)
	sB := math.Sin(B)
	cConst := R2 + R1*ct
	return [3]float64{
			cConst*(cB*cp+sA*sB*sp) - R1*cA*sB*st,
			-cConst*(cp*sB-cB*sA*sp) - R1*cA*cB*st,
			cA*(R2+R1*ct)*sp + R1*sA*st,
		},
		cp*ct*sB - cA*ct*sp - sA*st + cB*(cA*st-ct*sA*sp)
}

func projection(X [3]float64) [2]float64 {
	return [2]float64{
		K1 * X[0] / (K2 + X[2]),
		K1 * X[1] / (K2 + X[2]),
	}
}

func main() {
	A := 0.0
	B := 0.0

	running := true

	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for running {
		var luminescence [screen_width][screen_height]int

		for j := 0; j < screen_height; j++ {
			for i := 0; i < screen_width; i++ {
				termbox.SetCell(i*2-1, j, rune(" "[0]), termbox.ColorWhite, termbox.ColorDefault)
				termbox.SetCell(i*2, j, rune(" "[0]), termbox.ColorWhite, termbox.ColorDefault)
			}
		}

		var zbuffer [screen_width][screen_height]float64
		for theta := 0.0; theta < 2*math.Pi; theta += theta_spacing {
			for phi := 0.0; phi < 2*math.Pi; phi += phi_spacing {
				var d, L = donut(theta, phi, A, B)
				if L > -1 {
					var point = projection(d)
					var x int = int(point[0] + x_offset)
					var y int = int(point[1] + y_offset)
					if 1/(d[2]+K2) > zbuffer[x][y] {
						zbuffer[x][y] = 1 / (d[2] + K2)
						luminescence[x][y] = int(8 * L)
						if luminescence[x][y] >= 0 {
							termbox.SetCell(x*2-1, y, rune(".,-~:;=!*#$@"[luminescence[x][y]]), termbox.ColorWhite, termbox.ColorDefault)
							termbox.SetCell(x*2, y, rune(".,-~:;=!*#$@"[luminescence[x][y]]), termbox.ColorWhite, termbox.ColorDefault)
						}
					}
				}
			}
		}

		termbox.Flush()
		A += A_spacing
		B += B_spacing
		time.Sleep(delay * time.Millisecond)
	}

	termbox.Close()

}
