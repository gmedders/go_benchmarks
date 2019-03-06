Some benchmarks of performing simple operations in Go.

```
BenchmarkIncrementIntByValue-8                 	2000000000	         0.30 ns/op
BenchmarkIncrementIntByRef-8                   	2000000000	         1.59 ns/op
BenchmarkIncrementIntByRefGoroutineMutex-8     	 5000000	       296 ns/op
BenchmarkIncrementIntByRefGoroutineAtomic-8    	 5000000	       292 ns/op
BenchmarkIncrementIntByRefGoroutineChannel-8   	 2000000	       912 ns/op
```

### Some useful commands

```
go test
go test -race
go test -bench=.
```
