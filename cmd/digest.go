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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aswinkarthik93/csv-digest/pkg/digest"
	"github.com/spf13/cobra"
)

// digestCmd represents the digest command
var digestCmd = &cobra.Command{
	Use:   "digest <csv-file>",
	Short: "Takes in a csv and creates a digest of each line",
	Long: `Takes a Csv file and creates a digest for each line.
The tool can output to stdout or a file in plaintext.
It can also serialize the output as a binary file for any other go program to consume directly`,
	Run: func(cmd *cobra.Command, args []string) {
		runDigest()
	},
}

func runDigest() {
	if str, err := json.Marshal(config); err == nil && debug {
		fmt.Println(string(str))
	} else if err != nil {
		log.Fatal(err)
	}

	baseConfig := digest.DigestConfig{
		KeyPositions: config.GetKeyPositions(),
		Encoder:      config.GetEncoder(),
		Reader:       config.GetBase(),
		Writer:       os.Stdout,
	}

	inputConfig := digest.DigestConfig{
		KeyPositions: config.GetKeyPositions(),
		Encoder:      config.GetEncoder(),
		Reader:       config.GetInput(),
		Writer:       os.Stdout,
		SourceMap:    true,
	}

	base, _, err := digest.Create(baseConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Generated Base Digest")

	change, sourceMap, err := digest.Create(inputConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Generated Input file Digest")

	additions, modifications := digest.Compare(base, change)

	fmt.Println(fmt.Sprintf("Additions Count: %d", len(additions)))
	for _, addition := range additions {
		fmt.Println(sourceMap[addition])
	}

	fmt.Println("")
	fmt.Println(fmt.Sprintf("Modifications Count: %d", len(modifications)))
	for _, modification := range modifications {
		fmt.Println(sourceMap[modification])
	}
}

var debug bool

func init() {
	rootCmd.AddCommand(digestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// digestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	digestCmd.Flags().StringVarP(&config.Base, "base", "b", "", "Input csv to be used as base")
	digestCmd.Flags().StringVarP(&config.Input, "input", "i", "", "The new csv file on which diff should be done")
	digestCmd.Flags().StringVarP(&config.Encoder, "encoder", "e", "json", "Encoder to use to output the digest. Available Encoders: "+strings.Join(GetEncoders(), ","))
	digestCmd.Flags().IntSliceVarP(&config.KeyPositions, "key-positions", "k", []int{0}, "Primary key positions of the Input CSV as comma separated values Eg: 1,2")
	digestCmd.Flags().BoolVarP(&debug, "debug", "", false, "Debug mode")

	digestCmd.MarkFlagRequired("base")
	digestCmd.MarkFlagRequired("input")
}
