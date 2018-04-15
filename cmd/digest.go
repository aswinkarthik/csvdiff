// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/aswinkarthik93/csv-digest/pkg/digest"
	"github.com/aswinkarthik93/csv-digest/pkg/encoder"
	"github.com/spf13/cobra"
)

// digestCmd represents the digest command
var digestCmd = &cobra.Command{
	Use:   "digest <csv-file>",
	Short: "Takes in a csv and creates a digest of each line",
	Long: `Takes a Csv file and creates a digest for each line.
The tool can output to stdout or a file in plaintext.
It can also serialize the output as a binary file for any other go program to consume directly`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return nil
		} else if len(args) > 1 {
			return errors.New("requires exactly one arg - the csv file")
		}
		return errors.New("requires atleast one arg - the csv file")
	},
	Run: func(cmd *cobra.Command, args []string) {
		runDigest(args[0])
	},
}

func runDigest(csvFile string) {
	config := digest.DigestConfig{
		KeyPositions: primaryKeyPositions(),
		Encoder:      encoder.JsonEncoder{},
		Reader:       os.Stdin,
		Writer:       os.Stdout,
	}

	err := digest.DigestForFile(config)
	if err != nil {
		log.Fatal(err)
	}
}

func primaryKeyPositions() []int {
	return []int{0}
}

func init() {
	rootCmd.AddCommand(digestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// digestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
