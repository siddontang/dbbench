# Sysbench toolbox 

## Benchmark 

The script `./bench.sh` can help you run sysbench to benchmark different databases. Below are some examples.

### OLTP

```bash
# Load data into MySQL and TiDB
PORT=3306 ./bench.sh prepare mysql
PORT=4000 ./bench.sh prepare tidb

# Run some oltp workloads
PORT=3306 ./bench.sh run mysql oltp_read_write

# Save the result to another directory 
OUTPUT=./20190601 PORT=3306 ./bench.sh  run mysql oltp_update_index
```

### TPCC

```bash
PORT=3306 SCALE=100 ./bench.sh prepare mysql tpcc 
PORT=3306 SCALE=100 ./bench.sh run mysql tpcc 
```

### Blob

```bash
# TODO: support more blob commands soon
PORT=3306 BLOB_LENGTH=1000 ./bench.sh prepare mysql blob 
PORT=3306 BLOB_LENGTH=1000 ./bench.sh run mysql blob 
```