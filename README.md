# csvdiff

[![Build Status](https://travis-ci.org/aswinkarthik/csvdiff.svg?branch=master)](https://travis-ci.org/aswinkarthik/csvdiff)
[![Go Doc](https://godoc.org/github.com/aswinkarthik/csvdiff?status.svg)](https://godoc.org/github.com/aswinkarthik/csvdiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/aswinkarthik/csvdiff)](https://goreportcard.com/report/github.com/aswinkarthik/csvdiff)
[![codecov](https://codecov.io/gh/aswinkarthik/csvdiff/branch/master/graph/badge.svg)](https://codecov.io/gh/aswinkarthik/csvdiff)
[![Downloads](https://img.shields.io/github/downloads/aswinkarthik/csvdiff/total.svg)](https://github.com/aswinkarthik/csvdiff/releases)
[![Latest release](https://img.shields.io/github/release/aswinkarthik/csvdiff.svg)](https://github.com/aswinkarthik/csvdiff/releases)

A fast diff tool for comparing csv files.

## What is csvdiff?

Csvdiff is a difftool to compute changes between two csv files.

- It is not a traditional diff tool. It is **most suitable** for comparing csv files dumped from **database tables**. GNU diff tool is orders of magnitude faster on comparing line by line.
- Supports selective comparison of fields in a row.
- Supports specifying group of columns as primary-key i.e uniquely identify a row.
- Support ignoring columns e.g ignore columns like `created_at` timestamps.
- Compares csvs of million records csv in under 2 seconds.
- Supports lot of output formats e.g colored git style output or JSON for post-processing.

## Why?

I wanted to compare if the rows of a table before and after a given time and see what is the new changes that came in. Also, I wanted to selectively compare columns ignoring columns like `created_at` and `updated_at`. All I had was just the dumped csv files.

## Demo

[![asciicast](https://asciinema.org/a/YNO5G0b2qL92MZWmb2IeiXveN.svg)](https://asciinema.org/a/YNO5G0b2qL92MZWmb2IeiXveN?speed=2&autoplay=1&size=medium&rows=20&cols=150)

## Usage

```diff
$ csvdiff base.csv delta.csv
# Additions (1)
+ 24564,907,completely-newsite.com,com,19827,32902,completely-newsite.com,com,1621,909,19787,32822
# Modifications (1)
- 69,48,aol.com,com,97543,225532,aol.com,com,70,49,97328,224491
+ 69,1048,aol.com,com,97543,225532,aol.com,com,70,49,97328,224491
# Deletions (1)
- 1618,907,deleted-website.com,com,19827,32902,deleted-website.com,com,1621,909,19787,32822
```


```bash
Differentiates two csv files and finds out the additions and modifications.
Most suitable for csv files created from database tables

Usage:
  csvdiff <base-csv> <delta-csv> [flags]

Flags:
      --columns ints          Selectively compare positions in CSV Eg: 1,2. Default is entire row
  -o, --format string         Available (rowmark|rowmark-with-header|json|legacy-json|diff|word-diff|color-words) (default "diff")
  -h, --help                  help for csvdiff
      --ignore-columns ints   Inverse of --columns flag. This cannot be used if --columns are specified
      --include ints          Include positions in CSV to display Eg: 1,2. Default is entire row
  -p, --primary-key ints      Primary key positions of the Input CSV as comma separated values Eg: 1,2 (default [0])
  -s, --separator string      use specific separator (\t, or any one character string) (default ",")
      --time                  Measure time
  -t, --toggle                Help message for toggle
      --version               version for csvdiff
```

## Installation

### Homebrew

```bash
brew tap thecasualcoder/stable
brew install csvdiff
```

### Using binaries

```bash
# binary will be $GOPATH/bin/csvdiff
curl -sfL https://raw.githubusercontent.com/aswinkarthik/csvdiff/master/install.sh | sh -s -- -b $GOPATH/bin

# or install it into ./bin/
curl -sfL https://raw.githubusercontent.com/aswinkarthik/csvdiff/master/install.sh | sh -s

# In alpine linux (as it does not come with curl by default)
wget -O - -q https://raw.githubusercontent.com/aswinkarthik/csvdiff/master/install.sh | sh -s
```

### Using source code

```bash
go get -u github.com/aswinkarthik/csvdiff
```

## Use case

- Cases where you have a base database dump as csv. If you receive the changes as another database dump as csv, this tool can be used to figure out what are the additions and modifications to the original database dump. The `additions.csv` can be used to create an `insert.sql` and with the `modifications.csv` an `update.sql` data migration.
- The delta file can either contain just the changes or the entire table dump along with the changes.

## Supported

- Additions
- Modifications
- Deletions
- Non comma separators

## Not Supported

- Cannot be used as a generic difftool. Requires a column to be used as a primary key from the csv.

## Formats

There are a number of formats supported

- `diff`: Git's diff style
- `word-diff`: Git's --word-diff style 
- `color-words`: Git's --color-words style
- `json`: JSON serialization of result
- `legacy-json`: JSON serialization of result in old format
- `rowmark`: Marks each row with ADDED or MODIFIED status.
- `rowmark-with-header`: Marks each row with ADDED or MODIFIED status. Always print the first line(csv header).

## Miscellaneous features

- The `--primary-key` in an integer array. Specify comma separated positions if the table has a compound key. Using this primary key, it can figure out modifications. If the primary key changes, it is an addition.

```bash
% csvdiff base.csv delta.csv --primary-key 0,1
```

- If you want to compare only few columns in the csv when computing hash,

```bash
% csvdiff base.csv delta.csv --primary-key 0,1 --columns 2
```

- Supports JSON format for post processing

```bash
% csvdiff examples/base-small.csv examples/delta-small.csv --format json | jq '.'
{
  "Additions": [
    "24564,907,completely-newsite.com,com,19827,32902,completely-newsite.com,com,1621,909,19787,32822"
  ],
  "Modifications": [{
    "Original": "69,1048,aol.com,com,97543,225532,aol.com,com,70,49,97328,224491",
    "Current":  "69,1049,aol.com,com,97543,225532,aol.com,com,70,49,97328,224491"
  }],
  "Deletions": [
    "1615,905,deleted-website.com,com,19833,33110,deleted-website.com,com,1613,902,19835,33135"
  ]
}
```

## Build locally

```bash
$ git clone https://github.com/aswinkarthik/csvdiff
$ go get ./...
$ go build

# To run tests
$ go get github.com/stretchr/testify/assert
$ go test -v ./...
```

## Algorithm

- Creates a map of <uint64, uint64> for both base and delta file
  - `key` is a hash of the primary key values as csv
  - `value` is a hash of the entire row
- Two maps as initial processing output
  - base-map
  - delta-map
- The delta map is compared with the base map. As long as primary key is unchanged, they row will have same `key`. An entry in delta map is a
  - **Addition**, if the base-map's does not have a `value`.
  - **Modification**, if the base-map's `value` is different.
  - **Deletions**, if the base-map has no match on the delta map.

## Credits

- Uses 64 bit [xxHash](https://cyan4973.github.io/xxHash/) algorithm, an extremely fast non-cryptographic hash algorithm, for creating the hash. Implementations from [cespare](https://github.com/cespare/xxhash)
- Used [Majestic million](https://blog.majestic.com/development/majestic-million-csv-daily/) data for demo.
