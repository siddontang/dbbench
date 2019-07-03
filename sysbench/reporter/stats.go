package reporter

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"fmt"

	"github.com/siddontang/dbbench/pkg/stats"
)

type record struct {
	Second          int32
	Threads         int32
	TPS             float64
	QPS             float64
	ReadQPS         float64
	WriteQPS        float64
	OtherQPS        float64
	Latency         float64
	LatencyLimit    float64
	ErrorPerSec     float64
	ReconnectPerSec float64
}

const (
	recordParseFormat = "[ %ds ] thds: %d tps: %f qps: %f (r/w/o: %f/%f/%f) lat (ms,%f%%): %f err/s: %f reconn/s: %f"
)

func parseRecord(str string) (record, error) {
	r := record{}

	_, err := fmt.Sscanf(
		str,
		recordParseFormat,
		&r.Second,
		&r.Threads,
		&r.TPS,
		&r.QPS,
		&r.ReadQPS,
		&r.WriteQPS,
		&r.OtherQPS,
		&r.LatencyLimit,
		&r.Latency,
		&r.ErrorPerSec,
		&r.ReconnectPerSec)
	if err != nil {
		return record{}, err
	}
	return r, nil
}

func newRecord(line string) (*stats.Record, error) {
	r, err := parseRecord(line)
	if err != nil {
		return nil, err
	}

	s := stats.Record{
		TPS: r.TPS,
		QPS: r.QPS,
		// Assume we use P95 now
		P95: r.Latency,
	}

	return &s, nil
}

func parseSummaryOPS(line string, r *stats.Record) {
	line = strings.TrimSpace(line)
	tp := stats.None
	if strings.HasPrefix(line, "transactions:") {
		tp = stats.TPS
	} else if strings.HasPrefix(line, "queries:") {
		tp = stats.QPS
	}

	if tp != stats.TPS && tp != stats.QPS {
		return
	}

	seps := strings.Split(line, "(")
	seps = strings.Split(seps[1], "per sec")
	v, _ := strconv.ParseFloat(strings.TrimSpace(seps[0]), 64)
	if tp == stats.TPS {
		r.TPS = v
		return
	}

	r.QPS = v
}

func parseSummaryLatency(line string, r *stats.Record) {
	line = strings.TrimSpace(line)
	// We only care P95 now
	if !strings.HasPrefix(line, "95th percentile:") {
		return
	}

	seps := strings.Split(line, ":")
	v, _ := strconv.ParseFloat(strings.TrimSpace(seps[1]), 64)
	r.P95 = v
}

func parse(workload string, pathName string, s *stats.DBStat) error {
	file, err := os.Open(pathName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	handleSummary := false

	summaryRecord := &stats.Record{}

	for scanner.Scan() {
		line := scanner.Text()
		if handleSummary {
			parseSummaryOPS(line, summaryRecord)
			parseSummaryLatency(line, summaryRecord)
			continue
		}

		if strings.HasPrefix(line, "SQL statistics") {
			handleSummary = true
			continue
		}

		if r, err := newRecord(strings.TrimSpace(line)); err == nil {
			s.Progress[""] = append(s.Progress[""], r)
		}
	}

	s.Summary[""] = summaryRecord

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

type reporter struct {
}

func (r reporter) NewDBStat(name string, db string, workload string, pathName string) (*stats.DBStat, error) {
	s := stats.NewDBStat(name, db, workload, name)

	if err := parse(workload, pathName, s); err != nil {
		return nil, err
	}

	return s, nil
}

func (r reporter) ParseName(fileName string) (db string, workload string) {
	// We must ensure that the base filename of the sysbench log must be the format of db_workload.log
	// Now the common workload may be oltp_read_write, oltp_update_non_index, etc.	fileName := path.Base(pathName)
	seps := strings.Split(fileName, "_")

	db = seps[0]
	workload = strings.TrimSuffix(fileName[len(db)+1:], ".log")
	return db, workload
}

func (r reporter) StatTypes() []stats.StatType {
	return []stats.StatType{
		stats.TPS,
		stats.QPS,
		stats.P95,
	}
}

func init() {
	stats.RegisterReporter("sysbench", reporter{})
}
