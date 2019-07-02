#!/bin/bash

TYPE=$1
DRIVER=$2
RUN_TYPE=$3

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

drop_mysql() {
    mysql -h ${HOST} -P ${PORT} -u ${DB_USER} -e "drop database if exists ${DB}"
    mysql -h ${HOST} -P ${PORT} -u ${DB_USER} -e "create database if not exists ${DB}"
}

drop_pgsql() {
    dropdb -h ${HOST} -p ${PORT} -U ${DB_USER} --if-exists -w ${DB}
    createdb -h ${HOST} -p ${PORT} -U ${DB_USER} -w ${DB}
}

if [ $DRIVER == "tidb" ]; then
    mysql -h ${HOST} -P ${PORT} -u ${DB_USER} -e "set global tidb_disable_txn_auto_retry = off"
fi


case ${TYPE} in
    prepare)
        case ${DB_DRIVER} in
            mysql)
            drop_mysql
            ;;
            pgsql)
            drop_pgsql
            ;;
        esac
        sysbench ${OPTS} oltp_point_select --tables=${TABLES} --table-size=${TABLE_SIZE} prepare
        ;;
    run)
        sysbench ${OPTS} ${RUN_TYPE} --tables=${TABLES} --table-size=${TABLE_SIZE} run 2>&1 | tee ${OUTPUT}/${DRIVER}_${RUN_TYPE}.log
        ;;
    *)
        echo "type must be prepare|run"
        exit 1
        ;;
esac
