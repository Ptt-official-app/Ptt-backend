## Router benchmark

Ptt-backend/misc/router-benchmark\
`go test -bench=. -benchmem`
```go
goos: linux
goarch: amd64
pkg: github.com/Ptt-official-app/Ptt-backend/misc/router-benchmark
cpu: Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz
Benchmark_ServeMux-7             3167156               374.2 ns/op            80 B/op          1 allocs/op
Benchmark_gorillamux-7            760047              1560 ns/op            1312 B/op         10 allocs/op
Benchmark_httprouter-7           3195236               371.7 ns/op           504 B/op          5 allocs/op
```
