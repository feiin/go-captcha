package captcha

import (
	// "fmt"
	"image"
	"image/draw"

	"embed"
	// "fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	// "image"
	"image/color"
	"math"
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
	Config Config
}

//DrawBackgroud 设置画板背景
func (c *Canvas) DrawBackgroud() {
	backgroupColor := image.NewUniform(c.Config.BackgroupColor)
	draw.Draw(c, c.Bounds(), backgroupColor, image.ZP, draw.Over)
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
	var b float64 = 0
	if (p2.X - p1.X) != 0 {
		k = float64(p2.Y-p1.Y) / float64(p2.X-p1.X)
	}

	b = float64(p1.Y) - k*float64(p1.X)

	sx := 1
	offsetX := p1.X

	if p2.X < p1.X {
		sx = -1
	}

	if k == 0 && p2.X == p1.X {

		tmpY := p1.Y

		if p2.Y < p1.Y {
			sx = -1
		}

		for {

			c.Set(offsetX, tmpY, color)

			if tmpY == p2.Y {
				break
			}

			tmpY = tmpY + sx

		}

		return
	}

	for {

		offsetY := int(k*float64(offsetX) + b)
		c.Set(offsetX, offsetY, color)

		if offsetX == p2.X {
			break
		}
		offsetX = offsetX + sx

	}
}

//DrawBezierLines 画贝塞尔干扰线
func (c *Canvas) DrawBezierLines(lineCount int) {

	for lineCount > 0 {

		lineColor := randomColor()

		x1 := rand.Intn(c.Width / 4)
		y1 := randomInt(c.Height/4, int(float64(c.Height)*0.9))

		x2 := randomInt(c.Width/2, int(float64(c.Width)*0.9))
		y2 := randomInt(c.Height/4, int(float64(c.Height)*0.9))

		cx := randomInt(c.Width/4, int(float64(c.Width)*0.7))
		cy := randomInt(c.Height/4, int(float64(c.Height)*0.6))

		// B(t) = (1-t)P0 + 2t(1-t)P1 + tP2,   t∈[0,1]
		// x = Math.pow(1-t, 2) * x1 + 2 * t * (1-t) * cx + Math.pow(t, 2) * x2
		// y = Math.pow(1-t, 2) * y1 + 2 * t * (1-t) * cy + Math.pow(t, 2) * y2
		var t float64
		for t < 1 {

			x := math.Pow(1-t, 2.0)*float64(x1) + 2*t*(1-t)*float64(cx) + math.Pow(t, 2.0)*float64(x2)
			y := math.Pow(1-t, 2)*float64(y1) + 2*t*(1-t)*float64(cy) + math.Pow(t, 2)*float64(y2)

			c.Set(int(x), int(y), lineColor)
			t = t + 0.001
		}

		lineCount--
	}

}

func (c *Canvas) DrawString(text string) {

	var drawPos []DrawPos
	for _, ch := range text {

		pos := c.randomFontPosition(56)

		fontImg, areaPoint := c.DrawFont(string(ch), c.Config.FrontColors)

		minX := areaPoint.MinX
		maxX := areaPoint.MaxX
		minY := areaPoint.MinY
		maxY := areaPoint.MaxY
		width := maxX - minX
		height := maxY - minY
		nW := fontImg.Bounds().Max.X
		nH := fontImg.Bounds().Max.Y
		for x := 0; x < nW; x++ {
			for y := 0; y < nH; y++ {
				co := fontImg.At(x, y)
				if _, _, _, a := co.RGBA(); a > 0 {
					if x >= minX && x <= maxX && y >= minY && y <= maxY {
						c.Set(pos.X+(x-minX), pos.Y-height+(y-minY), fontImg.At(x, y))
					}
				}
			}
		}

		var dp DrawPos
		dp.X = minX + pos.X
		dp.Y = pos.Y - height
		dp.Width = width
		dp.Height = height

		drawPos = append(drawPos, dp)

	}

	// lineColor := randomColor()

	// for _, linePos := range drawPos {
	// 	p1 := Point{X: linePos.X, Y: linePos.Y}
	// 	p2 := Point{X: linePos.X + linePos.Width, Y: linePos.Y}

	// 	p3 := Point{X: linePos.X, Y: linePos.Y + linePos.Height}
	// 	p4 := Point{X: linePos.X + linePos.Width, Y: linePos.Y + linePos.Height}
	// 	c.DrawLine(p1, p2, lineColor)
	// 	c.DrawLine(p1, p3, lineColor)
	// 	c.DrawLine(p3, p4, lineColor)
	// 	c.DrawLine(p2, p4, lineColor)

	// }
	// fmt.Printf("drawpos %+v %s", drawPos, text)

}

//DrawFont 绘制验证码字
func (c *Canvas) DrawFont(fontText string, fontColors []color.Color) (*Palette, *AreaPoint) {
	fontSize := 56

	rand.Seed(time.Now().UnixNano())

	fntColorIndex := rand.Intn(len(fontColors))

	rColor := fontColors[fntColorIndex]
	p := []color.Color{
		color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x00},
		rColor,
	}

	canvas := NewPalette(image.Rect(0, 0, fontSize, fontSize), p)

	dc := freetype.NewContext()
	dc.SetDPI(float64(72))
	dc.SetFont(defaultFont)
	dc.SetClip(canvas.Bounds())
	dc.SetDst(canvas)

	// 文字大小
	dc.SetFontSize(float64(56))

	// 文字颜色
	fontColor := image.NewUniform(rColor)
	dc.SetSrc(fontColor)

	// 画文本
	pt := freetype.Pt(0, fontSize) // 字出现的位置
	_, err := dc.DrawString(string(fontText), pt)

	if err != nil {
		panic(err)
	}

	//旋转角度
	canvas.Rotate(randomInt(-30, 30))
	ap := c.calcImageSpace(canvas)
	return canvas, ap

}

func (c *Canvas) calcImageSpace(pa *Palette) *AreaPoint {
	nW := pa.Bounds().Max.X
	nH := pa.Bounds().Max.Y
	// 计算裁剪的最小及最大的坐标
	minX := nW
	maxX := 0
	minY := nH
	maxY := 0
	for x := 0; x < nW; x++ {
		for y := 0; y < nH; y++ {
			co := pa.At(x, y)
			if _, _, _, a := co.RGBA(); a > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}

				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	minX = int(math.Max(0, float64(minX-2)))
	maxX = int(math.Min(float64(nW), float64(maxX+2)))
	minY = int(math.Max(0, float64(minY-2)))
	maxY = int(math.Min(float64(nH), float64(maxY+2)))

	return &AreaPoint{
		minX,
		maxX,
		minY,
		maxY,
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
