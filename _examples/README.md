# Colly examples

This folder provides easy to understand code snippets on how to get started with colly.

To execute an example run `go run [example/example.go]`


## Demo

```
$ go run rate_limit/rate_limit.go
[000001] 1 [     1 - request] map["url":"https://httpbin.org/delay/2?n=4"] (60.872µs)
[000002] 1 [     2 - request] map["url":"https://httpbin.org/delay/2?n=2"] (154.425µs)
[000003] 1 [     3 - request] map["url":"https://httpbin.org/delay/2?n=0"] (158.374µs)
[000004] 1 [     5 - request] map["url":"https://httpbin.org/delay/2?n=3"] (426.999µs)
[000005] 1 [     4 - request] map["url":"https://httpbin.org/delay/2?n=1"] (448.75µs)
[000007] 1 [     2 - response] map["url":"https://httpbin.org/delay/2?n=2" "status":"OK"] (2.855764394s)
[000008] 1 [     2 - scraped] map["url":"https://httpbin.org/delay/2?n=2"] (2.855797868s)
[000006] 1 [     1 - response] map["url":"https://httpbin.org/delay/2?n=4" "status":"OK"] (2.855756753s)
[000009] 1 [     1 - scraped] map["url":"https://httpbin.org/delay/2?n=4"] (2.855819581s)
[000010] 1 [     3 - response] map["status":"OK" "url":"https://httpbin.org/delay/2?n=0"] (5.002065299s)
[000011] 1 [     3 - scraped] map["url":"https://httpbin.org/delay/2?n=0"] (5.002103755s)
[000012] 1 [     5 - response] map["status":"OK" "url":"https://httpbin.org/delay/2?n=3"] (5.012080614s)
[000013] 1 [     5 - scraped] map["url":"https://httpbin.org/delay/2?n=3"] (5.012101056s)
[000014] 1 [     4 - response] map["url":"https://httpbin.org/delay/2?n=1" "status":"OK"] (7.155725591s)
[000015] 1 [     4 - scraped] map["url":"https://httpbin.org/delay/2?n=1"] (7.155759136s)

```
