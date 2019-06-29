package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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

type StatType int

const (
	OPS StatType = iota
	P99
)

func (s StatType) String() string {
	switch s {
	case OPS:
		return "ops"
	case P99:
		return "p99(us)"
	default:
		return ""
	}
}

type stat struct {
	ops float64
	p99 float64
}

func statFieldFunc(c rune) bool {
	return c == ':' || c == ','
}

func newStat(line string) (*stat, error) {
	kvs := strings.FieldsFunc(line, statFieldFunc)
	s := stat{}
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
			s.ops = v
		case "99th(us)":
			s.p99 = v
		default:
		}
	}
	return &s, nil
}

func (s *stat) Value(tp StatType) float64 {
	switch tp {
	case OPS:
		return s.ops
	case P99:
		return s.p99
	default:
		perr(fmt.Errorf("unsupported stat type"))
		return 0.0
	}
}

type dbStat struct {
	name     string
	db       string
	workload string
	summary  map[string]*stat
	progress map[string][]*stat
}

func (s *dbStat) Operations() []string {
	names := make([]string, 0, len(s.summary))
	for name, _ := range s.summary {
		names = append(names, name)
	}
	return names
}

func (s *dbStat) parse(pathName string) error {
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

		stat, err := newStat(strings.TrimSpace(seps[1]))
		if err != nil {
			return err
		}

		if handleSummary {
			// handle summary logs
			s.summary[op] = stat
		} else {
			s.progress[op] = append(s.progress[op], stat)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func newDBStat(db string, workload string, pathName string) (*dbStat, error) {
	s := new(dbStat)
	s.summary = make(map[string]*stat, 1)
	s.progress = make(map[string][]*stat, 1)

	// We assume we put all logs in one unique directory in each benchmark.
	// E.g, we can use Git commit as the parent directory for benchmarking special version,
	// use datetime for benchmarking different databases.
	s.name = fmt.Sprintf("%s-%s", db, path.Base(filepath.Dir(pathName)))
	s.db = db
	s.workload = workload

	if err := s.parse(pathName); err != nil {
		return nil, err
	}

	return s, nil
}

type dbStats []*dbStat

func (a dbStats) Len() int      { return len(a) }
func (a dbStats) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a dbStats) Less(i, j int) bool {
	if a[i].db < a[j].db {
		return true
	} else if a[i].name < a[j].name {
		return true
	}
	return false
}
