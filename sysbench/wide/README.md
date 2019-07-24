# Wide Table benchmark

## Prepare

```bash
./wide_point_select.lua --table-size=1000 --db-driver=mysql --mysql-socket=/tmp/mysql.sock --mysql-user=root --columns=100 prepare 
```

## Run

```bash
./wide_point_select.lua --table-size=1000 --db-driver=mysql --mysql-socket=/tmp/mysql.sock --mysql-user=root --columns=100 run 
```