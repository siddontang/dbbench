package plot

import (
	"fmt"
	"path"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/siddontang/dbbench/pkg/stats"
)

type barValue struct {
	name  string
	value float64
}

type lineValue struct {
	name   string
	points plotter.XYs
}

func newBarValue(op string, tp stats.StatType, s *stats.DBStat) barValue {
	stat, ok := s.Summary[op]
	value := 0.0
	if ok {
		value = stat.Value(tp)
	}

	return barValue{
		name:  s.Name,
		value: value}
}

func newLineValue(op string, tp stats.StatType, s *stats.DBStat) lineValue {
	stats, ok := s.Progress[op]
	if !ok {
		return lineValue{
			name:   s.Name,
			points: plotter.XYs{},
		}
	}

	xys := make(plotter.XYs, 0, len(stats))
	for i, stat := range stats {
		xys = append(xys, plotter.XY{
			float64(i),
			stat.Value(tp),
		})
	}
	return lineValue{
		name:   s.Name,
		points: xys,
	}
}

// Meta is the meta to plot the charts
type Meta struct {
	OutputDir string

	// StatTypes is the array of statistics we want to plot
	StatTypes []stats.StatType

	Workload string
	OP       string

	XLength int
	YLength int
}

func (m Meta) outputPath(tp stats.StatType, suffix string) string {
	var baseName string
	if m.OP == "" {
		baseName = fmt.Sprintf("%s_%s_%s.png", m.Workload, tp, suffix)
	} else {
		baseName = fmt.Sprintf("%s_%s_%s_%s.png", m.Workload, m.OP, tp, suffix)
	}
	return path.Join(m.OutputDir, m.Workload, baseName)
}

func (m Meta) name() string {
	if m.OP == "" {
		return m.Workload
	}

	return fmt.Sprintf("%s %s", m.Workload, m.OP)
}

func plotBarCharts(m Meta, stats stats.DBStats, tp stats.StatType, outputFile string) error {
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = m.name()
	p.HideX()
	p.Y.Label.Text = tp.String()

	w := vg.Points(20)

	barValues := make([]barValue, len(stats))
	for i := 0; i < len(stats); i++ {
		barValues[i] = newBarValue(m.OP, tp, stats[i])
	}

	bars := make([]*plotter.BarChart, len(stats))
	for i := 0; i < len(stats); i++ {
		bars[i], err = plotter.NewBarChart(plotter.Values{
			barValues[i].value,
		}, w)
		bars[i].LineStyle.Width = vg.Length(0)
		bars[i].Color = plotutil.Color(i)
		bars[i].Offset = w * vg.Length(i-len(stats)/2)
	}

	vals := make([]plot.Plotter, len(stats))
	for i := 0; i < len(stats); i++ {
		vals[i] = bars[i]
	}
	p.Add(vals...)
	for i := 0; i < len(stats); i++ {
		p.Legend.Add(barValues[i].name, bars[i])
	}
	p.Legend.Top = true

	return p.Save(vg.Length(m.XLength)*vg.Inch, vg.Length(m.YLength)*vg.Inch, outputFile)
}

func plotLineCharts(m Meta, stats stats.DBStats, tp stats.StatType, outputFile string) error {
	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = m.name()
	p.X.Label.Text = "time"
	p.Y.Label.Text = tp.String()

	vals := make([]interface{}, 0, len(stats)*2)
	for i := 0; i < len(stats); i++ {
		lineValue := newLineValue(m.OP, tp, stats[i])
		vals = append(vals, lineValue.name, lineValue.points)
	}

	err = plotutil.AddLinePoints(p, vals...)
	if err != nil {
		return err
	}

	return p.Save(vg.Length(m.XLength)*vg.Inch, vg.Length(m.YLength)*vg.Inch, outputFile)
}

// PlotCharts plots the charts from the DBStats
func PlotCharts(m Meta, stats stats.DBStats) error {
	for _, tp := range m.StatTypes {
		outputFile := m.outputPath(tp, "prog")
		if err := plotLineCharts(m, stats, tp, outputFile); err != nil {
			return err
		}

		outputFile = m.outputPath(tp, "total")
		if err := plotBarCharts(m, stats, tp, outputFile); err != nil {
			return err
		}
	}

	return nil
}
