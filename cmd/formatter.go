package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/fatih/color"
)

const (
	rowmark          = "rowmark"
	jsonFormat       = "json"
	legacyJSONFormat = "legacy-json"
	lineDiff         = "diff"
	wordDiff         = "word-diff"
	colorWords       = "color-words"
)

var allFormats = []string{rowmark, jsonFormat, legacyJSONFormat, lineDiff, wordDiff, colorWords}

// Formatter can print the differences to stdout
// and accompanying metadata to stderr
type Formatter struct {
	stdout io.Writer
	stderr io.Writer
	config Config
}

// NewFormatter can be used to create a new formatter
func NewFormatter(stdout, stderr io.Writer, config Config) *Formatter {
	return &Formatter{stdout: stdout, stderr: stderr, config: config}
}

// Format can be used to format the differences based on config
// to appropriate writers
func (f *Formatter) Format(diff digest.Differences) error {
	switch f.config.Format {
	case legacyJSONFormat:
		return f.legacyJSON(diff)
	case jsonFormat:
		return f.json(diff)
	case rowmark:
		return f.rowMark(diff)
	case lineDiff:
		return f.lineDiff(diff)
	case wordDiff:
		return f.wordDiff(diff)
	case colorWords:
		return f.colorWords(diff)
	default:
		return fmt.Errorf("formatter not found")
	}
}

// JSONFormatter formats diff to as a JSON Object
// { "Additions": [...], "Modifications": [...] }
func (f *Formatter) legacyJSON(diff digest.Differences) error {
	// jsonDifference is a struct to represent legacy JSON format
	type jsonDifference struct {
		Additions     []string
		Modifications []string
	}

	includes := config.GetIncludeColumnPositions()

	additions := make([]string, 0, len(diff.Additions))
	for _, addition := range diff.Additions {
		additions = append(additions, includes.MapToValue(addition))
	}

	modifications := make([]string, 0, len(diff.Modifications))
	for _, modification := range diff.Modifications {
		modifications = append(modifications, includes.MapToValue(modification.Current))
	}

	jsonDiff := jsonDifference{Additions: additions, Modifications: modifications}
	data, err := json.MarshalIndent(jsonDiff, "", "  ")

	if err != nil {
		return fmt.Errorf("error when serializing with JSON formatter: %v", err)
	}

	_, err = f.stdout.Write(data)

	if err != nil {
		return fmt.Errorf("error when writing to writer with JSON formatter: %v", err)
	}

	return nil
}

// JSONFormatter formats diff to as a JSON Object
// { "Additions": [...], "Modifications": [{ "Original": [...], "Current": [...]}]}
func (f *Formatter) json(diff digest.Differences) error {
	includes := config.GetIncludeColumnPositions()

	additions := make([]string, 0, len(diff.Additions))
	for _, addition := range diff.Additions {
		additions = append(additions, includes.MapToValue(addition))
	}

	type modification struct {
		Original string
		Current  string
	}

	type jsonDifference struct {
		Additions     []string
		Modifications []modification
	}

	modifications := make([]modification, 0, len(diff.Modifications))
	for _, mods := range diff.Modifications {
		m := modification{Original: includes.MapToValue(mods.Original), Current: includes.MapToValue(mods.Current)}
		modifications = append(modifications, m)
	}

	data, err := json.MarshalIndent(jsonDifference{Additions: additions, Modifications: modifications}, "", "  ")

	if err != nil {
		return fmt.Errorf("error when serializing with JSON formatter: %v", err)
	}

	_, err = f.stdout.Write(data)

	if err != nil {
		return fmt.Errorf("error when writing to writer with JSON formatter: %v", err)
	}

	return nil
}

// RowMarkFormatter formats diff by marking each row as
// ADDED/MODIFIED. It mutates the row and adds as a new column.
func (f *Formatter) rowMark(diff digest.Differences) error {
	fmt.Fprintf(f.stderr, "Additions %d\n", len(diff.Additions))
	fmt.Fprintf(f.stderr, "Modifications %d\n", len(diff.Modifications))
	fmt.Fprintf(f.stderr, "Rows:\n")

	includes := config.GetIncludeColumnPositions()

	additions := make([]string, 0, len(diff.Additions))
	for _, addition := range diff.Additions {
		additions = append(additions, includes.MapToValue(addition))
	}

	modifications := make([]string, 0, len(diff.Modifications))
	for _, modification := range diff.Modifications {
		modifications = append(modifications, includes.MapToValue(modification.Current))
	}

	for _, added := range additions {
		fmt.Fprintf(f.stdout, "%s,%s\n", added, "ADDED")
	}

	for _, modified := range modifications {
		fmt.Fprintf(f.stdout, "%s,%s\n", modified, "MODIFIED")
	}

	return nil
}

// lineDiff is git-style line diff
func (f *Formatter) lineDiff(diff digest.Differences) error {
	includes := config.GetIncludeColumnPositions()

	blue := color.New(color.FgBlue).FprintfFunc()
	red := color.New(color.FgRed).FprintfFunc()
	green := color.New(color.FgGreen).FprintfFunc()

	blue(f.stderr, "# Additions (%d)\n", len(diff.Additions))
	for _, addition := range diff.Additions {
		green(f.stdout, "+ %s\n", includes.MapToValue(addition))
	}
	blue(f.stderr, "# Modifications (%d)\n", len(diff.Modifications))
	for _, modification := range diff.Modifications {
		red(f.stdout, "- %s\n", includes.MapToValue(modification.Original))
		green(f.stdout, "+ %s\n", includes.MapToValue(modification.Current))
	}

	return nil
}

// wordDiff is git-style --word-diff
func (f *Formatter) wordDiff(diff digest.Differences) error {
	return f.wordLevelDiffs(diff, "[-%s-]", "{+%s+}")
}

// colorWords is git-style --color-words
func (f *Formatter) colorWords(diff digest.Differences) error {
	return f.wordLevelDiffs(diff, "%s", "%s")
}

func (f *Formatter) wordLevelDiffs(diff digest.Differences, deletionFormat, additionFormat string) error {
	includes := config.GetIncludeColumnPositions()
	if len(includes) <= 0 {
		includes = config.GetValueColumns()
	}
	blue := color.New(color.FgBlue).SprintfFunc()
	red := color.New(color.FgRed).SprintfFunc()
	green := color.New(color.FgGreen).SprintfFunc()

	fmt.Fprintln(f.stderr, blue("# Additions (%d)", len(diff.Additions)))
	for _, addition := range diff.Additions {
		fmt.Fprintln(f.stdout, green(additionFormat, includes.MapToValue(addition)))
	}

	fmt.Fprintln(f.stderr, blue("# Modifications (%d)", len(diff.Modifications)))
	for _, modification := range diff.Modifications {
		result := make([]string, 0, len(modification.Current))
		for i := 0; i < len(includes) || i < len(modification.Current); i++ {
			if modification.Original[i] != modification.Current[i] {
				removed := red(deletionFormat, modification.Original[i])
				added := green(additionFormat, modification.Current[i])
				result = append(result, fmt.Sprintf("%s%s", removed, added))
			} else {
				result = append(result, modification.Current[i])
			}
		}
		fmt.Fprintln(f.stdout, strings.Join(result, digest.Separator))
	}

	return nil

}
