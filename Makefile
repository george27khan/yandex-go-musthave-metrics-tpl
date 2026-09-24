APP_AGENT=./bin/agent
APP_SERVER=./bin/server

.PHONY: build build-agent build-server run run-agent run-server clean

build: build-agent build-server

build-agent:
	mkdir -p bin
	go build -o $(APP_AGENT) ./cmd/agent

build-server:
	mkdir -p bin
	go build -o $(APP_SERVER) ./cmd/server

run: build
	$(APP_SERVER) -a=localhost:8080 &
	$(APP_AGENT) -a=localhost:8080 -p=2 -r=10

run-server: build-server
	$(APP_SERVER) -a=localhost:8080

run-agent: build-agent
	$(APP_AGENT) -a=localhost:8080 #-p=2 -r=10

clean:
	rm -rf bin