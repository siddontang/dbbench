#!/bin/bash

TYPE=$1
DB=$2

# go-ycsb path
CMD=${CMD:-go-ycsb}
# Output direcotry to save logs
OUTPUT=${OUTPUT:-./logs/}

WORKLOADS=(a b c d e)

RECORDCOUNT=1000000
OPERATIONCOUNT=1000000
THREADCOUNT=100
FIELDCOUNT=10
FIELDLENGTH=100
MAXSCANLENGTH=10

PROPS="-p recordcount=${RECORDCOUNT} \
    -p operationcount=${OPERATIONCOUNT} \
    -p threadcount=${THREADCOUNT} \
    -p fieldcount=${FIELDCOUNT} \
    -p fieldlength=${FIELDLENGTH} \
    -p maxscanlength=${MAXSCANLENGTH}"
PROPS+=" ${@:3}"

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

if [ ${TYPE} == 'load' ]; then 
    echo "clear data before load"
    PROPS+=" -p dropdata=true"
fi 

echo ${TYPE} ${DB} ${PROPS}

if [ ${TYPE} == 'load' ]; then 
    $CMD load ${DB} -p=workload=core ${PROPS} | tee ${OUTPUT}/${BENCH_DB}_load.log
elif [ ${TYPE} == 'run' ]; then
    for workload in ${WORKLOADS[@]}
    do 
        $CMD run ${DB} -P ./workloads/workload${workload} ${PROPS} | tee ${OUTPUT}/${BENCH_DB}_workload${workload}.log
    done
else
    echo "invalid type ${TYPE}"
    exit 1
fi 

