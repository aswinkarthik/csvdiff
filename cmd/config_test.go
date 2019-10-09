package cmd_test

import (
	"os"
	"testing"

	"github.com/aswinkarthik/csvdiff/cmd"
	"github.com/spf13/afero"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/stretchr/testify/assert"
)

func TestPrimaryKeyPositions(t *testing.T) {
	type testCase struct {
		name string
		in   []int
		out  digest.Positions
	}
	testCases := []testCase{
		{
			name: "should return primary key columns",
			in:   []int{0, 1},
			out:  []int{0, 1},
		},
		{
			name: "should return primary key columns as default input is empty",
			in:   []int{},
			out:  []int{0},
		},
		{
			name: "should return primary key columns as default input is nil",
			in:   []int{},
			out:  []int{0},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			setupFiles(t, fs)
			ctx, err := cmd.NewContext(fs,
				tt.in,
				nil,
				nil,
				nil,
				"json",
				"/base.csv",
				"/delta.csv",
			)
			assert.NoError(t, err)
			assert.Equal(t, tt.out, ctx.GetPrimaryKeys())

		})
	}
}

func TestValueColumnPositions(t *testing.T) {
	type testCase struct {
		name string
		in   []int
		out  digest.Positions
	}
	testCases := []testCase{
		{
			name: "should return value columns",
			in:   []int{0, 1},
			out:  []int{0, 1},
		},
		{
			name: "should return value columns as empty if input is empty",
			in:   []int{},
			out:  []int{},
		},
		{
			name: "should return value columns as empty if input is nil",
			in:   []int{},
			out:  []int{},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			setupFiles(t, fs)
			ctx, err := cmd.NewContext(fs,
				nil,
				tt.in,
				nil,
				nil,
				"json",
				"/base.csv",
				"/delta.csv",
			)
			assert.NoError(t, err)
			assert.Equal(t, tt.out, ctx.GetValueColumns())

		})
	}
}

func TestNewContext(t *testing.T) {

	t.Run("should validate format", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		setupFiles(t, fs)

		t.Run("empty format", func(t *testing.T) {
			_, err := cmd.NewContext(
				fs,
				nil,
				nil,
				nil,
				nil,
				"",
				"/base.csv",
				"/delta.csv",
			)

			assert.EqualError(t, err, "validation failed: specified format is not valid")
		})

		t.Run("valid format", func(t *testing.T) {
			_, err := cmd.NewContext(
				fs,
				nil,
				nil,
				nil,
				nil,
				"rowmark",
				"/base.csv",
				"/delta.csv",
			)

			assert.NoError(t, err)
		})

		t.Run("case-insensitive valid format", func(t *testing.T) {
			_, err := cmd.NewContext(
				fs,
				nil,
				nil,
				nil,
				nil,
				"jSOn",
				"/base.csv",
				"/delta.csv",
			)

			assert.NoError(t, err)
		})

	})

	t.Run("should validate base file existence", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, err := cmd.NewContext(
			fs,
			nil,
			nil,
			nil,
			nil,
			"json",
			"/base.csv",
			"/delta.csv",
		)
		assert.EqualError(t, err, "error in base-file: open "+string(os.PathSeparator)+"base.csv: file does not exist")
	})

	t.Run("should validate if base file is a csv file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		{
			err := fs.Mkdir("/base.csv", os.ModePerm)
			assert.NoError(t, err)
		}

		_, err := cmd.NewContext(
			fs,
			nil,
			nil,
			nil,
			nil,
			"json",
			"/base.csv",
			"/delta.csv",
		)
		assert.EqualError(t, err, "error in base-file: unable to process headers from csv file. EOF reached. invalid CSV file")
	})
	t.Run("should validate if delta file is a csv file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		{
			assert.NoError(t, afero.WriteFile(fs, "/base.csv", []byte("id"), os.ModePerm))
			err := fs.Mkdir("/delta.csv", os.ModePerm)
			assert.NoError(t, err)
		}

		_, err := cmd.NewContext(
			fs,
			nil,
			nil,
			nil,
			nil,
			"json",
			"/base.csv",
			"/delta.csv",
		)
		assert.EqualError(t, err, "error in delta-file: unable to process headers from csv file. EOF reached. invalid CSV file")
	})

	t.Run("should validate if both base and delta file exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		setupFiles(t, fs)

		_, err := cmd.NewContext(
			fs,
			nil,
			nil,
			nil,
			nil,
			"json",
			"/base.csv",
			"/delta.csv",
		)
		assert.NoError(t, err)
	})

	t.Run("should validate if positions are within the limits of the csv file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		{
			baseContent := []byte("id,name,age,desc")
			err := afero.WriteFile(fs, "/base.csv", baseContent, os.ModePerm)
			assert.NoError(t, err)
		}
		{
			deltaContent := []byte("id,name,age,desc")
			err := afero.WriteFile(fs, "/delta.csv", deltaContent, os.ModePerm)
			assert.NoError(t, err)
		}

		t.Run("primary key positions", func(t *testing.T) {
			_, err := cmd.NewContext(
				fs,
				[]int{4},
				nil,
				nil,
				nil,
				"json",
				"/base.csv",
				"/delta.csv",
			)

			assert.EqualError(t, err, "validation failed: --primary-key positions are out of bounds")
		})

		t.Run("include positions", func(t *testing.T) {
			_, err := cmd.NewContext(
				fs,
				nil,
				nil,
				nil,
				[]int{4},
				"json",
				"/base.csv",
				"/delta.csv",
			)

			assert.EqualError(t, err, "validation failed: --include positions are out of bounds")
		})

		t.Run("value positions", func(t *testing.T) {
			_, err := cmd.NewContext(
				fs,
				nil,
				[]int{4},
				nil,
				nil,
				"json",
				"/base.csv",
				"/delta.csv",
			)

			assert.EqualError(t, err, "validation failed: --columns positions are out of bounds")
		})

		t.Run("inequal base and delta files", func(t *testing.T) {
			{
				deltaContent := []byte("id,name,age,desc,size")
				err := afero.WriteFile(fs, "/delta.csv", deltaContent, os.ModePerm)
				assert.NoError(t, err)
			}

			_, err := cmd.NewContext(
				fs,
				nil,
				nil,
				nil,
				nil,
				"json",
				"/base.csv",
				"/delta.csv",
			)
			assert.EqualError(t, err, "base-file and delta-file columns count do not match")
		})
	})

	t.Run("should pass only one of columns or ignore columns", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		setupFiles(t, fs)

		_, err := cmd.NewContext(
			fs,
			nil,
			[]int{0},
			[]int{0},
			nil,
			"jSOn",
			"/base.csv",
			"/delta.csv",
		)

		assert.EqualError(t, err, "only one of --columns or --ignore-columns")
	})
}

