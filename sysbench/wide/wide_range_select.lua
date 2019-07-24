#!/usr/bin/env sysbench

require("wide_common")

local get_table_num = get_table_num
local get_id = get_id
local get_col_num = get_col_num

function event()
	local tnum = get_table_num()
	local from_id = get_id()
	local to_id = from_id + sysbench.opt.range_size
	local cnum = get_col_num()
	local query = string.format([[
SELECT COUNT(c%d) FROM sbtest%d WHERE id BETWEEN %d and %d
]], cnum, tnum, from_id, to_id)

	con:query(query)
end