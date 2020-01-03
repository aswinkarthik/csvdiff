package digest

import "io"

// Config represents configurations that can be passed
// to create a Digest.
//
// Key: The primary key positions
// Value: The Value positions that needs to be compared for diff
// Include: Include these positions in output. It is Value positions by default.
type Config struct {
	Key         Positions
	Value       Positions
	Include     Positions
	Reader      io.Reader
	Separator   rune
	LazyQuotes  bool
}

// NewConfig creates an instance of Config struct.
func NewConfig(
	r io.Reader,
	primaryKey Positions,
	valueColumns Positions,
	includeColumns Positions,
	separator rune,
	lazyQuotes bool,
) *Config {
	if len(includeColumns) == 0 {
		includeColumns = valueColumns
	}

	return &Config{
		Reader:     r,
		Key:        primaryKey,
		Value:      valueColumns,
		Include:    includeColumns,
		Separator:  separator,
		LazyQuotes: lazyQuotes,
	}
}
