default: build

build:
	go build -o ./bin/ycsb_reporter ./ycsb/reporter/* 
