# csvdiff

[![Build Status](https://travis-ci.org/aswinkarthik93/csvdiff.svg?branch=master)](https://travis-ci.org/aswinkarthik93/csvdiff)
[![Go Doc](https://godoc.org/github.com/aswinkarthik93/csvdiff?status.svg)](https://godoc.org/github.com/aswinkarthik93/csvdiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/aswinkarthik93/csvdiff)](https://goreportcard.com/report/github.com/aswinkarthik93/csvdiff)
[![Downloads](https://img.shields.io/github/downloads/aswinkarthik93/csvdiff/latest/total.svg)](https://github.com/aswinkarthik93/csvdiff/releases)
[![Latest release](https://img.shields.io/github/release/aswinkarthik93/csvdiff.svg)](https://github.com/aswinkarthik93/csvdiff/releases)

A Blazingly fast diff tool for comparing csv files.

## What is csvdiff?

Csvdiff is a difftool to compute changes between two csv files.

* It is not a traditional diff tool. It is most suitable for comparing csv files dumped from database tables.
* Supports specifying group of columns as primary-key.
* Supports selective comparison of fields in a row.
* Process a million records csv in under 2 seconds

## Demo

![demo](/demo/csvdiff.gif)

## Usage

```bash
$ csvdiff run --base base.csv --delta delta.csv
# Additions: 1
...

# Modifications: 20
...
```

## Installation

- For MacOS

```bash
curl -sL https://github.com/aswinkarthik93/csvdiff/releases/download/v0.1.2/csvdiff_0.1.2_darwin_amd64.tar.gz | tar xfz -
```

- For centos

```bash
yum install https://github.com/aswinkarthik93/csvdiff/releases/download/v0.1.2/csvdiff_0.1.2_linux_64-bit.rpm
```

- For debian

```
curl -sL https://github.com/aswinkarthik93/csvdiff/releases/download/v0.1.2/csvdiff_0.1.2_linux_64-bit.deb -O
dpkg --install csvdiff_*_linux_64-bit.deb
```

- For Linux

```bash
curl -sL https://github.com/aswinkarthik93/csvdiff/releases/download/v0.1.2/csvdiff_0.1.2_linux_amd64.tar.gz | tar xfz -
```

- For [Windows](https://github.com/aswinkarthik93/csvdiff/releases/download/v0.1.2/csvdiff_0.1.2_windows_amd64.tar.gz)

- Build using Go

```
go get -u github.com/aswinkarthik93/csvdiff
```

## Usecase

- Cases where you have a base database dump as csv. If you receive the changes as another database dump as csv, this tool can be used to figure out what are the additions and modifications to the original database dump. The `additions.csv` can be used to create an `insert.sql` and with the `modifications.csv` an `update.sql` data migration.
- The delta file can either contain just the changes or the entire table dump along with the changes.

## Supported

- Additions
- Modifications

## Not Supported

- Deletions
- Non comma separators
- Cannot be used as a generic difftool. Requires a column to be used as a primary key from the csv.

## Miscellaneous features

- The `--primary-key` in an integer array. Specify comma separated positions if the table has a compound key. Using this primary key, it can figure out modifications. If the primary key changes, it is an addition.

```bash
% csvdiff run --base base.csv --delta delta.csv --primary-key 0,1
```

- If you want to compare only few columns in the csv when computing hash,

```bash
% csvdiff run --base base.csv --delta delta.csv --primary-key 0,1 --value-columns 2
```

- **Additions** and **Modifications** can be written to files directly instead of STDOUT.

```bash
% csvdiff run --base base.csv --delta delta.csv --additions additions.csv --modifications modifications.csv
```

## Build locally

```bash
$ git clone https://github.com/aswinkarthik93/csvdiff
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

## Credits

- Uses 64 bit [xxHash](https://cyan4973.github.io/xxHash/) algorithm, an extremely fast non-cryptographic hash algorithm, for creating the hash. Implementations from [cespare](https://github.com/cespare/xxhash)
- Used [Majestic million](https://blog.majestic.com/development/majestic-million-csv-daily/) data for demo.

_Benchmark tests can be found [here](/benchmark)._
