# sched
Simple Golang parallel scheduler.

## Usage
### Add one to each element in arr
```go
arr := make([]uint8, 64*1024*1024) // 64M
_, _ = NewTask(arr, func(_ int, x []byte) ([]byte, error) {
    for i := range x {
        arr[i]++
    }
    return arr, nil
}, false, false).Collect(2*1024*1024, true, true) // 2M batch
```
### Get HTTP contents for earh URL
```go
contents, err := NewTask(urls, func(_ int, urls []string) ([]string, error) {
    for i, u := range url {
        str, err := ... // get content
        if err != nil {
            return nil, err
        }
        urls[i] = str   // overlap
    }
    return urls, nil    // return result
}, false, false).Collect(4, false, false) // 4 URLs as a group
```

## Bechmark
A simple self-incrasing process is performed on a 64M uint8 array.
```c
goos: darwin
goarch: arm64
pkg: github.com/fumiama/sched
cpu: Apple M1
BenchmarkSched/single-8         	1000000000	         0.02992 ns/op	2242952686132.55 MB/s	       0 B/op	       0 allocs/op
BenchmarkSched/para512K-8       	1000000000	         0.02483 ns/op	2702751323377.81 MB/s	       0 B/op	       0 allocs/op
BenchmarkSched/para1M-8         	1000000000	         0.02492 ns/op	2693026104055.06 MB/s	       0 B/op	       0 allocs/op
BenchmarkSched/para2M-8         	1000000000	         0.02133 ns/op	3146490138908.33 MB/s	       0 B/op	       0 allocs/op
PASS
ok  	github.com/fumiama/sched	1.119s
```
