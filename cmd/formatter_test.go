package cmd_test

import (
	"bytes"
	"testing"

	"github.com/aswinkarthik/csvdiff/cmd"
	"github.com/aswinkarthik/csvdiff/pkg/digest"

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

	var stdout bytes.Buffer

	formatter = &cmd.JSONFormatter{Stdout: &stdout}

	err := formatter.Format(diff)
	assert.NoError(t, err)
	assert.Equal(t, expected, stdout.String())
}

func TestRowMarkFormatter(t *testing.T) {
	var formatter cmd.Formatter
	diff := digest.Difference{
		Additions:     []string{"additions"},
		Modifications: []string{"modification"},
	}
	expectedStdout := `additions,ADDED
modification,MODIFIED
`
	expectedStderr := `Additions 1
Modifications 1
Rows:
`

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	formatter = &cmd.RowMarkFormatter{Stdout: &stdout, Stderr: &stderr}

	err := formatter.Format(diff)

	assert.NoError(t, err)
	assert.Equal(t, expectedStdout, stdout.String())
	assert.Equal(t, expectedStderr, stderr.String())
}
