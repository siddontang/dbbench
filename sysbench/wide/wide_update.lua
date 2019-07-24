#!/usr/bin/env sysbench

require("wide_common")

local get_table_num = get_table_num
local get_id = get_id
local get_col_num = get_col_num
local get_c_value = get_c_value

function event()
	local tnum = get_table_num()
	local id = get_id()
	local cnum = get_col_num()
	local c_val = get_c_value()
	local query = string.format([[
UPDATE sbtest%d SET c%d= '%s' WHERE id = %d
]], tnum, cnum, c_val, id)

	con:query(query)
end