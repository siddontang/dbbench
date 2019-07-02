# Sysbench toolbox 

## Benchmark 

The script `./bench.ch` can help you run sysbench to benchmark different databases. Below are some examples:

```bash
# Load data into MySQL and TiDB
PORT=3306 ./bench.sh prepare mysql
PORT=4000 ./bench.sh prepare tidb

# Run some oltp workloads
PORT=3306 ./bench.sh run mysql oltp_read_write

# Save the result to another directory 
OUTPUT=./20190601 PORT=3306 ./bench.sh  run mysql oltp_update_index
```