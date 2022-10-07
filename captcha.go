package captcha

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

type CaptchaStyle int32

const (
	CaptchaStyle_Normal   CaptchaStyle = 0
	CaptchaStyle_Behavior CaptchaStyle = 1
)

type Config struct {
	FontColors     []color.Color
	BackgroupColor color.Color //default white
	FontSize       int         //default 56
	MaxRotate      int         //default 30
	Style          CaptchaStyle
	Font           *truetype.Font //字体
}

type Captcha struct {
	Width  int
	Height int
	Config Config
}

type CaptchaResult struct {
	DrawRects   []DrawRect `json:"draw_rect"`
	Text        string     `json:"text"`
	ImageBase64 string     `json:"image_base64"`
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

//SetFont 设置字体颜色
func (cp *Captcha) SetFont(font *truetype.Font) {
	if font == nil {
		return
	}

	cp.Config.Font = font
}

//SetFontSize 字体大小
func (cp *Captcha) SetFontSize(fontSize int) {
	cp.Config.FontSize = fontSize
}

//SetFontSize 字体大小
func (cp *Captcha) SetSytle(style CaptchaStyle) {
	cp.Config.Style = style
}

//SetBackgroundColor 设置背景颜色
func (cp *Captcha) SetBackgroundColor(color color.Color) {
	cp.Config.BackgroupColor = color
}

//GenCaptchaImage 生成文字的验证码图片
func (cp *Captcha) GenCaptchaImage(text string) (CaptchaResult, error) {

	result := CaptchaResult{}

	c := Canvas{
		Width:  cp.Width,
		Height: cp.Height,
		NRGBA:  image.NewNRGBA(image.Rect(0, 0, cp.Width, cp.Height)),
		Config: cp.Config,
	}

	// text := "3567"
	c.DrawBackgroud()

	c.DrawLines(5)

	drawRects := c.DrawString(text)

	c.DrawCircles(120)

	// c.DrawBezierLines(5)

	imageBytes, err := encodingWithPng(c)

	if err != nil {
		return result, err
	}

	result.DrawRects = drawRects
	result.ImageBase64 = fmt.Sprintf("data:%s;base64,%s", "image/png", base64.StdEncoding.EncodeToString(imageBytes))
	result.Text = text

	// writeImageFile("./previews/test.png", imageBytes)
	return result, nil

}

//GenNormalRandomCaptcha 随机验证码 - 普通验证码
func (cp *Captcha) GenRandomNormalCaptcha(length int) (CaptchaResult, error) {
	var buff bytes.Buffer

	for i := 0; i < length; i++ {
		ix := rand.Intn(length)
		buff.WriteByte(Alphabets[ix])
	}

	return cp.GenCaptchaImage(buff.String())

}

//GenBehaviorCaptcha 生成中文点击验证码 - 点击行为验证码
func (cp *Captcha) GenBehaviorCNCaptcha() (CaptchaResult, error) {

	ix := rand.Intn(len(CNChars))
	return cp.GenCaptchaImage(CNChars[ix])

}

//ValidBehaviorCaptcha 验证行为点击验证码是否正确
func ValidBehaviorCaptcha(cr CaptchaResult, userPoints []Point) bool {
	//字符数量不相等
	return ValidBehaviorRects(cr.DrawRects, userPoints)
}

//ValidBehaviorRects 验证行为点是否在验证码矩形中
func ValidBehaviorRects(rects []DrawRect, userPoints []Point) bool {

	if len(userPoints) != len(rects) {
		return false
	}

	for i := 0; i < len(userPoints); i++ {
		realRect := rects[i]
		userPos := userPoints[i]

		if !PointInRect(userPos, realRect) {
			return false
		}
	}
	return true

}

//PointInRect 验证点是否在矩形中
func PointInRect(p Point, rect DrawRect) bool {
	if p.X >= rect.X && p.X <= rect.X+rect.Width && p.Y >= rect.Y && p.Y <= rect.Y+rect.Height {
		return true
	}
	return false
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

//writeImageFile tmp
func writeImageFile(filepath string, imageBytes []byte) {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(imageBytes)
}
