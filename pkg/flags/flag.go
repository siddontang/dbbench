package flags

import "strings"

// ArrayFlags parses multi same flags into an arry.
type ArrayFlags []string

// String implemetns flag.Value interface.
func (i *ArrayFlags) String() string {
	return strings.Join(*i, ",")
}

// Set implemetns flag.Value interface.
func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// SetFlags parses multi same flags into a set.
type SetFlags map[string]struct{}

// String implemetns flag.Value interface.
func (i *SetFlags) String() string {
	names := make([]string, 0, len(*i))
	for name, _ := range *i {
		names = append(names, name)
	}
	return strings.Join(names, ",")
}

// Set implemetns flag.Value interface.
func (i *SetFlags) Set(value string) error {
	(*i)[value] = struct{}{}
	return nil
}
