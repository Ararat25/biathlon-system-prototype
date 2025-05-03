.PHONY: run
run:
	cd cmd && ./main

.PHONY: build
build:
	cd cmd && go build -o main

.PHONY: start
start: build run