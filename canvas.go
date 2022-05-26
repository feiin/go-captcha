package captcha

import (
	// "fmt"

	"embed"
	// "fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	// "math"
	"math/rand"
	"time"
)

//go:embed fonts/*
var fontFS embed.FS

var defaultFont *truetype.Font

func init() {

	fontdata, _ := fontFS.ReadFile("fonts/SourceHanSerifCN-Light.ttf")
	defaultFont, _ = freetype.ParseFont(fontdata)

	rand.Seed(time.Now().UnixNano())

	// fmt.Printf("defaultFont %+v :err:%+v", defaultFont, err)
}

type Canvas struct {
	*image.NRGBA
	Width  int
	Height int
}

func (c *Canvas) DrawLines(lineCount int) {

	for lineCount > 0 {

		p1 := Point{X: rand.Intn(c.Width), Y: rand.Intn(c.Height)}
		p2 := Point{X: rand.Intn(c.Width), Y: rand.Intn(c.Height)}

		c.DrawLine(p1, p2, randomColor())
		lineCount--
	}
}

//DrawLine 画干扰线
func (c *Canvas) DrawLine(p1, p2 Point, color color.Color) {

	var k float64 = 0
	b := p1.Y
	if (p2.X - p1.X) != 0 {
		k = float64(p2.Y-p1.Y) / float64(p2.X-p1.X)
	}

	b = p1.Y - int(k*float64(p1.X))

	sx := 1
	offsetX := p1.X

	if p2.X < p1.X {
		sx = -1
	}

	for {

		offsetY := int(k*float64(offsetX)) + b
		c.Set(offsetX, offsetY, color)

		if offsetX == p2.X {
			break
		}
		offsetX = offsetX + sx

	}
}

func (c *Canvas) DrawString(text string) {

	for _, ch := range text {

		dc := freetype.NewContext()
		dc.SetDPI(float64(72))
		dc.SetFont(defaultFont)
		dc.SetClip(c.Bounds())
		dc.SetDst(c)

		// 文字大小
		dc.SetFontSize(float64(56))

		// 文字颜色
		fontColor := image.NewUniform(randomColor())
		dc.SetSrc(fontColor)

		pos := c.randomFontPosition(56)
		// 画文本
		pt := freetype.Pt(pos.X, pos.Y) // 字出现的位置
		_, err := dc.DrawString(string(ch), pt)
		if err != nil {
			panic(err)
		}
	}

}

func (c *Canvas) randomFontPosition(fontSize int) Point {
	minX := fontSize
	minY := fontSize

	maxX := c.Width - fontSize
	maxY := c.Height - fontSize

	x := randomInt(minX, maxX)
	y := randomInt(minY, maxY)

	return Point{x, y}

}

//randomInt 返回随机数 [min,max)
func randomInt(min, max int) int {
	if min >= max || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

func randomColor() color.RGBA {
	red := rand.Intn(256)
	green := rand.Intn(256)
	blue := rand.Intn(256)

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

//DrawCircles 随机产生circleCount干扰点
func (c *Canvas) DrawCircles(circleCount int) {

	for circleCount > 0 {

		cc := randomColor()
		x := rand.Intn(c.Width)
		y := rand.Intn(c.Height)

		r := rand.Intn(3)

		fillCircle := rand.Intn(5) > 1

		c.DrawCircle(x, y, r, fillCircle, cc)

		circleCount = circleCount - 1

	}

}

func (c *Canvas) drawCircle8(xc, yc, x, y int, cc color.Color) {
	c.Set(xc+x, yc+y, cc)
	c.Set(xc-x, yc+y, cc)
	c.Set(xc+x, yc-y, cc)
	c.Set(xc-x, yc-y, cc)
	c.Set(xc+y, yc+x, cc)
	c.Set(xc-y, yc+x, cc)
	c.Set(xc+y, yc-x, cc)
	c.Set(xc-y, yc-x, cc)
}

// DrawCircle 画圆
func (c *Canvas) DrawCircle(xc, yc, r int, fill bool, cc color.Color) {
	size := c.Bounds().Size()
	if xc+r < 0 || xc-r >= size.X || yc+r < 0 || yc-r >= size.Y {
		return
	}
	x, y, d := 0, r, 3-2*r
	for x <= y {
		if fill {
			for yi := x; yi <= y; yi++ {
				c.drawCircle8(xc, yc, x, yi, cc)
			}
		} else {
			c.drawCircle8(xc, yc, x, y, cc)
		}
		if d < 0 {
			d = d + 4*x + 6
		} else {
			d = d + 4*(x-y) + 10
			y--
		}
		x++
	}
}
