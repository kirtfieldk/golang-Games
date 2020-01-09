package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600

func clear(pixles []byte) {
	for i := range pixles {
		pixles[i] = 0
	}
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

func (ball *ball) update(leftPad *paddle, rightPad *paddle) {
	ball.x += ball.xv
	ball.y += ball.yv
	// Handle collisions
	if ball.y-float32(ball.radius) < 0 || ball.y+float32(ball.radius) > float32(winHeight) {
		ball.yv = -ball.yv
	}
	if ball.x < 0 || ball.x > float32(winWidth) {
		ball.x = 300
		ball.y = 300
	}
	// Ball hits the edhe of left paddle
	if ball.x < float32(leftPad.x)+float32(rightPad.w/2) {
		if ball.y > leftPad.y-float32(leftPad.h/2) && ball.y < leftPad.y+float32(leftPad.h/2) {
			ball.xv = -ball.xv
		}
	}
	if int(ball.x) > int(rightPad.x)-int(rightPad.w/2) {
		if ball.y > rightPad.y-float32(rightPad.h/2) && ball.y < rightPad.y+float32(rightPad.h/2) {
			ball.xv = -ball.xv
		}
	}
}

type paddle struct {
	pos
	w     int
	h     int
	color color
}

func (paddle *paddle) draw(pixels []byte) {
	startx := int(paddle.x) - paddle.w/2
	starty := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startx+x, starty+y, paddle.color, pixels)
		}
	}
}
func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= 10
	} else if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += 10
	}

}
func (paddle *paddle) updateP2(keyState []uint8) {
	if keyState[sdl.SCANCODE_A] != 0 {
		paddle.y -= 10
	} else if keyState[sdl.SCANCODE_Q] != 0 {
		paddle.y += 10
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

	player1 := paddle{pos{100, 100}, 20, 100, color{255, 255, 255}}
	player2 := paddle{pos{700, 100}, 20, 100, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 20, 3, 3, color{255, 255, 255}}
	keyState := sdl.GetKeyboardState()
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)
		player1.update(keyState)
		player2.updateP2(keyState)
		ball.update(&player1, &player2)
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		// UPDATING
		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		sdl.Delay(16)
	}
}
