library(data.table)

csv1 = fread('majestic_million.csv')
csv2 = fread('majestic_million_diff.csv')

setkey(csv1,id)
setkey(csv2,id)

result <- merge(csv2, csv1, all.x=TRUE)

diff <- result[result$"col-1.x" != result$"col-1.y" | result$"col-2.x" != result$"col-2.y" | result$"col-3.x" != result$"col-3.y" | result$"col-4.x" != result$"col-4.y" | result$"col-5.x" != result$"col-5.y" | result$"col-6.x" != result$"col-6.y" | result$"col-7.x" != result$"col-7.y" | result$"col-8.x" != result$"col-8.y" | result$"col-9.x" != result$"col-9.y" | result$"col-10.x" != result$"col-10.y" | result$"col-11.x" != result$"col-11.y"]

diff
