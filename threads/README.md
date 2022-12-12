# Threads

## Setup

* `task test` - test in host OS, sends 1000 parallel requests to https://httpbin.org/delay/10
* `task docker-test` - test in linux container, sends 1000 parallel requests to https://httpbin.org/delay/10
* `task test-local-sleeper` - test in docker-compose env, sends 1000 parallel requests to neighbour container with 10s delay

## Measurements

### Windows host 8 CPU threads with GOMAXPROCS=2

```
=== RUN   TestThreads
2022/12/10 15:10:05 OS: windows Arch: amd64, GO_MAX_PROCS: 2, Requests: 1000 x https://httpbin.org/delay/10
2022/12/10 15:10:05 threads: 6, goroutines: 2, fd: -1
2022/12/10 15:10:06 threads: 68, goroutines: 3845, fd: -1
2022/12/10 15:10:07 threads: 68, goroutines: 6667, fd: -1
2022/12/10 15:10:08 threads: 68, goroutines: 6531, fd: -1
2022/12/10 15:10:09 threads: 68, goroutines: 6019, fd: -1
2022/12/10 15:10:10 threads: 68, goroutines: 5913, fd: -1
2022/12/10 15:10:11 threads: 68, goroutines: 5407, fd: -1
2022/12/10 15:10:12 threads: 68, goroutines: 4727, fd: -1
2022/12/10 15:10:13 threads: 68, goroutines: 3827, fd: -1
2022/12/10 15:10:14 threads: 68, goroutines: 3051, fd: -1
2022/12/10 15:10:15 threads: 68, goroutines: 2699, fd: -1
2022/12/10 15:10:16 threads: 68, goroutines: 2149, fd: -1
2022/12/10 15:10:17 threads: 68, goroutines: 843, fd: -1
2022/12/10 15:10:18 threads: 68, goroutines: 517, fd: -1
2022/12/10 15:10:19 threads: 68, goroutines: 213, fd: -1
2022/12/10 15:10:20 threads: 68, goroutines: 213, fd: -1
2022/12/10 15:10:21 threads: 68, goroutines: 213, fd: -1
2022/12/10 15:10:22 threads: 68, goroutines: 107, fd: -1
2022/12/10 15:10:23 threads: 68, goroutines: 105, fd: -1
2022/12/10 15:10:24 threads: 68, goroutines: 105, fd: -1
2022/12/10 15:10:25 threads: 68, goroutines: 29, fd: -1
2022/12/10 15:10:26 threads: 68, goroutines: 11, fd: -1
2022/12/10 15:10:26 Response codes: map[200:999 502:1]
--- PASS: TestThreads (20.81s)
```

### Linux docker container 2 CPU threads

```
=== RUN   TestThreads
2022/12/10 12:07:44 OS: linux Arch: amd64, GO_MAX_PROCS: 2, Requests: 1000 x https://httpbin.org/delay/10
2022/12/10 12:07:44 threads: 5, goroutines: 2, fd: 8
2022/12/10 12:07:45 threads: 6, goroutines: 4293, fd: 914
2022/12/10 12:07:46 threads: 6, goroutines: 6561, fd: 2291
2022/12/10 12:07:47 threads: 6, goroutines: 6491, fd: 2256
2022/12/10 12:07:48 threads: 6, goroutines: 6359, fd: 2190
2022/12/10 12:07:49 threads: 6, goroutines: 5887, fd: 1954
2022/12/10 12:07:50 threads: 6, goroutines: 5699, fd: 1860
2022/12/10 12:07:51 threads: 6, goroutines: 4915, fd: 1468
2022/12/10 12:07:52 threads: 6, goroutines: 4505, fd: 1263
2022/12/10 12:07:53 threads: 6, goroutines: 3845, fd: 933
2022/12/10 12:07:54 threads: 6, goroutines: 3079, fd: 550
2022/12/10 12:07:55 threads: 6, goroutines: 2535, fd: 278
2022/12/10 12:07:56 threads: 6, goroutines: 587, fd: 185
2022/12/10 12:07:57 threads: 6, goroutines: 343, fd: 182
2022/12/10 12:07:57 Response codes: map[200:1000]
--- PASS: TestThreads (12.81s)
```

### Linux docker container 2 CPU threads, test with local sleeping server

```
=== RUN   TestThreads
2022/12/10 12:33:04 OS: linux Arch: amd64, GO_MAX_PROCS: 2, Requests: 1000 x http://sleeper:8080/?delay=10s
2022/12/10 12:33:04 threads: 5, goroutines: 2, fd: 8
2022/12/10 12:33:05 threads: 5, goroutines: 3003, fd: 1008
2022/12/10 12:33:06 threads: 5, goroutines: 3003, fd: 1008
2022/12/10 12:33:07 threads: 5, goroutines: 3003, fd: 1008
2022/12/10 12:33:08 threads: 5, goroutines: 3003, fd: 1008
2022/12/10 12:33:09 threads: 5, goroutines: 3003, fd: 1008
2022/12/10 12:33:10 threads: 5, goroutines: 3003, fd: 1008
2022/12/10 12:33:11 threads: 5, goroutines: 3003, fd: 1008
2022/12/10 12:33:12 threads: 5, goroutines: 3003, fd: 1008
2022/12/10 12:33:13 threads: 5, goroutines: 3003, fd: 1008
2022/12/10 12:33:14 threads: 5, goroutines: 3003, fd: 1008
2022/12/10 12:33:15 threads: 6, goroutines: 5, fd: 8
2022/12/10 12:33:15 Response codes: map[200:1000]
--- PASS: TestThreads (10.22s)
```

## Useful links

* [Post about Golang and system threads](https://www.sobyte.net/post/2021-06/golang-number-of-threads-in-the-running-program/)