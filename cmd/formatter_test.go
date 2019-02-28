package cmd_test

import (
	"bytes"
	"testing"

	"github.com/aswinkarthik/csvdiff/cmd"
	"github.com/aswinkarthik/csvdiff/pkg/digest"

	"github.com/stretchr/testify/assert"
)

func TestJSONFormat(t *testing.T) {
	diff := digest.Differences{
		Additions:     []digest.Addition{[]string{"additions"}},
		Modifications: []digest.Modification{digest.Modification{Current: []string{"modification"}}},
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
	var stderr bytes.Buffer

	formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Config{Format: "json"})

	err := formatter.Format(diff)
	assert.NoError(t, err)
	assert.Equal(t, expected, stdout.String())
}

func TestRowMarkFormatter(t *testing.T) {
	diff := digest.Differences{
		Additions:     []digest.Addition{[]string{"additions"}},
		Modifications: []digest.Modification{digest.Modification{Current: []string{"modification"}}},
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

	formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Config{Format: "rowmark"})

	err := formatter.Format(diff)

	assert.NoError(t, err)
	assert.Equal(t, expectedStdout, stdout.String())
	assert.Equal(t, expectedStderr, stderr.String())
}

func TestLineDiff(t *testing.T) {
	diff := digest.Differences{
		Additions:     []digest.Addition{[]string{"additions"}},
		Modifications: []digest.Modification{digest.Modification{Original: []string{"original"}, Current: []string{"modification"}}},
	}
	expectedStdout := `+ additions
- original
+ modification
`
	expectedStderr := `# Additions (1)
# Modifications (1)
`

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Config{Format: "diff"})

	err := formatter.Format(diff)

	assert.NoError(t, err)
	assert.Equal(t, expectedStdout, stdout.String())
	assert.Equal(t, expectedStderr, stderr.String())
}

func TestWordDiff(t *testing.T) {
	diff := digest.Differences{
		Additions:     []digest.Addition{[]string{"additions"}},
		Modifications: []digest.Modification{digest.Modification{Original: []string{"original"}, Current: []string{"modification"}}},
	}
	expectedStdout := `{+additions+}
[-original-]{+modification+}
`
	expectedStderr := `# Additions (1)
# Modifications (1)
`

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Config{Format: "word-diff"})

	err := formatter.Format(diff)

	assert.NoError(t, err)
	assert.Equal(t, expectedStdout, stdout.String())
	assert.Equal(t, expectedStderr, stderr.String())
}

func TestColorWords(t *testing.T) {
	diff := digest.Differences{
		Additions:     []digest.Addition{[]string{"additions"}},
		Modifications: []digest.Modification{digest.Modification{Original: []string{"original"}, Current: []string{"modification"}}},
	}
	expectedStdout := `additions
originalmodification
`
	expectedStderr := `# Additions (1)
# Modifications (1)
`

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Config{Format: "color-words"})

	err := formatter.Format(diff)

	assert.NoError(t, err)
	assert.Equal(t, expectedStdout, stdout.String())
	assert.Equal(t, expectedStderr, stderr.String())
}
func TestWrongFormatter(t *testing.T) {
	diff := digest.Differences{}
	formatter := cmd.NewFormatter(nil, nil, cmd.Config{Format: "random-str"})

	err := formatter.Format(diff)

	assert.Error(t, err)
}
