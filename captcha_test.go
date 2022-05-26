package captcha

import (
	"testing"
)

func TestGenCaptcha(t *testing.T) {

	cpt := NewCaptcha(460, 200)
	image, err := cpt.GenCaptchaImage()
	t.Logf("GenCaptchaImage %s %+v", image, err)
}
