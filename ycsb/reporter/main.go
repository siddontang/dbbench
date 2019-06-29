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
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	logPaths  arrayFlags
	outputDir string
)

func init() {
	flag.Var(&logPaths, "p", "Log files or directories")
	flag.StringVar(&outputDir, "o", "./output", "Output directory")
	flag.IntVar(&plotXLength, "x", 8, "X axis Inch length of Output chart")
	flag.IntVar(&plotYLength, "y", 4, "Y axis Inch length of Output chart")
}

func perr(err error) {
	if err == nil {
		return
	}

	fmt.Printf("meet err: %v\n", err)
	debug.PrintStack()
	os.Exit(1)
}

// We must ensure that the base filename of the YCSB log must be the format of db_workload.log
// Now the common workload is load, workloada, workloadb, ... workloadf, if you want to use you own workload,
// please use a unique workload name.
func parseName(pathName string) (db string, workload string) {
	// check db and workload from file name, the name format is:
	// 	1. db_load.log
	// 	2. db_workloadx.log

	fileName := path.Base(pathName)
	seps := strings.Split(fileName, "_")

	if len(seps) != 2 {
		return "", ""
	}

	db = seps[0]
	workload = strings.TrimSuffix(seps[1], ".log")
	return db, workload
}

func main() {
	flag.Parse()
	workloads := make(map[string]dbStats)

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
			}

			s, err := newDBStat(db, workload, path)
			perr(err)

			workloads[workload] = append(workloads[workload], s)

			return nil
		})
		perr(err)
	}

	for workload, stats := range workloads {
		sort.Sort(stats)

		os.MkdirAll(path.Join(outputDir, workload), 0755)

		operations := stats[0].Operations()

		for _, op := range operations {
			err := plotCharts(workload, op, stats)
			perr(err)
		}
	}
}
