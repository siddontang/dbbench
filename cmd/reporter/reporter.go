package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"sort"

	"github.com/siddontang/dbbench/pkg/flags"
	"github.com/siddontang/dbbench/pkg/plot"
	"github.com/siddontang/dbbench/pkg/stats"

	// Register different benchmark
	_ "github.com/siddontang/dbbench/sysbench/reporter"
	_ "github.com/siddontang/dbbench/ycsb/reporter"
)

var (
	logPaths        flags.ArrayFlags
	filterDBs       flags.SetFlags
	filterWorkloads flags.SetFlags
	outputDir       string
	onlyDBName      bool
	plotXLength     int
	plotYLength     int

	commandLine *flag.FlagSet
)

func init() {
	commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	commandLine.Var(&logPaths, "p", "Log files or directories")
	commandLine.StringVar(&outputDir, "o", "./output", "Output directory")
	commandLine.IntVar(&plotXLength, "x", 8, "X axis Inch length of Output chart")
	commandLine.IntVar(&plotYLength, "y", 4, "Y axis Inch length of Output chart")
	filterDBs = make(flags.SetFlags, 0)
	filterWorkloads = make(flags.SetFlags, 0)
	commandLine.Var(&filterDBs, "d", "Filter database")
	commandLine.Var(&filterWorkloads, "w", "Filter workload")
	commandLine.BoolVar(&onlyDBName, "i", false, "Use only db name for identification in the chart")

	commandLine.Usage = func() {
		fmt.Fprintf(commandLine.Output(), "Usage of %s benchmark_name:\n", os.Args[0])
		commandLine.PrintDefaults()
	}
}

func perr(err error) {
	if err == nil {
		return
	}

	fmt.Printf("meet err: %v\n", err)
	debug.PrintStack()
	os.Exit(1)
}

func isFiltered(db string, workload string) bool {
	if len(filterDBs) == 0 && len(filterWorkloads) == 0 {
		return true
	}
	_, ok1 := filterDBs[db]
	_, ok2 := filterWorkloads[workload]

	if len(filterDBs) > 0 && !ok1 {
		return false
	} else if len(filterWorkloads) > 0 && !ok2 {
		return false
	}

	return true
}

func main() {
	name := os.Args[1]
	r := stats.GetReporter(name)
	if r == nil {
		commandLine.Usage()
		os.Exit(1)
	}

	commandLine.Parse(os.Args[2:])

	outputDir = path.Join(outputDir, name)

	workloads := make(map[string]stats.DBStats)

	for _, logPath := range logPaths {
		err := filepath.Walk(logPath, func(pathName string, f os.FileInfo, err error) error {
			perr(err)

			if f.IsDir() {
				return nil
			}

			pathName, err = filepath.Abs(pathName)
			perr(err)

			db, workload := r.ParseName(path.Base(pathName))
			if db == "" || workload == "" {
				// invalid format
				return nil
			} else if !isFiltered(db, workload) {
				// we don't care these db and workload
				return nil
			}

			// We assume we put all logs in one unique directory in each benchmark.
			// E.g, we can use Git commit as the parent directory for benchmarking special version.
			name := db
			if !onlyDBName {
				name = fmt.Sprintf("%s-%s", db, path.Base(filepath.Dir(pathName)))
			}
			s, err := r.NewDBStat(name, db, workload, pathName)
			perr(err)

			workloads[workload] = append(workloads[workload], s)

			return nil
		})
		perr(err)
	}

	m := plot.Meta{
		OutputDir: outputDir,
		XLength:   plotXLength,
		YLength:   plotYLength,
		StatTypes: r.StatTypes(),
	}

	for workload, stats := range workloads {
		sort.Sort(stats)

		os.MkdirAll(path.Join(outputDir, workload), 0755)

		operations := stats[0].Operations()

		for _, op := range operations {
			m.Workload = workload
			m.OP = op
			err := plot.PlotCharts(m, stats)
			perr(err)
		}
	}
}
