package captcha

import (
	"image"
	"image/color"
	"math"
)

type Palette struct {
	*image.Paletted
}

func NewPalette(r image.Rectangle, p color.Palette) *Palette {
	return &Palette{
		image.NewPaletted(r, p),
	}
}

//Rotate 旋转角度
func (p *Palette) Rotate(angle int) {
	width := p.Bounds().Max.X
	height := p.Bounds().Max.Y
	r := width / 2
	retImg := image.NewPaletted(image.Rect(0, 0, width, height), p.Palette)
	for x := 0; x <= retImg.Bounds().Max.X; x++ {
		for y := 0; y <= retImg.Bounds().Max.Y; y++ {
			tx, ty := p.angleSwapPoint(float64(x), float64(y), float64(r), float64(angle))
			retImg.SetColorIndex(x, y, p.ColorIndexAt(int(tx), int(ty)))
		}
	}

	nW := retImg.Bounds().Max.X
	nH := retImg.Bounds().Max.Y
	for x := 0; x < nW; x++ {
		for y := 0; y < nH; y++ {
			p.SetColorIndex(x, y, retImg.ColorIndexAt(x, y))
		}
	}
}

//angleSwapPoint 坐标转换
func (p *Palette) angleSwapPoint(x, y, r, angle float64) (tarX, tarY float64) {
	x -= r
	y = r - y
	sinVal := math.Sin(angle * (math.Pi / 180))
	cosVal := math.Cos(angle * (math.Pi / 180))
	tarX = x*cosVal + y*sinVal
	tarY = -x*sinVal + y*cosVal
	tarX += r
	tarY = r - tarY
	return
}
