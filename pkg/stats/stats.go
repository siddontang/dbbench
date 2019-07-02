package stats

import (
	"fmt"
	"path"
	"path/filepath"
)

// StatType is the type of different statistics
type StatType int

// Different statistics for different benchmark tools.
const (
	None StatType = iota
	OPS
	P99
	P95
	TPS
	QPS
)

// String returns the name of statistics
func (s StatType) String() string {
	switch s {
	case OPS:
		return "ops"
	case P99:
		return "p99(ms)"
	case P95:
		return "p95(ms)"
	case TPS:
		return "tps"
	case QPS:
		return "qps"
	default:
		return ""
	}
}

// Value saves the statistics value
type Record struct {
	OPS float64
	TPS float64
	QPS float64
	P99 float64
	P95 float64
}

// Value returns the value of different statistics
func (r *Record) Value(tp StatType) float64 {
	switch tp {
	case OPS:
		return r.OPS
	case P99:
		return r.P99
	case P95:
		return r.P95
	case TPS:
		return r.TPS
	case QPS:
		return r.QPS
	default:
		return 0.0
	}
}

// DBStat holds all statistics in one benchmark
type DBStat struct {
	// Name is the unqiue name used in plotting later.
	Name string
	// DB is the database name
	DB string
	// Workload is the benchmark workload name
	Workload string
	// Summary holds the final output summary record
	// The key of the map is the operation in the benchmark.
	// E.g, in go-ycsb, the operation may be INSERT, READ
	Summary map[string]*Record

	// Progress holds the in progess record in benchmarking
	Progress map[string][]*Record
}

// Operations returns all the operations in the test
func (s *DBStat) Operations() []string {
	names := make([]string, 0, len(s.Summary))
	for name, _ := range s.Summary {
		names = append(names, name)
	}
	return names
}

// NewDBStat creates a DBStat.
// We assume we put all logs in one unique directory in each benchmark.
// E.g, we can use Git commit as the parent directory for benchmarking special version,
// use datetime for benchmarking different databases.
// If pathName is empty, we will use db as the name of DBStat.
func NewDBStat(db string, workload string, pathName string) *DBStat {
	s := new(DBStat)
	s.Summary = make(map[string]*Record, 1)
	s.Progress = make(map[string][]*Record, 1)

	// We assume we put all logs in one unique directory in each benchmark.
	// E.g, we can use Git commit as the parent directory for benchmarking special version,
	// use datetime for benchmarking different databases.
	if pathName != "" {
		s.Name = fmt.Sprintf("%s-%s", db, path.Base(filepath.Dir(pathName)))
	} else {
		s.Name = db
	}
	s.DB = db
	s.Workload = workload

	return s
}

// DBStats is the array of DBStat.
type DBStats []*DBStat

func (a DBStats) Len() int      { return len(a) }
func (a DBStats) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a DBStats) Less(i, j int) bool {
	if a[i].DB < a[j].DB {
		return true
	} else if a[i].Name < a[j].Name {
		return true
	}
	return false
}
