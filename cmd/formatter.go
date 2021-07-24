package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/fatih/color"
	"io"
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
	ctx    Context
}

// NewFormatter can be used to create a new formatter
func NewFormatter(stdout, stderr io.Writer, ctx Context) *Formatter {
	if ctx.separator == rune(0) {
		ctx.separator = ','
	}
	return &Formatter{stdout: stdout, stderr: stderr, ctx: ctx}
}

// Format can be used to format the differences based on ctx
// to appropriate writers
func (f *Formatter) Format(diff digest.Differences) error {
	switch f.ctx.format {
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
		Deletions     []string
	}

	includes := f.ctx.GetIncludeColumnPositions()

	additions := make([]string, 0, len(diff.Additions))
	for _, addition := range diff.Additions {
		additions = append(additions, includes.String(addition, f.ctx.separator))
	}

	modifications := make([]string, 0, len(diff.Modifications))
	for _, modification := range diff.Modifications {
		modifications = append(modifications, includes.String(modification.Current, f.ctx.separator))
	}

	deletions := make([]string, 0, len(diff.Deletions))
	for _, deletion := range diff.Deletions {
		deletions = append(deletions, includes.String(deletion, f.ctx.separator))
	}

	jsonDiff := jsonDifference{Additions: additions, Modifications: modifications, Deletions: deletions}
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
	includes := f.ctx.GetIncludeColumnPositions()

	additions := make([]string, 0, len(diff.Additions))
	for _, addition := range diff.Additions {
		additions = append(additions, includes.String(addition, f.ctx.separator))
	}

	deletions := make([]string, 0, len(diff.Deletions))
	for _, deletion := range diff.Deletions {
		deletions = append(deletions, includes.String(deletion, f.ctx.separator))
	}

	type modification struct {
		Original string
		Current  string
	}

	type jsonDifference struct {
		Additions     []string
		Modifications []modification
		Deletions     []string
	}

	modifications := make([]modification, 0, len(diff.Modifications))
	for _, mods := range diff.Modifications {
		m := modification{Original: includes.String(mods.Original, f.ctx.separator), Current: includes.String(mods.Current, f.ctx.separator)}
		modifications = append(modifications, m)
	}

	data, err := json.MarshalIndent(jsonDifference{Additions: additions, Modifications: modifications, Deletions: deletions}, "", "  ")

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
	_, _ = fmt.Fprintf(f.stderr, "Additions %d\n", len(diff.Additions))
	_, _ = fmt.Fprintf(f.stderr, "Modifications %d\n", len(diff.Modifications))
	_, _ = fmt.Fprintf(f.stderr, "Deletions %d\n", len(diff.Deletions))
	_, _ = fmt.Fprintf(f.stderr, "Rows:\n")

	includes := f.ctx.GetIncludeColumnPositions()

	additions := make([]string, 0, len(diff.Additions))
	for _, addition := range diff.Additions {
		additions = append(additions, includes.String(addition, f.ctx.separator))
	}

	modifications := make([]string, 0, len(diff.Modifications))
	for _, modification := range diff.Modifications {
		modifications = append(modifications, includes.String(modification.Current, f.ctx.separator))
	}

	deletions := make([]string, 0, len(diff.Deletions))
	for _, deletion := range diff.Deletions {
		deletions = append(deletions, includes.String(deletion, f.ctx.separator))
	}

	for _, added := range additions {
		_, _ = fmt.Fprintf(f.stdout, "%s,%s\n", added, "ADDED")
	}

	for _, modified := range modifications {
		_, _ = fmt.Fprintf(f.stdout, "%s,%s\n", modified, "MODIFIED")
	}

	for _, deleted := range deletions {
		_, _ = fmt.Fprintf(f.stdout, "%s,%s\n", deleted, "DELETED")
	}

	return nil
}

// lineDiff is git-style line diff
func (f *Formatter) lineDiff(diff digest.Differences) error {
	includes := f.ctx.GetIncludeColumnPositions()

	blue := color.New(color.FgBlue).FprintfFunc()
	red := color.New(color.FgRed).FprintfFunc()
	green := color.New(color.FgGreen).FprintfFunc()

	blue(f.stderr, "# Additions (%d)\n", len(diff.Additions))
	for _, addition := range diff.Additions {
		green(f.stdout, "+ %s\n", includes.String(addition, f.ctx.separator))
	}
	blue(f.stderr, "# Modifications (%d)\n", len(diff.Modifications))
	for _, modification := range diff.Modifications {
		red(f.stdout, "- %s\n", includes.String(modification.Original, f.ctx.separator))
		green(f.stdout, "+ %s\n", includes.String(modification.Current, f.ctx.separator))
	}
	blue(f.stderr, "# Deletions (%d)\n", len(diff.Deletions))
	for _, deletion := range diff.Deletions {
		red(f.stdout, "- %s\n", includes.String(deletion, f.ctx.separator))
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
	includes := f.ctx.GetIncludeColumnPositions()
	if len(includes) <= 0 {
		includes = f.ctx.GetValueColumns()
	}
	blue := color.New(color.FgBlue).SprintfFunc()
	red := color.New(color.FgRed).SprintfFunc()
	green := color.New(color.FgGreen).SprintfFunc()

	_, _ = fmt.Fprintln(f.stderr, blue("# Additions (%d)", len(diff.Additions)))
	for _, addition := range diff.Additions {
		_, _ = fmt.Fprintln(f.stdout, green(additionFormat, includes.String(addition, f.ctx.separator)))
	}

	_, _ = fmt.Fprintln(f.stderr, blue("# Modifications (%d)", len(diff.Modifications)))
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
		_, _ = fmt.Fprintln(f.stdout, includes.String(result, f.ctx.separator))
	}

	_, _ = fmt.Fprintln(f.stderr, blue("# Deletions (%d)", len(diff.Deletions)))
	for _, deletion := range diff.Deletions {
		_, _ = fmt.Fprintln(f.stdout, red(deletionFormat, includes.String(deletion, f.ctx.separator)))
	}

	return nil

}
