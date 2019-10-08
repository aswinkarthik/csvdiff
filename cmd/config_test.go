package cmd_test

import (
	"github.com/spf13/afero"
	"os"
	"testing"

	"github.com/aswinkarthik/csvdiff/cmd"
	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/stretchr/testify/assert"
)

func TestPrimaryKeyPositions(t *testing.T) {
	config := cmd.Context{PrimaryKeyPositions: []int{0, 1}}
	assert.Equal(t, digest.Positions([]int{0, 1}), config.GetPrimaryKeys())

	config = cmd.Context{PrimaryKeyPositions: []int{}}
	assert.Equal(t, digest.Positions([]int{0}), config.GetPrimaryKeys())

	config = cmd.Context{}
	assert.Equal(t, digest.Positions([]int{0}), config.GetPrimaryKeys())
}

func TestValueColumnPositions(t *testing.T) {
	config := cmd.Context{ValueColumnPositions: []int{0, 1}}
	assert.Equal(t, digest.Positions([]int{0, 1}), config.GetValueColumns())

	config = cmd.Context{ValueColumnPositions: []int{}}
	assert.Equal(t, digest.Positions([]int{}), config.GetValueColumns())

	config = cmd.Context{}
	assert.Equal(t, digest.Positions([]int{}), config.GetValueColumns())
}

func TestConfigValidate(t *testing.T) {
	validConfig := func(t *testing.T, fs afero.Fs) *cmd.Context {
		_, err := fs.Create("/base.csv")
		assert.NoError(t, err)
		_, err = fs.Create("/delta.csv")
		assert.NoError(t, err)
		return &cmd.Context{Format: "json", BaseFilename: "/base.csv", DeltaFilename: "/delta.csv"}
	}

	t.Run("should validate format", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		config := validConfig(t, fs)

		config.Format = ""
		assert.Error(t, config.Validate(fs))

		config.Format = "rowmark"
		assert.NoError(t, config.Validate(fs))

		config.Format = "rowMARK"
		assert.NoError(t, config.Validate(fs))

		config.Format = "json"
		assert.NoError(t, config.Validate(fs))
	})

	t.Run("should validate base file existence", func(t *testing.T) {
		fs := afero.NewMemMapFs()

		config := &cmd.Context{Format: "json", BaseFilename: "/base.csv", DeltaFilename: "/delta.csv"}
		err := config.Validate(fs)
		assert.EqualError(t, err, "base-file /base.csv does not exits")
	})

	t.Run("should validate if base file or delta file is a file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		err := fs.Mkdir("/base.csv", os.ModePerm)
		assert.NoError(t, err)

		config := &cmd.Context{Format: "json", BaseFilename: "/base.csv", DeltaFilename: "/delta.csv"}
		err = config.Validate(fs)
		assert.EqualError(t, err, "base-file /base.csv should be a file")

		_, err = fs.Create("/valid-base.csv")
		err = fs.Mkdir("/delta.csv", os.ModePerm)
		assert.NoError(t, err)

		config = &cmd.Context{Format: "json", BaseFilename: "/valid-base.csv", DeltaFilename: "/delta.csv"}
		err = config.Validate(fs)
		assert.EqualError(t, err, "delta-file /delta.csv should be a file")
	})

	t.Run("should validate if both base and delta file exist", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, err := fs.Create("/base.csv")
		assert.NoError(t, err)
		_, err = fs.Create("/delta.csv")
		assert.NoError(t, err)

		config := &cmd.Context{Format: "json", BaseFilename: "/base.csv", DeltaFilename: "/delta.csv"}
		err = config.Validate(fs)
		assert.NoError(t, err)
	})
}

func TestConfig_DigestConfig(t *testing.T) {
	t.Run("should create digest ctx", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_, err := fs.Create("/base.csv")
		assert.NoError(t, err)
		_, err = fs.Create("/delta.csv")
		assert.NoError(t, err)

		valueColumns := digest.Positions{0, 1, 2}
		primaryColumns := digest.Positions{0, 1}
		includeColumns := digest.Positions{2}
		config := &cmd.Context{
			Format:                 "json",
			BaseFilename:           "/base.csv",
			DeltaFilename:          "/delta.csv",
			ValueColumnPositions:   valueColumns,
			PrimaryKeyPositions:    primaryColumns,
			IncludeColumnPositions: includeColumns,
		}

		baseConfig, err := config.BaseDigestConfig(fs)

		assert.NoError(t, err)
		assert.NotNil(t, baseConfig.Reader)
		assert.Equal(t, baseConfig.Value, valueColumns)
		assert.Equal(t, baseConfig.Key, primaryColumns)
		assert.Equal(t, baseConfig.Include, includeColumns)

		deltaConfig, err := config.DeltaDigestConfig(fs)

		assert.NoError(t, err)
		assert.NotNil(t, deltaConfig.Reader)
		assert.Equal(t, deltaConfig.Value, valueColumns)
		assert.Equal(t, deltaConfig.Key, primaryColumns)
		assert.Equal(t, deltaConfig.Include, includeColumns)
	})
}
