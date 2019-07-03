# YCSB toolbox 

## Benchmark 

The script `./bench.sh` can help you run [go-ycsb](https://github.com/pingcap/go-ycsb) to benchmark different databases. Below are some examples:

```bash
# Load data into TiKV with Raw mode
./bench.sh load raw -p tikv.pd=127.0.0.1:2379

# Run workloads a, b, ... e for TiKV with Raw mode
./bench.sh run raw -p tikv.pd=127.0.0.1:2379

# Run benchmark and output the result to the OUTPUT directory
OUTPUT=./20190601 ./bench.sh load raw
```