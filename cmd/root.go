// Copyright © 2018 aswinkarthik93
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
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aswinkarthik93/csvdiff/pkg/digest"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var timed bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "csvdiff <base-csv> <delta-csv>",
	Short: "A diff tool for database tables dumped as csv files",
	Long: `Differentiates two csv files and finds out the additions and modifications.
Most suitable for csv files created from database tables`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Pass 2 files. Usage: csvdiff <base-csv> <delta-csv>")
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if timed {
			defer timeTrack(time.Now(), "csvdiff")
		}

		baseFile := newReadCloser(args[0])
		defer baseFile.Close()
		deltaFile := newReadCloser(args[1])
		defer deltaFile.Close()

		baseConfig := digest.NewConfig(baseFile, config.GetPrimaryKeys(), config.GetValueColumns())
		deltaConfig := digest.NewConfig(deltaFile, config.GetPrimaryKeys(), config.GetValueColumns())

		diff := digest.Diff(baseConfig, deltaConfig)

		fmt.Printf("Additions %d\n", len(diff.Additions))
		fmt.Printf("Modifications %d\n", len(diff.Modifications))
		fmt.Println("Rows:")

		for _, added := range diff.Additions {
			fmt.Printf("%s,%s\n", added, "ADDED")
		}

		for _, modified := range diff.Modifications {
			fmt.Printf("%s,%s\n", modified, "MODIFIED")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.csvdiff.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().IntSliceVarP(&config.PrimaryKeyPositions, "primary-key", "p", []int{0}, "Primary key positions of the Input CSV as comma separated values Eg: 1,2")
	rootCmd.Flags().IntSliceVarP(&config.ValueColumnPositions, "columns", "", []int{}, "Selectively compare positions in CSV Eg: 1,2. Default is entire row")

	rootCmd.Flags().BoolVarP(&timed, "time", "", false, "Measure time")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".csvdiff" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".csvdiff")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
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
