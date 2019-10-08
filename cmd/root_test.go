package cmd

import (
	"bytes"
	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRunContext(t *testing.T) {
	t.Run("should find diff in happy path", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		{
			baseContent := []byte(`id,name,age,desc
0,tom,2,developer
2,ryan,20,qa
4,emin,40,pm

`)
			err := afero.WriteFile(fs, "/base.csv", baseContent, os.ModePerm)
			assert.NoError(t, err)
		}
		{
			deltaContent := []byte(`id,name,age,desc
0,tom,2,developer
1,caprio,3,developer
2,ryan,23,qa
`)
			err := afero.WriteFile(fs, "/delta.csv", deltaContent, os.ModePerm)
			assert.NoError(t, err)
		}

		ctx := Context{
			Format:                 "json",
			BaseFilename:           "/base.csv",
			DeltaFilename:          "/delta.csv",
			PrimaryKeyPositions:    digest.Positions{0},
			ValueColumnPositions:   digest.Positions{1, 2},
			IncludeColumnPositions: digest.Positions{0, 1, 2},
		}
		outStream := &bytes.Buffer{}
		errStream := &bytes.Buffer{}

		err := runContext(ctx, fs, outStream, errStream)
		expected := `{
  "Additions": [
    "1,caprio,3"
  ],
  "Modifications": [
    {
      "Original": "2,ryan,20",
      "Current": "2,ryan,23"
    }
  ],
  "Deletions": [
    "4,emin,40"
  ]
}`

		assert.NoError(t, err)
		assert.Equal(t, expected, outStream.String())

	})
}
