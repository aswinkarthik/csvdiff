package cmd_test

import (
	"bytes"
	"testing"

	"github.com/aswinkarthik93/csvdiff/cmd"
	"github.com/aswinkarthik93/csvdiff/pkg/digest"

	"github.com/stretchr/testify/assert"
)

func TestJSONFormat(t *testing.T) {
	var formatter cmd.Formatter
	diff := digest.Difference{
		Additions:     []string{"additions"},
		Modifications: []string{"modification"},
	}
	expected := `{
  "Additions": [
    "additions"
  ],
  "Modifications": [
    "modification"
  ]
}`

	var buffer bytes.Buffer

	formatter = &cmd.JSONFormatter{}

	formatter.Format(diff, &buffer)
	assert.Equal(t, expected, buffer.String())
}

func TestRowMarkFormatter(t *testing.T) {
	var formatter cmd.Formatter
	diff := digest.Difference{
		Additions:     []string{"additions"},
		Modifications: []string{"modification"},
	}
	expected := `Additions 1
Modifications 1
Rows:
additions,ADDED
modification,MODIFIED
`

	var buffer bytes.Buffer

	formatter = &cmd.RowMarkFormatter{}

	formatter.Format(diff, &buffer)
	assert.Equal(t, expected, buffer.String())
}
