package multibarplotter

import (
	gosort "github.com/Rakiiii/goSort"
	"math"
)
type Barer interface{
	//Len returns amount of bars
	Len() int

	//BarVal returns bar characteristics
	BarVal(int)(float64,[]float64,float64)

	//Xmin returns min X 
	Xmin()float64
}

//Bars implements Barer interface using slice
type Bars []struct{
	X float64
	Y []float64
	Ymin float64
}

//Len implements Len method of Barer interface 
func (b Bars) Len()int{
	return len(b)
}

//BarVal implements BarVal method of Barer interface
func (b Bars)BarVal(i int)(float64 , []float64,float64){
	tmp := make([]float64,len(b[i].Y))
	copy(tmp,b[i].Y)
	return b[i].X, gosort.QuicksortFloat64(tmp),b[i].Ymin
}

//Xmin implements Xmin method of Barer interface
func (b Bars)Xmin()float64{
	xmin := math.Inf(1)
	for _,i := range b{
		if i.X < xmin{
			xmin = i.X
		}
	} 
	return xmin
}

//CopyBars copies an Barer
func CopyBars(data Barer)Bars{
	cpy :=make(Bars,data.Len())
	for i := range cpy{
		cpy[i].X,cpy[i].Y,cpy[i].Ymin = data.BarVal(i)
	}
	return cpy
}


type CofNormalizer struct{
	cof float64
}

func (c CofNormalizer)Normalize(min, max, x float64) float64{
	return c.cof*(x - min) / (max - min)
}