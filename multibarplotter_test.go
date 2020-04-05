package multibarplotter

import (
	"testing"
	"math/rand"

	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

var (
	black = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	blue  = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	red   = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green = color.RGBA{R: 0, G: 200, B: 100, A: 255}
)

func TestMultiBarPlotter(t *testing.T){
	rand.Seed(int64(0))
	b := getRandomData(15)

	color := make([]color.Color,4)
	color[0] = black
	color[1] = blue
	color[2] = red
	color[3] = green

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "MultiBar"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	bars,err := NewMultiBarPlotter(b,vg.Length(10),color)
	if err != nil {
		panic(err)
	}

	p.Add(bars)

	p.Legend.Add("Part 1",bars.GetSubLegend(0))
	p.Legend.Add("Part 2",bars.GetSubLegend(1))
	p.Legend.Add("Part 3",bars.GetSubLegend(2))
	p.Legend.Add("Part 4",bars.GetSubLegend(3))
	p.Legend.Top = true
	if err := p.Save(600, 600, "MultiBar.png"); err != nil {
		panic(err)
	}
}

func getRandomData(n int)Bars{
	newBars := make(Bars,n+3)
	for i := 0 ; i < n ; i++{
		newBars[i].Ymin = 0.0
		newBars[i].X = float64(rand.Intn(30))
		newBars[i].Y = make([]float64,4)
		firstFloat := rand.Float64()
		for j := range newBars[i].Y{
			cof := 0.0
			for cof == 0.0{
				cof = float64(rand.Intn(6))
			}
			newBars[i].Y[j] = firstFloat * cof
		}
	}
	newBars[14].X = 15.0
	return newBars
}