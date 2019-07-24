#!/bin/bash

DRIVER=$1
RUN_TYPE=$2
TYPE=$3

ROOT=$(pwd)

EXT_LIBS=("tpcc" "blob" "bank" "wide")

LUA_PATH=""
for name in ${EXT_LIBS[@]}; do 
    LUA_PATH+="$ROOT/$name/?.lua;"
done 

# add default libs
LUA_PATH+=";"
export LUA_PATH

# Output direcotry to save logs
OUTPUT=${OUTPUT:-./logs/}

mkdir -p ${OUTPUT}

THREADS=${THREADS:-32}
TABLES=${TABLES:-1}
TABLE_SIZE=${TABLE_SIZE:-1000000}

HOST=${HOST:-127.0.0.1}
PORT=${PORT:-4000}
DB_USER=${DB_USER:-root}
DB=${DB:-sbtest}

RNAD_TYPE=${RNAD_TYPE:-uniform}

# Used in prepare
DROPDATA=1

TIME=${TIME:-600}
REPORT_INTERVAL=${REPORT_INTERVAL:-10}

OPTS="--report-interval=${REPORT_INTERVAL} \
    --time=${TIME} \
    --rand-type=${RNAD_TYPE} \
    --threads=${THREADS} "

DB_DRIVER=mysql 

# For blob
BLOB_LENGTH=${BLOB_LENGTH:-10240}

# For tpcc
SCALE=${SCALE:-100}
USE_FK=${USE_FK:-1}

# For wide
COLUMNS=${COLUMNS:-100}

COMMAND_OPTS=""

case ${RUN_TYPE} in
    tpcc)
        COMMAND_OPTS=" --tables=${TABLES} \
        --scale=${SCALE} \
        --use-fk=${USE_FK} "
    ;;
    blob_*)
        COMMAND_OPTS=" --tables=${TABLES} \
        --table-size=${TABLE_SIZE} \
        --blob-length=${BLOB_LENGTH} "
    ;;
    wide_*)
        COMMAND_OPTS=" --tables=${TABLES} \
        --table-size=${TABLE_SIZE} \
        --columns=${COLUMNS}"
    ;;
    *)
        COMMAND_OPTS=" --tables=${TABLES} \
        --table-size=${TABLE_SIZE} "
    ;;
esac

COMMAND_OPTS+=" ${@:4}"

case ${DRIVER} in
    mysql|tidb)
        DB_DRIVER=mysql
        OPTS+="--mysql-host=${HOST} \
            --mysql-port=${PORT} \
            --mysql-user=${DB_USER} \
            --mysql-db=${DB} \
            --db-driver=mysql"
        ;;
    pg|cockroachdb)
        DB_DRIVER=pgsql
        OPTS+="--pgsql-host=${HOST} \
            --pgsql-port=${PORT} \
            --pgsql-user=${DB_USER} \
            --pgsql-db=${DB} \
            --db-driver=pgsql"
        ;;
    *)
        ;;
esac

if [ $DRIVER == "tidb" ]; then
    mysql -h ${HOST} -P ${PORT} -u ${DB_USER} -e "set global tidb_disable_txn_auto_retry = off"

    if [ $RUN_TYPE == "tpcc" ]; then 
        mysql -h ${HOST} -P ${PORT} -u ${DB_USER} -e "set global sql_mode = ''"
    fi 

    OPTS+=" --mysql-ignore-errors=8002"
fi

case ${TYPE} in
    prepare)
        case ${DB_DRIVER} in
            mysql)
            mysql -h ${HOST} -P ${PORT} -u ${DB_USER} -e "create database if not exists ${DB}"
            ;;
            pgsql)
            createdb -h ${HOST} -p ${PORT} -U ${DB_USER} -w ${DB}
            echo $?
            ;;
        esac
        sysbench ${OPTS} ${RUN_TYPE} ${COMMAND_OPTS} cleanup
        sysbench ${OPTS} ${RUN_TYPE} ${COMMAND_OPTS} prepare
        ;;
    run)
        sysbench ${OPTS} ${RUN_TYPE} ${COMMAND_OPTS} run 2>&1 | tee ${OUTPUT}/${DRIVER}_${RUN_TYPE}.log
        ;;
    cleanup)
        sysbench ${OPTS} ${RUN_TYPE} ${COMMAND_OPTS} cleanup
        ;;
    *)
        echo "type must be prepare|run|cleanup"
        exit 1
        ;;
esac
