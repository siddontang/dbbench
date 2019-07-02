// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/siddontang/dbbench/pkg/flags"
	"github.com/siddontang/dbbench/pkg/plot"
	"github.com/siddontang/dbbench/pkg/stats"
)

var (
	logPaths        flags.ArrayFlags
	filterDBs       flags.SetFlags
	filterWorkloads flags.SetFlags
	outputDir       string
	onlyDBName      bool
	plotXLength     int
	plotYLength     int
)

func init() {
	flag.Var(&logPaths, "p", "Log files or directories")
	flag.StringVar(&outputDir, "o", "./output", "Output directory")
	flag.IntVar(&plotXLength, "x", 8, "X axis Inch length of Output chart")
	flag.IntVar(&plotYLength, "y", 4, "Y axis Inch length of Output chart")
	filterDBs = make(flags.SetFlags, 0)
	filterWorkloads = make(flags.SetFlags, 0)
	flag.Var(&filterDBs, "d", "Filter database")
	flag.Var(&filterWorkloads, "w", "Filter workload")
	flag.BoolVar(&onlyDBName, "i", false, "Use only db name for identification in the chart")
}

func perr(err error) {
	if err == nil {
		return
	}

	fmt.Printf("meet err: %v\n", err)
	debug.PrintStack()
	os.Exit(1)
}

// We must ensure that the base filename of the sysbench log must be the format of db_workload.log
// Now the common workload may be oltp_read_write, oltp_update_non_index, etc.
func parseName(pathName string) (db string, workload string) {
	fileName := path.Base(pathName)
	seps := strings.Split(fileName, "_")

	db = seps[0]
	workload = strings.TrimSuffix(fileName[len(db)+1:], ".log")
	return db, workload
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
	flag.Parse()
	workloads := make(map[string]stats.DBStats)

	for _, logPath := range logPaths {
		err := filepath.Walk(logPath, func(path string, f os.FileInfo, err error) error {
			perr(err)

			if f.IsDir() {
				return nil
			}

			path, err = filepath.Abs(path)
			perr(err)

			db, workload := parseName(path)
			if db == "" || workload == "" {
				// invalid format
				return nil
			} else if !isFiltered(db, workload) {
				// we don't care these db and workload
				return nil
			}

			s, err := newDBStat(db, workload, path)
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
		StatTypes: []stats.StatType{
			stats.TPS,
			stats.QPS,
			stats.P95,
		},
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
