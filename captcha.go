package captcha

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Config struct {
	FontColors     []color.Color
	BackgroupColor color.Color
	FontSize       int
	MaxRotate      int
}

type Captcha struct {
	Width  int
	Height int
	Config Config
}

func NewCaptcha(width, height int) *Captcha {
	cp := Captcha{Width: width, Height: height}

	config := Config{}
	config.FontColors = []color.Color{
		color.RGBA{R: uint8(66), G: uint8(153), B: uint8(244), A: uint8(255)},
		color.RGBA{R: uint8(234), G: uint8(67), B: uint8(53), A: uint8(255)},
		color.RGBA{R: uint8(251), G: uint8(188), B: uint8(5), A: uint8(255)},
		color.RGBA{R: uint8(52), G: uint8(168), B: uint8(83), A: uint8(255)},
	}
	config.BackgroupColor = color.White
	config.FontSize = 56
	config.MaxRotate = 30
	cp.Config = config
	return &cp

}

//SetFontColors 设置字体颜色
func (cp *Captcha) SetFontColors(colors ...color.Color) {
	if len(colors) == 0 {
		return
	}

	cp.Config.FontColors = cp.Config.FontColors[:0]
	for _, cr := range colors {
		cp.Config.FontColors = append(cp.Config.FontColors, cr)
	}

}

//SetBackgroundColor 设置背景颜色
func (cp *Captcha) SetBackgroundColor(color color.Color) {
	cp.Config.BackgroupColor = color
}

func (cp *Captcha) GenCaptchaImage() (string, error) {

	c := Canvas{
		Width:  cp.Width,
		Height: cp.Height,
		NRGBA:  image.NewNRGBA(image.Rect(0, 0, cp.Width, cp.Height)),
		Config: cp.Config,
	}

	c.DrawBackgroud()

	c.DrawLines(5)

	c.DrawString("3567中")

	c.DrawCircles(120)

	// c.DrawBezierLines(5)

	imageBytes, err := encodingWithPng(c)

	if err != nil {
		return "", err
	}

	writeImageFile("./test.png", imageBytes)
	return fmt.Sprintf("data:%s;base64,%s", "image/png", base64.StdEncoding.EncodeToString(imageBytes)), nil

}

func encodingWithPng(img image.Image) (encodeResult []byte, err error) {

	var buf bytes.Buffer

	if err := png.Encode(&buf, img); err != nil {
		return encodeResult, err
	}

	encodeResult = buf.Bytes()

	buf.Reset()
	return encodeResult, nil

}

func writeImageFile(filepath string, imageBytes []byte) {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(imageBytes)
}
