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
	"time"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	timed   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "csvdiff <base-csv> <delta-csv>",
	Short: "A diff tool for database tables dumped as csv files",
	Long: `Differentiates two csv files and finds out the additions and modifications.
Most suitable for csv files created from database tables`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate args
		if len(args) != 2 {
			return fmt.Errorf("Pass 2 files. Usage: csvdiff <base-csv> <delta-csv>")
		}

		if err := isValidFile(args[0]); err != nil {
			return err
		}

		if err := isValidFile(args[1]); err != nil {
			return err
		}

		// Validate flags
		if err := config.Validate(); err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if timed {
			defer timeTrack(time.Now(), "csvdiff")
		}

		baseFile := newReadCloser(args[0])
		defer baseFile.Close()
		deltaFile := newReadCloser(args[1])
		defer deltaFile.Close()

		baseConfig := digest.NewConfig(
			baseFile,
			config.GetPrimaryKeys(),
			config.GetValueColumns(),
			config.GetIncludeColumnPositions(),
		)
		deltaConfig := digest.NewConfig(
			deltaFile,
			config.GetPrimaryKeys(),
			config.GetValueColumns(),
			config.GetIncludeColumnPositions(),
		)

		diff, err := digest.Diff(baseConfig, deltaConfig)

		if err != nil {
			fmt.Fprintf(os.Stderr, "csvdiff failed: %v\n", err)
			os.Exit(2)
		}

		config.Formatter().Format(diff, os.Stdout)

		return
	},
}

func isValidFile(path string) error {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", path)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory. Please pass a file", path)
	}

	if err != nil {
		return fmt.Errorf("error reading path: %v", err)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Version = Version()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().IntSliceVarP(&config.PrimaryKeyPositions, "primary-key", "p", []int{0}, "Primary key positions of the Input CSV as comma separated values Eg: 1,2")
	rootCmd.Flags().IntSliceVarP(&config.ValueColumnPositions, "columns", "", []int{}, "Selectively compare positions in CSV Eg: 1,2. Default is entire row")
	rootCmd.Flags().IntSliceVarP(&config.IncludeColumnPositions, "include", "", []int{}, "Include positions in CSV to display Eg: 1,2. Default is entire row")
	rootCmd.Flags().StringVarP(&config.Format, "format", "", "rowmark", "Available (rowmark|json)")

	rootCmd.Flags().BoolVarP(&timed, "time", "", false, "Measure time")
}

func newReadCloser(filename string) io.ReadCloser {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	return file
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Fprintln(os.Stderr, fmt.Sprintf("%s took %s", name, elapsed))
}
