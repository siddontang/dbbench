uses longblob instead varchar for sbtest tables

executes only UPDATE sbtest%u SET c=? WHERE id=? statements

Prepare:

    ./blob_point_select.lua --table-size=1000 --db-driver=mysql --mysql-socket=/tmp/mysql.sock --mysql-user=root prepare 

Run:

    ./blob_point_select.lua --table-size=1000 --db-driver=mysql --mysql-socket=/tmp/mysql.sock --mysql-user=root --blob-length=578655 run
