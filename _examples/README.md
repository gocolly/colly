# Colly examples

This folder provides easy to understand code snippets on how to get started with colly.

To execute an example run `go run [example/example.go]`


## Demo

```
$ go run rate_limit/rate_limit.go
Starting https://httpbin.org/delay/2     2017-10-25 02:32:12.542918968 +0200 CEST m=+0.001826149
Starting https://httpbin.org/delay/2?n=1 2017-10-25 02:32:12.543011175 +0200 CEST m=+0.001918365
Starting https://httpbin.org/delay/2?n=2 2017-10-25 02:32:12.543070662 +0200 CEST m=+0.001977846
Starting https://httpbin.org/delay/2?n=0 2017-10-25 02:32:12.543141774 +0200 CEST m=+0.002048966
Starting https://httpbin.org/delay/2?n=3 2017-10-25 02:32:12.543234032 +0200 CEST m=+0.002141228
Finished https://httpbin.org/delay/2     2017-10-25 02:32:15.991943006 +0200 CEST m=+3.450850142
Finished https://httpbin.org/delay/2?n=1 2017-10-25 02:32:16.003512763 +0200 CEST m=+3.462419895
Finished https://httpbin.org/delay/2?n=2 2017-10-25 02:32:18.119386433 +0200 CEST m=+5.578293577
Finished https://httpbin.org/delay/2?n=0 2017-10-25 02:32:18.135157808 +0200 CEST m=+5.594064941
Finished https://httpbin.org/delay/2?n=3 2017-10-25 02:32:20.247735566 +0200 CEST m=+7.706642699
```
