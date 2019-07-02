default: build

build:
	go build -o ./bin/ycsb_reporter ./ycsb/reporter/* 
	go build -o ./bin/sysbench_reporter ./sysbench/reporter/* 
