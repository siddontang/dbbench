#!/bin/bash

DB=$1
TYPE=$2

# go-ycsb path
CMD=${CMD:-go-ycsb}
# Output direcotry to save logs
OUTPUT=${OUTPUT:-./logs/}

RECORDCOUNT=${RECORDCOUNT:-1000000}
OPERATIONCOUNT=${OPERATIONCOUNT:-1000000}
THREADCOUNT=${THREADCOUNT:-100}
FIELDCOUNT=${FIELDCOUNT:-10}
FIELDLENGTH=${FIELDLENGTH:-100}
MAXSCANLENGTH=${MAXSCANLENGTH:-10}

PROPS="-p recordcount=${RECORDCOUNT} \
    -p operationcount=${OPERATIONCOUNT} \
    -p threadcount=${THREADCOUNT} \
    -p fieldcount=${FIELDCOUNT} \
    -p fieldlength=${FIELDLENGTH} \
    -p maxscanlength=${MAXSCANLENGTH}"

mkdir -p ${OUTPUT} 

BENCH_DB=${DB}

case ${DB} in
    mysql)
        ;;
    mysql8)
        DB="mysql"
        ;;
    mariadb)
        DB="mysql"
        ;;
    pg)
        ;;
    tikv)
        PROPS+=" -p tikv.type=txn"
        ;;
    raw)
        PROPS+=" -p tikv.type=raw"
        DB="tikv"
        ;;
    tidb)
        PROPS+=" -p mysql.port=4000"
        ;;
    cockroach)
        PROPS+=" -p pg.port=26257"
        ;;
    cassandra)
        ;;
    scylla)
        ;;
    *)
    ;;
esac

PROPS+=" ${@:3}"

if [ ${TYPE} == 'load' ]; then 
    echo "clear data before load"
    PROPS+=" -p dropdata=true"
fi 

echo ${TYPE} ${DB} ${PROPS}

if [ ${TYPE} == 'load' ]; then 
    $CMD load ${DB} -p=workload=core ${PROPS} | tee ${OUTPUT}/${BENCH_DB}_load.log
else
    $CMD run ${DB} -P ./workloads/${TYPE} ${PROPS} | tee ${OUTPUT}/${BENCH_DB}_${TYPE}.log
fi 

