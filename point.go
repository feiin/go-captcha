package captcha

type Point struct {
	X int
	Y int
}

type DrawRect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func NewPoint(x int, y int) *Point {
	return &Point{X: x, Y: y}
}

type AreaPoint struct {
	MinX, MaxX, MinY, MaxY int
}
