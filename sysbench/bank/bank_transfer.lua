#!/usr/bin/env sysbench

if sysbench.cmdline.command == nil then
    error("Command is required. Supported commands: prepare, run, cleanup, help")
end

sysbench.cmdline.options = {
    table_size = {"Number of rows per table", 10000},
    tables = {"Number of tables", 1}
}

function cmd_prepare()
    local drv = sysbench.sql.driver()
    local con = drv:connect()

    for i = sysbench.tid % sysbench.opt.threads + 1, sysbench.opt.tables, sysbench.opt
        .threads do 
        create_table(drv, con, i) 
    end
end

sysbench.cmdline.commands = {
    prepare = {cmd_prepare, sysbench.cmdline.PARALLEL_COMMAND}
}

function create_table(drv, con, table_num)
    print(string.format("Creating table 'account%d'...", table_num))

    local query = string.format([[
CREATE TABLE account%d(
  id INTEGER NOT NULL,
  balance INTEGER DEFAULT '1000' NOT NULL,
  PRIMARY KEY (id)
)]], table_num)

    con:query(query)

    if (sysbench.opt.table_size > 0) then
        print(string.format("Inserting %d records into 'account%d'",
                            sysbench.opt.table_size, table_num))
    end

    query = "INSERT INTO account" .. table_num .. "(id, balance) VALUES"

    con:bulk_insert_init(query)

    for i = 1, sysbench.opt.table_size do
        query = string.format("(%d, %d)", i, 1000)

        con:bulk_insert_next(query)
    end

    con:bulk_insert_done()
end

function cleanup()
    local drv = sysbench.sql.driver()
    local con = drv:connect()

    for i = 1, sysbench.opt.tables do
        print(string.format("Dropping table 'account%d'...", i))
        con:query("DROP TABLE IF EXISTS account" .. i)
    end
end

function thread_init()
    drv = sysbench.sql.driver()
    con = drv:connect()
end

function thread_done() 
    con:disconnect() 
end

function sysbench.hooks.before_restart_event(err)
    con:query("ROLLBACK")
end

local function get_table_num()
   return sysbench.rand.uniform(1, sysbench.opt.tables)
end

local function get_id()
   return sysbench.rand.default(1, sysbench.opt.table_size)
end

function event() 
    local from = get_id()
    local to = get_id()
    local table_num = get_table_num()
    local amount = sysbench.rand.default(1, 100)
    while(from == to)
    do
        to = get_id()
    end

    con:query("BEGIN")

    local rs = con:query(string.format([[
SELECT id, balance FROM account%d WHERE id IN (%d, %d) FOR UPDATE
]], table_num, from, to))

    assert(rs.nrows == 2)

    local row_from = rs:fetch_row()
    local row_to = rs:fetch_row()

    if row_from[1] ~= from then
        row_from, row_to = row_to, row_from
    end 

    if row_from[2] - amount < 0 then 
        con:query("ROLLBACK")
        return 
    end

    con:query(string.format([[
UPDATE account%d SET balance = balance - %d WHERE id = %d
]], table_num, amount, from))

    con:query(string.format([[
UPDATE account%d SET balance = balance + %d WHERE id = %d
]], table_num, amount, to))

    con:query("COMMIT") 
end
