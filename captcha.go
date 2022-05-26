package captcha

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	// "image/color"
	"image/png"
	"os"
)

type Captcha struct {
	Width  int
	Height int
}

func NewCaptcha(width, height int) *Captcha {
	return &Captcha{Width: width, Height: height}
}

func (cp *Captcha) GenCaptchaImage() (string, error) {
	c := Canvas{
		Width:  cp.Width,
		Height: cp.Height,
		NRGBA:  image.NewNRGBA(image.Rect(0, 0, cp.Width, cp.Height)),
	}

	c.DrawLines(5)

	c.DrawString("测试阿斯蒂")

	c.DrawCircles(120)

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
