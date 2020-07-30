# Benchmarking

It contains node benchmarking scripts

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

### Other commands
```
./monitoring-tools roothash -start 1 -end 10
```
