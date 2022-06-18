package captcha

import (
	"image/color"
	"testing"
)

func TestGenCaptcha(t *testing.T) {

	cpt := NewCaptcha(260, 100)
	// cpt.SetBackgroundColor(color.RGBA{R: uint8(20), G: uint8(8), B: uint8(100), A: uint8(255)})
	// cpt.Config.Style = CaptchaStyle_Behavior
	result, err := cpt.GenCaptchaImage("3567")
	t.Logf("GenCaptchaImage: %+v %+v", result, err)
}

func TestGenNormalCaptcha(t *testing.T) {

	cpt := NewCaptcha(260, 100)
	cpt.SetBackgroundColor(color.RGBA{R: uint8(20), G: uint8(8), B: uint8(100), A: uint8(255)})
	// cpt.Config.Style = CaptchaStyle_Behavior
	result, err := cpt.GenRandomNormalCaptcha(4)
	t.Logf("GenCaptchaImage: %+v %+v", result, err)
}

func TestGenGenBehaviorCaptcha(t *testing.T) {

	cpt := NewCaptcha(260, 100)
	cpt.SetBackgroundColor(color.RGBA{R: uint8(20), G: uint8(8), B: uint8(100), A: uint8(255)})
	// cpt.Config.Style = CaptchaStyle_Behavior
	cpt.Config.MaxRotate = 20
	result, err := cpt.GenBehaviorCaptcha()
	t.Logf("GenCaptchaImage: %+v %+v", result, err)
}

func TestPointInRect(t *testing.T) {
	rect := DrawRect{X: 10, Width: 50, Y: 50, Height: 60}
	point := Point{X: 10, Y: 40}

	if PointInRect(point, rect) {
		t.Errorf("PointInRect error")
	}

	point = Point{X: 10, Y: 50}
	if !PointInRect(point, rect) {
		t.Errorf("PointInRect error")
	}

	point = Point{X: 10, Y: 60}
	if !PointInRect(point, rect) {
		t.Errorf("PointInRect error")
	}

	point = Point{X: 10, Y: 110}
	if !PointInRect(point, rect) {
		t.Errorf("PointInRect error")
	}

	point = Point{X: 30, Y: 110}
	if !PointInRect(point, rect) {
		t.Errorf("PointInRect error")
	}

	point = Point{X: 60, Y: 110}
	if !PointInRect(point, rect) {
		t.Errorf("PointInRect error")
	}
}

func TestValidBehaviorCaptcha(t *testing.T) {
	cr := CaptchaResult{}
	cr.DrawRects = []DrawRect{
		DrawRect{X: 10, Y: 10, Width: 60, Height: 60},
		DrawRect{X: 20, Y: 60, Width: 60, Height: 60},
		DrawRect{X: 60, Y: 90, Width: 60, Height: 60},
	}

	ps := []Point{
		Point{X: 15, Y: 8},
	}

	if ValidBehaviorCaptcha(cr, ps) {
		t.Errorf("invalid error")
	}
	if !ValidBehaviorCaptcha(cr, ps) {
		t.Logf("invalid captcha")
	}

	ps = []Point{
		Point{X: 15, Y: 8},
		Point{X: 15, Y: 8},
		Point{X: 15, Y: 8},
	}
	if ValidBehaviorCaptcha(cr, ps) {
		t.Errorf("invalid error")
	}

	ps = []Point{
		Point{X: 15, Y: 12},
		Point{X: 20, Y: 62},
		Point{X: 60, Y: 8},
	}
	if ValidBehaviorCaptcha(cr, ps) {
		t.Errorf("invalid error")
	}

	ps = []Point{
		Point{X: 15, Y: 12},
		Point{X: 20, Y: 62},
		Point{X: 60, Y: 98},
	}
	if ValidBehaviorCaptcha(cr, ps) {
		t.Logf("valid captcha")
	}
}
