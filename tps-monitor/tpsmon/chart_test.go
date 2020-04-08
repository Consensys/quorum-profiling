package tpsmon

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"image/color"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestSimpleComparisonChart(t *testing.T){

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	var nomX []string
	var nomX1 []string

	var trs = []TPSRecord{
		{"1-Jan 14:01","00:00:01",200, 100, 200},
		{"1-Jan 14:02","00:00:02",500, 300, 6000},
		{"1-Jan 14:03","00:00:03",600, 400, 800},
		{"1-Jan 14:04","00:00:04",400, 450, 9000},
		{"1-Jan 14:05","00:00:05",300, 600, 10000},
	}

	var trs1 = []TPSRecord{
		{"1-Jan 14:01","00:00:01",300, 100, 200},
		{"1-Jan 14:02","00:00:02",570, 300, 6000},

	}

	pts := make(plotter.XYs, len(trs))
	pts1 := make(plotter.XYs, len(trs1))

	p.Title.Text = "TPS Comparison"
	p.X.Label.Text = "Time Elapsed"
	p.Y.Label.Text = "TPS"

	var X float64
	X = 0

	var X1 float64
	X1 = 0

	for i, r := range trs {
		pts[i].Y = float64(r.tps)
		pts[i].X = X
		X += 1.0
		nomX = append(nomX, r.ltime)
	}

	for i, r := range trs1 {
		pts1[i].Y = float64(r.tps)
		pts1[i].X = X1
		X1 += 1.0
		nomX1 = append(nomX1, r.rtime)
	}

	err = plotutil.AddLinePoints(p, "quorum 2.4.0", pts, "quorum 1.9.7", pts1)

	if err != nil {
		panic(err)
	}
	p.NominalX(nomX...)
	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, "/Users/amalraj.manigmail.com/Downloads/tps.png"); err != nil {
		panic(err)
	}

}
func TestSimpleComparisonChart1(t *testing.T){
	pts1 := randomPoints(25*60)
	pts2 := randomPoints(25*60)

	t1 := time.Now()
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	var nomX []string
	var X float64

	X = 0
	for i, _ := range pts1 {
		pts1[i].X = X
		X += 1.0
		lbl := ""
		if i%100 == 0 {
			lbl = "X"+strconv.Itoa(i)
		}
		nomX = append(nomX, lbl)
	}

	X = 0
	for i, _ := range pts2 {
		pts2[i].X = X
		X += 1.0
	}

	p.Title.Text = "TPS Comparison"
	p.X.Label.Text = "Time Elapsed"
	p.Y.Label.Text = "TPS"

	p.Add(plotter.NewGrid())

	// Make a line plotter and set its style.
	l1, err := plotter.NewLine(pts1)
	if err != nil {
		panic(err)
	}
	l1.LineStyle.Width = vg.Points(1)
	l1.LineStyle.Color = color.RGBA{B: 255, A: 255}

	// Make a line plotter and set its style.
	l2, err := plotter.NewLine(pts2)
	if err != nil {
		panic(err)
	}
	l2.LineStyle.Width = vg.Points(1)
	l2.LineStyle.Color = color.RGBA{R: 255, A: 255}

	p.Add(l1, l2)
	p.Legend.Add("quorum 2.4.0", l1)
	p.Legend.Add("quorum 1.9.7", l2)


	p.NominalX(nomX...)
	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 5*vg.Inch, "/Users/amalraj.manigmail.com/Downloads/tps2.png"); err != nil {
		panic(err)
	}

	t2 := time.Now().Sub(t1).Milliseconds()
	t.Logf("time taken: %d\n", t2)

}


// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}