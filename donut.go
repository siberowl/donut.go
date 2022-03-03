package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"
)

const (
	screen_width  = 60
	screen_height = 60
	theta_spacing = 0.07
	phi_spacing   = 0.02
	rot_spacing   = 0.03
	R1            = 1
	R2            = 2
	K2            = 10
	K1            = float64(screen_width) * K2 * 3 / (8 * (R1 + R2))
	offset        = 30
	delay         = 40
)

func donut(K2 float64, R1 float64, R2 float64, t float64, p float64, A float64, B float64) ([3]float64, float64) {
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
			cConst*(cp*sB-cB*sA*sp) + R1*cA*cB*st,
			cA*(R2+R1*ct)*sp + R1*sA*st,
		},
		cp*ct*sB - cA*ct*sp - sA*st + cB*(cA*st-ct*sA*sp)
}

func projection(K1 float64, K2 float64, X [3]float64) [2]float64 {
	return [2]float64{
		K1 * X[0] / (K2 + X[2]),
		K1 * X[1] / (K2 + X[2]),
	}
}

func main() {
	A := 0.0
	B := 0.0

	for {
		var render_string string
		var luminescence [screen_width][screen_height]int
		var zbuffer [screen_width][screen_height]float64
		for theta := 0.0; theta < 2*math.Pi; theta += theta_spacing {
			for phi := 0.0; phi < 2*math.Pi; phi += phi_spacing {
				var d, L = donut(K2, R1, R2, theta, phi, A, B)
				if L > 0 {
					var point = projection(K1, K2, d)
					var x int = int(point[0] + offset)
					var y int = int(point[1] + offset)
					if x > 0 && x < screen_width {
						if y > 0 && y < screen_height {
							if 1/(d[2]+K2) > zbuffer[x][y] {
								zbuffer[x][y] = 1 / (d[2] + K2)
								luminescence[x][y] = int(8 * L)
							}
						}
					}
				}
			}
		}

		render_string = ""
		for j := 0; j < screen_height; j++ {
			for i := 0; i < screen_width; i++ {
				if i < screen_width-1 {
					if luminescence[i][j] > 0 {
						render_string += string(".,-~:;=!*#$@"[luminescence[i][j]]) + string(".,-~:;=!*#$@"[luminescence[i][j]])
					} else {
						render_string += "  "
					}
				} else {
					render_string += "\n"
				}
			}
		}
		fmt.Print("\033[H\033[2J") // clear previous stdout
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
		fmt.Printf(render_string)

		A += rot_spacing
		B += rot_spacing
		time.Sleep(delay * time.Millisecond)
	}
}
