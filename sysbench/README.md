# Sysbench toolbox 

## Benchmark 

The script `./bench.sh` can help you run sysbench to benchmark different databases. Below are some examples.

### OLTP

```bash
# Load data into MySQL and TiDB
PORT=3306 ./bench.sh mysql oltp_point_select prepare
PORT=4000 ./bench.sh tidb oltp_point_select prepare

# Run some oltp workloads
PORT=3306 ./bench.sh mysql oltp_read_write run

# Save the result to another directory 
OUTPUT=./20190601 PORT=3306 ./bench.sh mysql oltp_update_index run
```

### TPCC

```bash
PORT=3306 SCALE=100 ./bench.sh mysql tpcc prepare
PORT=3306 SCALE=100 ./bench.sh mysql tpcc run
```

### Blob

```bash
# TODO: support more blob commands soon
PORT=3306 BLOB_LENGTH=1000 ./bench.sh mysql blob_point_select prepare
PORT=3306 BLOB_LENGTH=1000 ./bench.sh mysql blob_update_non_index run
```