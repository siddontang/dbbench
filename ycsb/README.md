# YCSB toolbox 

## Benchmark 

The script `./bench.sh` can help you run [go-ycsb](https://github.com/pingcap/go-ycsb) to benchmark different databases. Below are some examples:

```bash
# Load data into TiKV with Raw mode
./bench.sh raw load -p tikv.pd=127.0.0.1:2379

# Run workloads a for TiKV with Raw mode
./bench.sh raw workloada -p tikv.pd=127.0.0.1:2379

# Run benchmark and output the result to the OUTPUT directory
OUTPUT=./20190601 ./bench.sh raw load
```