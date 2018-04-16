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