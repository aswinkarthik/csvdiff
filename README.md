# csvdiff

A diff tool for database tables dumped as csv files.

## Usage

```bash
$ csvdiff run --base base.csv --delta delta.csv
2018/04/15 18:46:23 Generated Digest for delta
2018/04/15 18:46:24 Generated Digest for base
2018/04/15 18:46:24 Additions Count: 1
...

2018/04/15 18:46:24 Modifications Count: 20
...
```

## Usecase

- Cases where you have a base database dump as csv. If you receive the changes as another database dump as csv, this tool can be used to figure out what are the additions and modifications to the original database dump. The `additions.csv` can be used to create an `insert.sql` and with the `modifications.csv` an `update.sql` data migration.
- As the delta file, it supports passing of just the changes as well as the entire csv file.

## Supported

- Additions
- Modifications

## Not Supported

- Deletions
- Non comma separators
- Cannot be used as a generic difftool. Requires a column to be used as a primary key from the csv.

## Example

Dataset is used from this [blog](https://blog.majestic.com/development/majestic-million-csv-daily/)

- Base csv file

```bash
% cat ./examples/base-small.csv
15,12,wordpress.com,com,207790,792348,wordpress.com,com,15,12,207589,791634
43,1,europa.eu,eu,116613,353412,europa.eu,eu,41,1,119129,359818
69,48,aol.com,com,97543,225532,aol.com,com,70,49,97328,224491
1615,905,proboards.com,com,19833,33110,proboards.com,com,1613,902,19835,33135
1616,906,ccleaner.com,com,19831,32507,ccleaner.com,com,1614,903,19834,32463
1617,907,doodle.com,com,19827,32902,doodle.com,com,1621,909,19787,32822
```

- Delta csv file

```bash
% cat ./examples/delta-small.csv
15,12,wordpress.com,com,207790,792348,wordpress.com,com,15,12,207589,791634
43,1,europa.eu,eu,116613,353412,europa.eu,eu,41,1,119129,359818
69,1048,aol.com,com,97543,225532,aol.com,com,70,49,97328,224491
24564,907,completely-newsite.com,com,19827,32902,completely-newsite.com,com,1621,909,19787,32822
```

- On run of csvdiff

```bash
% csvdiff run --base ./examples/base-small.csv --delta ./examples/delta-small.csv --key-positions 0
2018/04/15 19:15:48 Generated Digest for base
2018/04/15 19:15:48 Generated Digest for delta
2018/04/15 19:15:48 Additions Count: 1
24564,907,completely-newsite.com,com,19827,32902,completely-newsite.com,com,1621,909,19787,32822
2018/04/15 19:15:48 Modifications Count: 1
69,1048,aol.com,com,97543,225532,aol.com,com,70,49,97328,224491
```

- The `--key-positions` in an integer array. Specify comma separated positions if the table has a compound key. Using this primary key, it can figure out modifications. If the primary key changes, it is an addition.

```bash
% csvdiff run --base base.csv --delta delta.csv --key-positions 0,1
```

- **Additions** and **Modifications** can be written to files directly instead of STDOUT.

```bash
% csvdiff run --base base.csv --delta delta.csv --additions additions.csv --modifications modifications.csv
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
