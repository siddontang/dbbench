package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/siddontang/dbbench/pkg/stats"
)

var (
	operations = map[string]struct{}{
		"INSERT":            struct{}{},
		"READ":              struct{}{},
		"UPDATE":            struct{}{},
		"SCAN":              struct{}{},
		"READ_MODIFY_WRITE": struct{}{},
		"DELETE":            struct{}{},
	}
)

func recordFieldFunc(c rune) bool {
	return c == ':' || c == ','
}

func newRecord(line string) (*stats.Record, error) {
	kvs := strings.FieldsFunc(line, recordFieldFunc)
	s := stats.Record{}
	if len(kvs)%2 != 0 {
		println(line)
	}
	for i := 0; i < len(kvs); i += 2 {
		v, err := strconv.ParseFloat(strings.TrimSpace(kvs[i+1]), 64)
		if err != nil {
			return nil, err
		}
		switch strings.TrimSpace(kvs[i]) {
		case "OPS":
			s.OPS = v
		case "99th(us)":
			s.P99 = v / 1000.0
		default:
		}
	}
	return &s, nil
}

func parse(workload string, pathName string, s *stats.DBStat) error {
	file, err := os.Open(pathName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	handleSummary := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Run finished") {
			handleSummary = true
			continue
		}

		seps := strings.Split(line, "-")
		op := strings.TrimSpace(seps[0])
		if _, ok := operations[op]; !ok {
			continue
		}

		if workload == "laod" {
			op = ""
		}

		r, err := newRecord(strings.TrimSpace(seps[1]))
		if err != nil {
			return err
		}

		if handleSummary {
			// handle summary logs
			s.Summary[op] = r
		} else {
			s.Progress[op] = append(s.Progress[op], r)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func newDBStat(db string, workload string, pathName string) (*stats.DBStat, error) {
	name := pathName
	if onlyDBName {
		name = ""
	}
	s := stats.NewDBStat(db, workload, name)

	if err := parse(workload, pathName, s); err != nil {
		return nil, err
	}

	return s, nil
}