func TestConfig_DigestConfig(t *testing.T) {
	t.Run("should create digest ctx", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		setupFiles(t, fs)

		valueColumns := digest.Positions{0, 1, 2}
		primaryColumns := digest.Positions{0, 1}
		includeColumns := digest.Positions{2}
		ctx, err := cmd.NewContext(
			fs,
			primaryColumns,
			valueColumns,
			nil,
			includeColumns,
			"jSOn",
			"/base.csv",
			"/delta.csv",
		)
		assert.NoError(t, err)

		baseConfig, err := ctx.BaseDigestConfig()

		assert.NoError(t, err)
		assert.NotNil(t, baseConfig.Reader)
		assert.Equal(t, valueColumns, baseConfig.Value)
		assert.Equal(t, primaryColumns, baseConfig.Key)
		assert.Equal(t, includeColumns, baseConfig.Include)

		deltaConfig, err := ctx.DeltaDigestConfig()

		assert.NoError(t, err)
		assert.NotNil(t, deltaConfig.Reader)
		assert.Equal(t, valueColumns, deltaConfig.Value)
		assert.Equal(t, primaryColumns, deltaConfig.Key)
		assert.Equal(t, includeColumns, deltaConfig.Include)
	})
	t.Run("should infer values columns as inverse of ignore columns digest ctx", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		setupFiles(t, fs)

		ignoreValueColumns := digest.Positions{0, 1, 2}
		primaryColumns := digest.Positions{0, 1}
		ctx, err := cmd.NewContext(
			fs,
			primaryColumns,
			nil,
			ignoreValueColumns,
			nil,
			"jSOn",
			"/base.csv",
			"/delta.csv",
		)
		assert.NoError(t, err)

		baseConfig, err := ctx.BaseDigestConfig()

		assert.NoError(t, err)
		assert.NotNil(t, baseConfig.Reader)
		assert.Equal(t, digest.Positions{3}, baseConfig.Value)
		assert.Equal(t, primaryColumns, baseConfig.Key)

		deltaConfig, err := ctx.DeltaDigestConfig()

		assert.NoError(t, err)
		assert.NotNil(t, deltaConfig.Reader)
		assert.Equal(t, digest.Positions{3}, deltaConfig.Value)
		assert.Equal(t, primaryColumns, deltaConfig.Key)
	})
}
func setupFiles(t *testing.T, fs afero.Fs) {
	{
		baseContent := []byte("id,name,age,desc")
		err := afero.WriteFile(fs, "/base.csv", baseContent, os.ModePerm)
		assert.NoError(t, err)
	}
	{
		deltaContent := []byte("id,name,age,desc")
		err := afero.WriteFile(fs, "/delta.csv", deltaContent, os.ModePerm)
		assert.NoError(t, err)
	}
}
