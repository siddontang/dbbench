# sysbench-tpcc

TPCC-like workload for sysbench 1.0.x.
**Make sure you are using sysbench 1.0.14 or better!**

# prepare data and tables

`
./tpcc.lua --mysql-socket=/tmp/mysql.sock --mysql-user=root --mysql-db=sbt --time=300 --threads=64 --report-interval=1 --tables=10 --scale=100 --db-driver=mysql prepare
`

## prepare for RocksDB

`
./tpcc.lua --mysql-socket=/tmp/mysql.sock --mysql-user=root --mysql-db=sbr --time=3000 --threads=64 --report-interval=1 --tables=10 --scale=100 --use_fk=0 --mysql_storage_engine=rocksdb --mysql_table_options='COLLATE latin1_bin' --trx_level=RC --db-driver=mysql prepare
`

# Run benchmark

`
./tpcc.lua --mysql-socket=/tmp/mysql.sock --mysql-user=root --mysql-db=sbt --time=300 --threads=64 --report-interval=1 --tables=10 --scale=100 --db-driver=mysql run
`

## Run for cockroachdb

`
./tpcc.lua --pgsql-host=127.0.0.1 --pgsql-port=26257 --pgsql-user=root --pgsql-db=sbt --time=300 --threads=64 --report-interval=1 --tables=10 --scale=100 --db-driver=pgsql --for-update="" --restart-rollback=0 run
`

# Cleanup 

`
./tpcc.lua --mysql-socket=/tmp/mysql.sock --mysql-user=root --mysql-db=sbt --time=300 --threads=64 --report-interval=1 --tables=10 --scale=100 --db-driver=mysql cleanup
`
