package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

//--Enum
type gameState int

const (
	start gameState = iota
	play
)

var state = start

//--
// Number Score
var nums = [][]byte{
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
	{
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1,
	},
	{
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
}

func drawNum(pos pos, color color, size int, number int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range nums[number] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}
func clear(pixles []byte) {
	for i := range pixles {
		pixles[i] = 0
	}
}

func lerp(a, b, pct float32) float32 {
	return a + pct*(b-a)
}

type color struct {
	r, g, b byte
}
type pos struct {
	x, y float32
}
type ball struct {
	pos
	radius int
	xv     float32
	yv     float32
	color  color
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.color, pixels)
			}
		}
	}
}
func getCenter() pos {
	return pos{float32(winWidth) / 2, float32(winHeight) / 2}
}
func (ball *ball) update(leftPad *paddle, rightPad *paddle, elapsedTime float32) {
	ball.x += ball.xv * elapsedTime
	ball.y += ball.yv * elapsedTime
	// Handle collisions
	if ball.y-float32(ball.radius) < 0 || ball.y+float32(ball.radius) > float32(winHeight) {
		ball.yv = -ball.yv
	}
	if ball.x < 0 {
		rightPad.score++
		ball.pos = getCenter()
		state = start
	} else if ball.x > float32(winWidth) {
		leftPad.score++
		ball.pos = getCenter()
		state = start
	}
	// Ball hits the edhe of left paddle
	if int(ball.x)-ball.radius < int(leftPad.x)+int(rightPad.w/2) {
		if ball.y > leftPad.y-float32(leftPad.h/2) && ball.y < leftPad.y+float32(leftPad.h/2) {
			ball.xv = -ball.xv
		}
	}
	if int(ball.x)+ball.radius > int(rightPad.x)-int(rightPad.w/2) {
		if ball.y > rightPad.y-float32(rightPad.h/2) && ball.y < rightPad.y+float32(rightPad.h/2) {
			ball.xv = -ball.xv
		}
	}
}

type paddle struct {
	pos
	w     int
	h     int
	speed float32
	color color
	score int
}

func (paddle *paddle) draw(pixels []byte) {
	startx := int(paddle.x) - paddle.w/2
	starty := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startx+x, starty+y, paddle.color, pixels)
		}
	}

	numX := lerp(paddle.x, getCenter().x, .2)
	drawNum(pos{numX, 35}, paddle.color, 10, paddle.score, pixels)

}
func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= paddle.speed * elapsedTime
	} else if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
	}

}
func (paddle *paddle) updateP2(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_Q] != 0 {
		paddle.y -= paddle.speed * elapsedTime
	} else if keyState[sdl.SCANCODE_A] != 0 {
		paddle.y += paddle.speed * elapsedTime
	}
}
func setPixel(x, y int, c color, pixles []byte) {
	index := (y*winWidth + x) * 4
	if index < (len(pixles)-4) && index >= 0 {
		pixles[index] = c.r
		pixles[index+1] = c.g
		pixles[index+2] = c.b

	}
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()
	window, err := sdl.CreateWindow("Title", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING, 800, 600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos{100, 100}, 20, 100, 300, color{255, 255, 255}, 0}
	player2 := paddle{pos{700, 100}, 20, 100, 300, color{255, 255, 255}, 0}
	ball := ball{pos{300, 300}, 20, 400, 400, color{255, 255, 255}}
	keyState := sdl.GetKeyboardState()
	var frameStart time.Time
	var elapsedTime float32
	elapsedTime = 1
	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		if state == play {
			player1.update(keyState, elapsedTime)
			player2.updateP2(keyState, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}
				state = play

			}
		}
		clear(pixels)
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		// UPDATING
		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime/1000))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}
}
