# Sysbench toolbox 

## Benchmark 

The script `./bench.sh` can help you run sysbench to benchmark different databases. Below are some examples:

```bash
# Load data into MySQL and TiDB
PORT=3306 ./bench.sh prepare mysql
PORT=4000 ./bench.sh prepare tidb

# Run some oltp workloads
PORT=3306 ./bench.sh run mysql oltp_read_write

# Save the result to another directory 
OUTPUT=./20190601 PORT=3306 ./bench.sh  run mysql oltp_update_index
```

## Reporter

The `sysbench-reporter` can help you generate some charts from the benchmark results. It can help you compare the performance of different databases, or check the performance changes of different versions for one database.

### Install

```bash
# Enter dbbench root directory

make 

# The sysbench-reporter is installed in the dbbench/bin directory
```

### Usage

Please notice that you must use the same workload to do the benchmark, and save the benchmark result to a unique directory for later comparision. For example, if you want to compare the performance of TiDB 2.1 and 3.0, you can use like below:

```bash
# For TiDB 2.1 
OUTPUT=./2.1 ./bench.sh run tidb oltp_point_select

# For TIDB 3.0
OUTPUT=./3.0 ./bench.sh run tidb oltp_point_select

# You must use different directories for benchmarking the same database, the sysbench-reporter will use 
# `tidb-3.0` and `tidb-2.1` (the name format is db-parentDir) to distinguish the results in the output charts. 
sysbench-reporter -p 2.1 -p 3.0 -o var 
```

If you want to compare different databases, you can save all data in one directory, e.g:

```bash
# Benchmark TiDB 
./bench.sh run tidb oltp_point_select
# Benchmark MySQL
./bench.sh run mysql oltp_point_select

# Passing -i here to tell reporter to use db name only (no need to include the parent directory) as identification in the chart, 
sysbench-reporter -p ./logs -o var -i
```