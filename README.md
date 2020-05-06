# Monitoring Tool
Used to log and monitor activities on the network by attaching to heimdall and bor nodes

### Build
```
go build
```

### Benchmarks
1. Subscribe to new blocks. Prints tx count in each block
```
./monitoring-tools txcount
```
2. Fire txs in Bulk
```
bash benchmarking/benchmark.sh
```
3. [Benchmarking Results](./benchmarking/README.md)
