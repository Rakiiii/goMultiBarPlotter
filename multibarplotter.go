package multibarplotter
import (
	"image/color"
	"errors"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/plotter"
)

//SubLegend struct sub struct for adding to legend
type SubLegend struct{
	color.Color
	draw.LineStyle
}

//Thumbnail implement Thumbnail method of thumbnailer interface
func (s SubLegend)Thumbnail(c *draw.Canvas){
	pts := []vg.Point{
		{c.Min.X, c.Min.Y},
		{c.Min.X, c.Max.Y},
		{c.Max.X, c.Max.Y},
		{c.Max.X, c.Min.Y},
	}
	poly := c.ClipPolygonY(pts)
	c.FillPolygon(s.Color, poly)

	pts = append(pts, vg.Point{X: c.Min.X, Y: c.Min.Y})
	outline := c.ClipLinesY(pts)
	c.StrokeLines(s.LineStyle, outline...)
}

type MultiBarPlotter struct{
	Bars

	//Width is width of the bars
	Width vg.Length

	//Colors is slice of color for set of bars part
	Colors []color.Color

	//LineStyle is the style of the outline of the bars
	draw.LineStyle
}

//NewMultiBarPlotter creates new multiBarPlotter with @br data set @width width and @cl color set 
func NewMultiBarPlotter(br Barer, width vg.Length, cl []color.Color)(*MultiBarPlotter,error){
	if width <= 0{
		return nil, errors.New("Width parameter was not positive")
	}

	clcopy := make([]color.Color,len(cl))
	copy(clcopy,cl)
	bars := CopyBars(br)

	return &MultiBarPlotter{
		Bars : bars,
		Width : width,
		Colors : clcopy,
		LineStyle : plotter.DefaultLineStyle,
	},nil
}



func (m *MultiBarPlotter)Plot(c draw.Canvas, plt *plot.Plot){
	trX , trY := plt.Transforms(&c)

	for _,bar := range m.Bars {
		if len(bar.Y) <= 0{
			continue
		}
		for i,y := range bar.Y{
			var ymin float64
			if i == 0{
				ymin = bar.Ymin
			}else{
				ymin = bar.Y[i-1]
			}
			//left X cord of bar on canvas
			catMin := trX(bar.X)-m.Width/2

			//check for possible draw
			if !c.ContainsX(catMin) {
				continue
			}
			//right X cord of bar on canvas
			catMax := catMin + m.Width
			
			//start Y pos of bar on canvas
			valMin := trY(ymin)
			//end Y pos of bar on canvas
			valMax := trY(y)

			pts := []vg.Point{
				{catMin,valMin},
				{catMin,valMax},
				{catMax,valMax},
				{catMax,valMin},
			}

			poly := c.ClipPolygonY(pts)
			c.FillPolygon(m.Colors[i%len(m.Colors)],poly)
			
		} 

		pts := []vg.Point{
			{trX(bar.X)-m.Width/2,trY(bar.Ymin)},
			{trX(bar.X)-m.Width/2,trY(bar.Y[len(bar.Y)-1])},
			{trX(bar.X)+m.Width/2,trY(bar.Y[len(bar.Y)-1])},
			{trX(bar.X)+m.Width/2,trY(bar.Ymin)},
			{trX(bar.X)-m.Width/2,trY(bar.Ymin)},
		}

		outline := c.ClipLinesY(pts)
		c.StrokeLines(m.LineStyle,outline...)
	}
	plt.X.Scale = plot.LinearScale{}
}

// DataRange implements the plot.DataRanger interface.
func (b *MultiBarPlotter) DataRange() (xmin, xmax, ymin, ymax float64) {
	ymin = 0
	xmin = b.Xmin()
	ymax = math.Inf(-1)
	xmax = math.Inf(-1)
	for _,bar := range b.Bars{
		if bar.X > xmax {
			xmax = bar.X
		}
		counter := 0.0
		for _,y := range bar.Y{
			counter += y
		}
		if counter > ymax {
			ymax = counter
		}
	}
	return xmin-1,xmax+1,ymin,ymax
}

//GetSubLegend return implementor of thumbnailer interface for adding to legend color with number @i 
func (m *MultiBarPlotter)GetSubLegend(i int)SubLegend{
	return SubLegend{
		Color:m.Colors[i],
		LineStyle:m.LineStyle,
	}
}

