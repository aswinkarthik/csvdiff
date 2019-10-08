package cmd_test

import (
	"bytes"
	"testing"

	"github.com/aswinkarthik/csvdiff/cmd"
	"github.com/aswinkarthik/csvdiff/pkg/digest"

	"github.com/stretchr/testify/assert"
)

func TestLegacyJSONFormat(t *testing.T) {
	diff := digest.Differences{
		Additions:     []digest.Addition{[]string{"additions"}},
		Modifications: []digest.Modification{{Current: []string{"modification"}}},
		Deletions:     []digest.Deletion{[]string{"deletions"}},
	}
	expected := `{
  "Additions": [
    "additions"
  ],
  "Modifications": [
    "modification"
  ],
  "Deletions": [
    "deletions"
  ]
}`

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Context{Format: "legacy-json"})

	err := formatter.Format(diff)
	assert.NoError(t, err)
	assert.Equal(t, expected, stdout.String())
}

func TestJSONFormat(t *testing.T) {
	diff := digest.Differences{
		Additions:     []digest.Addition{[]string{"additions"}},
		Modifications: []digest.Modification{{Original: []string{"original"}, Current: []string{"modification"}}},
		Deletions:     []digest.Deletion{[]string{"deletions"}},
	}
	expected := `{
  "Additions": [
    "additions"
  ],
  "Modifications": [
    {
      "Original": "original",
      "Current": "modification"
    }
  ],
  "Deletions": [
    "deletions"
  ]
}`

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Context{Format: "json"})

	err := formatter.Format(diff)
	assert.NoError(t, err)
	assert.Equal(t, expected, stdout.String())
}
func TestRowMarkFormatter(t *testing.T) {
	diff := digest.Differences{
		Additions:     []digest.Addition{[]string{"additions"}},
		Modifications: []digest.Modification{{Current: []string{"modification"}}},
		Deletions:     []digest.Deletion{[]string{"deletions"}},
	}
	expectedStdout := `additions,ADDED
modification,MODIFIED
deletions,DELETED
`
	expectedStderr := `Additions 1
Modifications 1
Deletions 1
Rows:
`

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Context{Format: "rowmark"})

	err := formatter.Format(diff)

	assert.NoError(t, err)
	assert.Equal(t, expectedStdout, stdout.String())
	assert.Equal(t, expectedStderr, stderr.String())
}

func TestLineDiff(t *testing.T) {
	diff := digest.Differences{
		Additions: []digest.Addition{[]string{"additions"}},
		Modifications: []digest.Modification{
			{
				Original: []string{"original", "comma,separated,value"},
				Current:  []string{"modification", "comma,separated,value-2"},
			},
		},
		Deletions: []digest.Deletion{{"deletion", "this-row-was-deleted"}},
	}
	expectedStdout := `+ additions
- original,"comma,separated,value"
+ modification,"comma,separated,value-2"
- deletion,this-row-was-deleted
`
	expectedStderr := `# Additions (1)
# Modifications (1)
# Deletions (1)
`

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Context{Format: "diff"})

	err := formatter.Format(diff)

	assert.NoError(t, err)
	assert.Equal(t, expectedStdout, stdout.String())
	assert.Equal(t, expectedStderr, stderr.String())
}

func TestWordDiff(t *testing.T) {
	t.Run("should cover single column happy path", func(t *testing.T) {
		diff := digest.Differences{
			Additions:     []digest.Addition{[]string{"additions"}},
			Modifications: []digest.Modification{{Original: []string{"original"}, Current: []string{"modification"}}},
			Deletions:     []digest.Deletion{{"deletions"}},
		}
		expectedStdout := `{+additions+}
[-original-]{+modification+}
[-deletions-]
`
		expectedStderr := `# Additions (1)
# Modifications (1)
# Deletions (1)
`

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Context{Format: "word-diff"})

		err := formatter.Format(diff)

		assert.NoError(t, err)
		assert.Equal(t, expectedStdout, stdout.String())
		assert.Equal(t, expectedStderr, stderr.String())
	})

	t.Run("should ouput only selective columns", func(t *testing.T) {
		diff := digest.Differences{
			Additions: []digest.Addition{[]string{"additions", "ignored-column"}},
			Modifications: []digest.Modification{
				{Original: []string{"original", "ignored-column"}, Current: []string{"modification", "ignored-column"}},
			},
			Deletions: []digest.Deletion{{"deletions", "ignored-column"}},
		}
		expectedStdout := `{+additions+}
[-original-]{+modification+}
[-deletions-]
`
		expectedStderr := `# Additions (1)
# Modifications (1)
# Deletions (1)
`

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Context{
			Format:                 "word-diff",
			IncludeColumnPositions: digest.Positions{0},
		})

		err := formatter.Format(diff)

		assert.NoError(t, err)
		assert.Equal(t, expectedStdout, stdout.String())
		assert.Equal(t, expectedStderr, stderr.String())

	})
}

func TestColorWords(t *testing.T) {
	diff := digest.Differences{
		Additions:     []digest.Addition{[]string{"additions"}},
		Modifications: []digest.Modification{{Original: []string{"original"}, Current: []string{"modification"}}},
		Deletions:     []digest.Deletion{{"deletions"}},
	}
	expectedStdout := `additions
originalmodification
deletions
`
	expectedStderr := `# Additions (1)
# Modifications (1)
# Deletions (1)
`

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	formatter := cmd.NewFormatter(&stdout, &stderr, cmd.Context{Format: "color-words"})

	err := formatter.Format(diff)

	assert.NoError(t, err)
	assert.Equal(t, expectedStdout, stdout.String())
	assert.Equal(t, expectedStderr, stderr.String())
}

func TestWrongFormatter(t *testing.T) {
	diff := digest.Differences{}
	formatter := cmd.NewFormatter(nil, nil, cmd.Context{Format: "random-str"})

	err := formatter.Format(diff)

	assert.Error(t, err)
}
