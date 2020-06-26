package tpsmon

import (
	"fmt"
	"image/color"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

type TPSServer struct {
	tm   *TPSMonitor
	port int
}

func NewTPSServer(tm *TPSMonitor, port int) TPSServer {
	s := TPSServer{
		tm:   tm,
		port: port,
	}
	go s.Start()
	return s

}

func (s TPSServer) Start() {
	http.HandleFunc("/tpsdata", s.PrintTPSData)
	http.HandleFunc("/tpschart", s.PrintTPSChart)
	log.Infof("started tps monitor server at port %d", s.port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.port), nil))
}

func (s TPSServer) PrintTPSData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "localTime,refTime,TPS,TxnCount,BlockCount\n")
	for _, v := range s.tm.tpsRecs {
		fmt.Fprintf(w, "%s", v.ReportString())
	}
}

func (s TPSServer) PrintTPSChart(w http.ResponseWriter, r *http.Request) {

	rmap := r.URL.Query()

	iw := rmap.Get("iw")
	ih := rmap.Get("ih")
	var cH, cW int
	cH = 15
	cW = 20

	if iw != "" {
		if tw, err := strconv.Atoi(iw); err == nil {
			cW = tw
		} else {
			log.Error("invalid canvas width", iw)
		}
	}
	if ih != "" {
		if th, err := strconv.Atoi(ih); err == nil {
			cH = th
		} else {
			log.Error("invalid canvas height", ih)
		}
	}

	if cH > 25 {
		cH = 25
	}

	if cW > 50 {
		cW = 50
	}

	t1 := time.Now()
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	tpsRecLen := len(s.tm.tpsRecs)
	pts := make(plotter.XYs, tpsRecLen)
	var nomX []string

	p.Title.Text = "TPS Comparison"
	p.X.Label.Text = "Time Elapsed"
	p.Y.Label.Text = "TPS"

	var X float64
	X = 0

	for i, r := range s.tm.tpsRecs {
		pts[i].Y = float64(r.tps)
		pts[i].X = X
		X += 1.0
		if tpsRecLen > 100 {
			if i%50 == 0 {
				nomX = append(nomX, r.rtime)
			} else {
				nomX = append(nomX, "")
			}
		} else {
			nomX = append(nomX, r.ltime)
		}
	}

	p.Add(plotter.NewGrid())

	// Make a line plotter and set its style.
	l1, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	l1.LineStyle.Width = vg.Points(1)
	l1.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(l1)
	p.Legend.Add("quorum", l1)

	p.NominalX(nomX...)

	// Create a Canvas for writing SVG images.
	c := vgsvg.New(vg.Length(cW)*vg.Inch, vg.Length(cH)*vg.Inch)

	// Draw to the Canvas.
	p.Draw(draw.New(c))

	// Write the Canvas to a io.Writer (in this case, os.Stdout).
	if _, err := c.WriteTo(w); err != nil {
		log.Errorf("error - writing image to responseWriter. %v", err)
	}

	t2 := time.Now().Sub(t1).Milliseconds()
	log.Infof("time taken: %d\n", t2)
}
