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
	"sync"
	"time"

	"github.com/aswinkarthik93/csvdiff/pkg/digest"
	spinner "github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// digestCmd represents the digest command
var digestCmd = &cobra.Command{
	Use:   "run",
	Short: "run diff between base-csv and delta-csv",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var (
	debug     bool
	noTime    bool
	noSpinner bool

	loader    *spinner.Spinner
	startTime time.Time
)

var newLine []byte

func init() {
	// rootCmd.AddCommand(digestCmd)
	newLine = []byte{'\n'}

	digestCmd.Flags().StringVarP(&config.Base, "base", "b", "", "The base csv file")
	digestCmd.Flags().StringVarP(&config.Delta, "delta", "d", "", "The delta csv file")
	digestCmd.Flags().IntSliceVarP(&config.PrimaryKeyPositions, "primary-key", "p", []int{0}, "Primary key positions of the Input CSV as comma separated values Eg: 1,2")
	digestCmd.Flags().IntSliceVarP(&config.ValueColumnPositions, "value-columns", "", []int{}, "Value key positions of the Input CSV as comma separated values Eg: 1,2. Default is entire row")
	digestCmd.Flags().StringVarP(&config.Additions, "additions", "a", "STDOUT", "Output stream for the additions in delta file")
	digestCmd.Flags().StringVarP(&config.Modifications, "modifications", "m", "STDOUT", "Output stream for the modifications in delta file")

	digestCmd.Flags().BoolVarP(&debug, "debug", "", false, "Debug mode")
	digestCmd.Flags().BoolVarP(&noTime, "no-time", "", false, "Do not measure time")
	digestCmd.Flags().BoolVarP(&noSpinner, "no-spinner", "", false, "Do not display spinner")

	digestCmd.MarkFlagRequired("base")
	digestCmd.MarkFlagRequired("delta")
}

func run() {
	if !noTime {
		defer timeTrack(time.Now(), "csvdiff")
	}
	if str, err := json.Marshal(config); err == nil && debug {
		fmt.Println(string(str))
	} else if err != nil {
		log.Fatal(err)
	}

	baseConfig := digest.NewConfig(config.GetBaseReader(), config.GetPrimaryKeys(), config.GetValueColumns())

	deltaConfig := digest.NewConfig(config.GetDeltaReader(), config.GetPrimaryKeys(), config.GetValueColumns())

	var wg sync.WaitGroup
	baseChannel := make(chan message)
	deltaChannel := make(chan message)

	if !noSpinner {
		loader = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		loader.Writer = os.Stderr
		loader.Color("cyan")
		loader.Start()
		loader.Suffix = " Computing hashes"
	}

	wg.Add(1)
	go generateInBackground("base", baseConfig, &wg, baseChannel)

	wg.Add(1)
	go generateInBackground("delta", deltaConfig, &wg, deltaChannel)

	wg.Add(1)
	go compareInBackgroud(baseChannel, deltaChannel, &wg)

	wg.Wait()
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

type message struct {
	digestMap map[uint64]uint64
	sourceMap map[uint64]string
}

func generateInBackground(name string, config *digest.Config, wg *sync.WaitGroup, channel chan<- message) {
	digest := digest.Create(config)

	if debug {
		log.Println("Generated Digest for " + name)
	}
	channel <- message{digestMap: digest}
	close(channel)
	wg.Done()
}

func compareInBackgroud(baseChannel, deltaChannel <-chan message, wg *sync.WaitGroup) {
	baseMessage := <-baseChannel
	deltaMessage := <-deltaChannel
	if !noSpinner {
		loader.Suffix = " Comparing hashes"
	}

	additions, modifications := digest.Compare(baseMessage.digestMap, deltaMessage.digestMap)

	aWriter := config.AdditionsWriter()
	mWriter := config.ModificationsWriter()
	defer aWriter.Close()
	defer mWriter.Close()

	if !noSpinner {
		loader.Stop()
	}
	print("Additions", aWriter, additions, deltaMessage.sourceMap)
	print("Modifications", mWriter, modifications, deltaMessage.sourceMap)
	wg.Done()
}

func print(recordType string, w io.Writer, positions []uint64, content map[uint64]string) {
	fmt.Println()
	fmt.Println(fmt.Sprintf("# %s: %d", recordType, len(positions)))
	fmt.Println()

	for _, i := range positions {
		w.Write([]byte(content[i]))
		w.Write(newLine)
	}
}
