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
	"io"
	"log"
	"os"
	"strings"
	"sync"

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

var debug bool

var newLine []byte

func init() {
	rootCmd.AddCommand(digestCmd)
	newLine = []byte{'\n'}

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
	digestCmd.Flags().StringVarP(&config.Additions, "additions", "a", "STDOUT", "Output stream for the additions in delta file")
	digestCmd.Flags().StringVarP(&config.Modifications, "modifications", "m", "STDOUT", "Output stream for the modifications in delta file")

	digestCmd.MarkFlagRequired("base")
	digestCmd.MarkFlagRequired("input")
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

	var wg sync.WaitGroup
	baseChannel := make(chan message)
	deltaChannel := make(chan message)

	wg.Add(1)
	go generateInBackground("base", baseConfig, &wg, baseChannel)

	wg.Add(1)
	go generateInBackground("delta", inputConfig, &wg, deltaChannel)

	wg.Add(1)
	go compareInBackgroud(baseChannel, deltaChannel, &wg)

	wg.Wait()
}

type message struct {
	digestMap map[uint64]uint64
	sourceMap map[uint64]string
}

func generateInBackground(name string, config digest.DigestConfig, wg *sync.WaitGroup, channel chan<- message) {
	digest, sourceMap, err := digest.Create(config)
	if err != nil {
		panic(err)
	}

	log.Println("Generated Digest for " + name)
	channel <- message{digestMap: digest, sourceMap: sourceMap}
	close(channel)
	wg.Done()
}

func compareInBackgroud(baseChannel, deltaChannel <-chan message, wg *sync.WaitGroup) {
	baseMessage := <-baseChannel
	deltaMessage := <-deltaChannel

	additions, modifications := digest.Compare(baseMessage.digestMap, deltaMessage.digestMap)

	aWriter := config.AdditionsWriter()
	mWriter := config.ModificationsWriter()
	defer aWriter.Close()
	defer mWriter.Close()

	fmt.Println()
	print("Additions", aWriter, additions, deltaMessage.sourceMap)
	fmt.Println()
	print("Modifications", mWriter, modifications, deltaMessage.sourceMap)
	fmt.Println()
	wg.Done()
}

func print(recordType string, w io.Writer, positions []uint64, content map[uint64]string) {
	log.Println(fmt.Sprintf("%s Count: %d", recordType, len(positions)))

	for _, i := range positions {
		w.Write([]byte(content[i]))
		w.Write(newLine)
	}
}
