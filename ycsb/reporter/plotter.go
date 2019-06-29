package main

import (
	"fmt"
	"path"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func (s *dbStat) BarValue(op string, tp StatType) barValue {
	stat, ok := s.summary[op]
	value := 0.0
	if ok {
		value = stat.Value(tp)
	}

	return barValue{
		name:  s.name,
		value: value}
}

func (s *dbStat) LineValue(op string, tp StatType) lineValue {
	stats, ok := s.progress[op]
	if !ok {
		return lineValue{
			name:   s.name,
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
		name:   s.name,
		points: xys,
	}
}

type plotMeta struct {
	workload string
	op       string
	statType StatType
}

func (m plotMeta) OutputPath(suffix string) string {
	var baseName string
	if m.workload == "load" {
		baseName = fmt.Sprintf("load_%s_%s.png", m.statType, suffix)
	} else {
		baseName = fmt.Sprintf("%s_%s_%s.png", m.workload, m.statType, suffix)
	}
	return path.Join(outputDir, m.workload, baseName)
}

func (m plotMeta) Name() string {
	if m.workload == "load" {
		return "load"
	}

	return fmt.Sprintf("%s %s", m.workload, m.op)
}

type barValue struct {
	name  string
	value float64
}

type lineValue struct {
	name   string
	points plotter.XYs
}

var (
	plotXLength int
	plotYLength int
)

func plotBarCharts(m plotMeta, stats dbStats, outputFile string) error {
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = m.Name()
	p.HideX()
	p.Y.Label.Text = m.statType.String()

	w := vg.Points(20)

	barValues := make([]barValue, len(stats))
	for i := 0; i < len(stats); i++ {
		barValues[i] = stats[i].BarValue(m.op, m.statType)
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

	return p.Save(vg.Length(plotXLength)*vg.Inch, vg.Length(plotYLength)*vg.Inch, outputFile)
}

func plotLineCharts(m plotMeta, stats dbStats, outputFile string) error {
	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = m.Name()
	p.X.Label.Text = "time"
	p.Y.Label.Text = m.statType.String()

	vals := make([]interface{}, 0, len(stats)*2)
	for i := 0; i < len(stats); i++ {
		lineValue := stats[i].LineValue(m.op, m.statType)
		vals = append(vals, lineValue.name, lineValue.points)
	}

	err = plotutil.AddLinePoints(p, vals...)
	if err != nil {
		return err
	}

	return p.Save(vg.Length(plotXLength)*vg.Inch, vg.Length(plotYLength)*vg.Inch, outputFile)
}

func plotCharts(workload string, op string, stats dbStats) error {
	m := plotMeta{
		workload: workload,
		op:       op,
	}

	statTypes := []StatType{
		OPS, P99,
	}

	for _, tp := range statTypes {
		m.statType = tp
		outputFile := m.OutputPath("prog")
		if err := plotLineCharts(m, stats, outputFile); err != nil {
			return err
		}

		outputFile = m.OutputPath("total")
		if err := plotBarCharts(m, stats, outputFile); err != nil {
			return err
		}
	}

	return nil
}
