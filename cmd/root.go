// Copyright © 2018 aswinkarthik
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/spf13/afero"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/spf13/cobra"
)

var (
	timed bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "csvdiff <base-csv> <delta-csv>",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "A diff tool for database tables dumped as csv files",
	Long: `Differentiates two csv files and finds out the additions and modifications.
Most suitable for csv files created from database tables`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// validate args
		if len(args) != 2 {
			return fmt.Errorf("pass 2 files. Usage: csvdiff <base-csv> <delta-csv>")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if timed {
			defer timeTrack(time.Now(), "csvdiff")
		}
		fs := afero.NewOsFs()
		baseFilename := args[0]
		deltaFilename := args[1]
		runeSeparator, err := parseSeparator(separator)
		if err != nil {
			return err
		}
		ctx, err := NewContext(
			fs,
			primaryKeyPositions,
			valueColumnPositions,
			ignoreValueColumnPositions,
			includeColumnPositions,
			format,
			baseFilename,
			deltaFilename,
			runeSeparator,
			lazyQuotes,
		)

		if err != nil {
			return err
		}
		defer ctx.Close()

		return runContext(ctx, os.Stdout, os.Stderr)
	},
}

func runContext(ctx *Context, outputStream, errorStream io.Writer) error {
	baseConfig, err := ctx.BaseDigestConfig()
	if err != nil {
		return fmt.Errorf("error opening base-file %s: %v", ctx.baseFilename, err)
	}
	deltaConfig, err := ctx.DeltaDigestConfig()
	if err != nil {
		return fmt.Errorf("error opening delta-file %s: %v", ctx.deltaFilename, err)
	}
	defer ctx.Close()

	diff, err := digest.Diff(baseConfig, deltaConfig)

	if err != nil {
		return err
	}

	return NewFormatter(outputStream, errorStream, *ctx).Format(diff)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Version = Version()
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprint(os.Stderr, color.RedString("csvdiff: command failed - %v\n\n", err))
		_ = rootCmd.Help()
		os.Exit(1)
	}
}

var (
	primaryKeyPositions        []int
	valueColumnPositions       []int
	ignoreValueColumnPositions []int
	includeColumnPositions     []int
	format                     string
	separator                  string
	lazyQuotes                 bool
)

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().IntSliceVarP(&primaryKeyPositions, "primary-key", "p", []int{0}, "Primary key positions of the Input CSV as comma separated values Eg: 1,2")
	rootCmd.Flags().IntSliceVarP(&valueColumnPositions, "columns", "", []int{}, "Selectively compare positions in CSV Eg: 1,2. Default is entire row")
	rootCmd.Flags().IntSliceVarP(&ignoreValueColumnPositions, "ignore-columns", "", []int{}, "Inverse of --columns flag. This cannot be used if --columns are specified")
	rootCmd.Flags().IntSliceVarP(&includeColumnPositions, "include", "", []int{}, "Include positions in CSV to display Eg: 1,2. Default is entire row")
	rootCmd.Flags().StringVarP(&format, "format", "o", "diff", fmt.Sprintf("Available (%s)", strings.Join(allFormats, "|")))
	rootCmd.Flags().StringVarP(&separator, "separator", "s", ",", "use specific separator (\\t, or any one character string)")

	rootCmd.Flags().BoolVarP(&timed, "time", "", false, "Measure time")
	rootCmd.Flags().BoolVar(&lazyQuotes, "lazyquotes", false, "allow unescaped quotes")
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("%s took %s", name, elapsed))
}

func parseSeparator(sep string) (rune, error) {
	if strings.HasPrefix(sep, "\\t") {
		return '\t', nil
	}

	runesep, _ := utf8.DecodeRuneInString(sep)
	if runesep == utf8.RuneError {
		return ' ', fmt.Errorf("unable to use %v (%q) as a separator", separator, separator)
	}

	return runesep, nil
}
