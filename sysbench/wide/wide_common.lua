#!/usr/bin/env sysbench
-- Copyright (C) 2006-2018 Alexey Kopytov <akopytov@gmail.com>

-- This program is free software; you can redistribute it and/or modify
-- it under the terms of the GNU General Public License as published by
-- the Free Software Foundation; either version 2 of the License, or
-- (at your option) any later version.

-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU General Public License for more details.

-- You should have received a copy of the GNU General Public License
-- along with this program; if not, write to the Free Software
-- Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA

-- -----------------------------------------------------------------------------
-- Common code for Wide table benchmarks.
-- -----------------------------------------------------------------------------

function init()
   assert(event ~= nil,
          "this script is meant to be included by other OLTP scripts and " ..
             "should not be called directly.")
end

if sysbench.cmdline.command == nil then
   error("Command is required. Supported commands: prepare, warmup, run, " ..
            "cleanup, help")
end

-- Command line options
sysbench.cmdline.options = {
   table_size =
      {"Number of rows per table", 10000},
   range_size =
      {"Range size for range SELECT queries", 100},
   columns =
      {"Number of columns per table", 100},
   tables =
      {"Number of tables", 1},
}

-- Prepare the dataset. This command supports parallel execution, i.e. will
-- benefit from executing with --threads > 1 as long as --tables > 1
function cmd_prepare()
   local drv = sysbench.sql.driver()
   local con = drv:connect()

   for i = sysbench.tid % sysbench.opt.threads + 1, sysbench.opt.tables,
   sysbench.opt.threads do
      create_table(drv, con, i)
   end
end

-- Preload the dataset into the server cache. This command supports parallel
-- execution, i.e. will benefit from executing with --threads > 1 as long as
-- --tables > 1
--
-- PS. Currently, this command is only meaningful for MySQL/InnoDB benchmarks
function cmd_warmup()
   local drv = sysbench.sql.driver()
   local con = drv:connect()

   assert(drv:name() == "mysql", "warmup is currently MySQL only")

   for i = sysbench.tid % sysbench.opt.threads + 1, sysbench.opt.tables,
   sysbench.opt.threads do
      local t = "sbtest" .. i
      print("Preloading table " .. t)
      con:query("ANALYZE TABLE sbtest" .. i)
   end
end

-- Implement parallel prepare and warmup commands, define 'prewarm' as an alias
-- for 'warmup'
sysbench.cmdline.commands = {
   prepare = {cmd_prepare, sysbench.cmdline.PARALLEL_COMMAND},
   warmup = {cmd_warmup, sysbench.cmdline.PARALLEL_COMMAND},
   prewarm = {cmd_warmup, sysbench.cmdline.PARALLEL_COMMAND}
}


-- Template strings of random digits with 11-digit groups separated by dashes

-- 10 groups, 119 characters
local c_value_template = "###########-###########-###########-" ..
   "###########-###########-###########-" ..
   "###########-###########-###########-" ..
   "###########"


function get_c_value()
   return sysbench.rand.string(c_value_template)
end

function create_table(drv, con, table_num)
   local query
   local cols = {}

   print(string.format("Creating table 'sbtest%d'...", table_num))

   for i = 1, sysbench.opt.columns do
      cols[i] = string.format("c%d", i)
   end

   query = string.format([[
CREATE TABLE sbtest%d(
  id INTEGER NOT NULL,
  %s VARCHAR(120) NOT NULL,
  PRIMARY KEY (id)
)]],
      table_num, table.concat(cols, " VARCHAR(120) NOT NULL,\n"))

   con:query(query)


   if (sysbench.opt.table_size > 0) then
      print(string.format("Inserting %d records into 'sbtest%d'",
                          sysbench.opt.table_size, table_num))
   end

   query = string.format([[
INSERT INTO sbtest%d (id, %s) VALUES
]], table_num, table.concat(cols, ","))
   
   con:bulk_insert_init(query)

   local c_val
   local c_vals = {}

   for i = 1, sysbench.opt.table_size do

      c_val = get_c_value()

      for i = 1, sysbench.opt.columns do
         c_vals[i] = string.format("'%s'", c_val)
      end

      query = string.format("(%d, %s)",
                            i,           
                            table.concat(c_vals, ","))

      con:bulk_insert_next(query)
   end

   con:bulk_insert_done()
end


function thread_init()
   drv = sysbench.sql.driver()
   con = drv:connect()
end

function thread_done()
   con:disconnect()
end

function cleanup()
   local drv = sysbench.sql.driver()
   local con = drv:connect()

   for i = 1, sysbench.opt.tables do
      print(string.format("Dropping table 'sbtest%d'...", i))
      con:query("DROP TABLE IF EXISTS sbtest" .. i )
   end
end

function get_table_num()
   return sysbench.rand.uniform(1, sysbench.opt.tables)
end

function get_id()
   return sysbench.rand.default(1, sysbench.opt.table_size)
end

function get_col_num()
   return sysbench.rand.uniform(1, sysbench.opt.columns)
end

function sysbench.hooks.before_restart_event(errdesc)
   con:query("ROLLBACK")
end
