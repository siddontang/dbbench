package stats

// Reporter defines how to report your own benchmark tool result.
type Reporter interface {
	// ParseName parses the db and workload from the result filename.
	// The filename format mostly is db_workload.log
	ParseName(fileName string) (db string, workload string)
	// NewDBStat parses the benchmark result file and saves the statistics into DBStat
	NewDBStat(name string, db string, workload string, filePath string) (*DBStat, error)
	// StatTypes returns the statistics we want to know
	StatTypes() []StatType
}

var reporters map[string]Reporter

// RegisterReporter registers the reporter with the unique name
func RegisterReporter(name string, r Reporter) {
	reporters[name] = r
}

// GetReporter gets the Reporter with the name
func GetReporter(name string) Reporter {
	return reporters[name]
}

func init() {
	reporters = make(map[string]Reporter)
}
