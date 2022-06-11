package captcha

import (
	"image/color"
	"testing"
)

func TestGenCaptcha(t *testing.T) {

	cpt := NewCaptcha(460, 200)
	cpt.SetBackgroundColor(color.RGBA{R: uint8(20), G: uint8(8), B: uint8(100), A: uint8(255)})

	image, err := cpt.GenCaptchaImage()
	t.Logf("GenCaptchaImage %s %+v", image, err)
}
