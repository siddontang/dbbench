#!/usr/bin/env sysbench

require("wide_common")

local get_table_num = get_table_num
local get_id = get_id
local get_col_num = get_col_num

function event()
	local tnum = get_table_num()
	local id = get_id()
	local cnum = get_col_num()
	local query = string.format([[
SELECT c%d FROM sbtest%d WHERE id = %d
]], cnum, tnum, id)

	con:query(query)
end