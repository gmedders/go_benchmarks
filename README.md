# Go benchmarks

This repo was motivated by some discussions about the overhead of different
ways to call a function in Go. Our eventual goal was to have an HTTP handler
that deserializes JSON, performs an operation that requires a lock (e.g., mutex, `chan`),
and returns a JSON response. Since I'm new to Go, I wanted to bring some data to the discussion.
I took a simple test case, incrementing an integer by 1, and used `go test`s
benchmarking feature to look at some scenarios. The goal was to incrementally
build complexity so that I could understand the relative overhead of different operations. I tested:

- Passing an integer to a function, incrementing the argument, and returning the value.
- Passing an integer [by pointer](http://goinbigdata.com/golang-pass-by-pointer-vs-pass-by-value/)
  and incrementing.
- Passing an integer by pointer to goroutines and using one of the following to control access to the integer:
  - `mutex`
  - `atomic`
  - `chan`
- Encoding the integer in JSON, deserializing to a `map[string]int`,
  incrementing the counter, and returning the serialized result.
- POSTing the JSON-encoded integer to an HTTP server, which does the same
  as above and returns the serialized result in the response.
- A simple POST to a server that performs no operation and immediately responds 200.

The results of those benchmarks performed on my laptop are:

```
[go_benchmark]$ go test -bench=.
goos: darwin
goarch: amd64
BenchmarkIncrementIntByValue-8                 	2000000000	         0.29 ns/op
BenchmarkIncrementIntByRef-8                   	2000000000	         1.58 ns/op
BenchmarkIncrementIntByRefGoroutineMutex-8     	 5000000	       289 ns/op
BenchmarkIncrementIntByRefGoroutineAtomic-8    	 5000000	       308 ns/op
BenchmarkIncrementIntByRefGoroutineChannel-8   	 2000000	       807 ns/op
BenchmarkIncrementIntInJSON-8                  	 1000000	      1890 ns/op
BenchmarkIncrementIntViaHTTP-8                 	   10000	    106486 ns/op
BenchmarkIncrementIntViaHTTPClient-8           	   10000	    113770 ns/op
BenchmarkNoOpHTTP-8                            	   20000	     95364 ns/op
```

As I mentioned these were from some preliminary exploration of Go. If you spot
any errors or have any suggestions, feel free to let me know.
