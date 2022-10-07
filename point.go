package captcha

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type DrawRect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func NewPoint(x int, y int) *Point {
	return &Point{X: x, Y: y}
}

type AreaPoint struct {
	MinX, MaxX, MinY, MaxY int
}
