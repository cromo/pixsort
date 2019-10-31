package pixsort

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"math"
	"os"

	svg "github.com/ajstarks/svgo"
)

func LoadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}

func GetPoints(im image.Image) (int, int, []Point) {
	var result []Point
	bounds := im.Bounds()
	w := bounds.Size().X
	h := bounds.Size().Y
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := im.At(x, y)
			r, _, _, a := c.RGBA()
			// fmt.Printf("R: %d, A: %d\n", r, a)
			if r < 65535 && a > 0 {
				result = append(result, Point{x, y})
			}
		}
	}
	return w, h, result
}

func CreateFrame(m, w, h int, points []Point) *image.Paletted {
	var palette []color.Color
	palette = append(palette, color.RGBA{255, 255, 255, 255})
	palette = append(palette, color.RGBA{0, 0, 0, 255})
	im := image.NewPaletted(image.Rect(0, 0, w*m, h*m), palette)
	for _, point := range points {
		for y := 0; y < m; y++ {
			for x := 0; x < m; x++ {
				im.SetColorIndex(point.X*m+x, point.Y*m+y, 1)
			}
		}
	}
	return im
}

func SaveGIF(path string, m, w, h int, points []Point) error {
	g := gif.GIF{}
	g.Image = append(g.Image, CreateFrame(m, w, h, nil))
	g.Delay = append(g.Delay, 10)
	for i := range points {
		g.Image = append(g.Image, CreateFrame(m, w, h, points[:i+1]))
		g.Delay = append(g.Delay, 10)
	}
	g.Image = append(g.Image, CreateFrame(m, w, h, points))
	g.Delay = append(g.Delay, 100)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gif.EncodeAll(file, &g)
}

func SaveSVG(path string, w, h int, points []Point) error {
	segments := groupSegments(points)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	canvas := svg.New(file)
	canvas.Start(w, h)
	// canvas.Path("M 1 1", "fill:none; stroke: red")
	// canvas.Path("M 1 2 L 2 2", "fill:none; stroke: blue")
	// canvas.Path("M 1 3 L 2 3 2 3", "fill:none; stroke: green")
	str := ""
	for _, s := range segments {
		// canvas.Circle(s[0].X, s[0].Y, 1, "fill:none;stroke:gray")
		str += fmt.Sprintf("\nM %d %d", s[0].X, s[0].Y)
		for _, point := range s[1:] {
			str += fmt.Sprintf(" %d %d", point.X, point.Y)
		}
	}
	canvas.Path(str, "fill:none;stroke:black;stroke-linecap:square")
	canvas.End()
	return nil
}

func groupSegments(ps []Point) [][]Point {
	ss := make([][]Point, 0)
	s := make([]Point, 0)
	for _, p := range ps {
		if len(s) == 0 {
			s = append(s, p)
			continue
		}
		last := s[len(s)-1]
		if math.Abs(float64(last.X-p.X)) <= 1 && math.Abs(float64(last.Y-p.Y)) <= 1 {
			s = append(s, p)
		} else {
			ss = append(ss, s)
			s = make([]Point, 1)
			s[0] = p
		}
	}
	if len(s) != 0 {
		ss = append(ss, s)
	}
	return ss
}
