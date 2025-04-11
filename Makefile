# Final binary name
BINARY = redraft

# Default config path (can be overridden)
CONFIG ?= config/node1.yaml

# Build the binary
build:
	go build -o $(BINARY) main.go

# Run node dynamically (pass CONFIG=... or use default)
run: build
	./$(BINARY) $(CONFIG)

# Named node shortcuts (these use static config files)
node1: CONFIG=config/node1.yaml
node1: run

node2: CONFIG=config/node2.yaml
node2: run

node3: CONFIG=config/node3.yaml
node3: run

# Optional: Start cluster in tmux panes (3 windows)
cluster: build
	@echo "Starting 3-node cluster in tmux"
	@tmux new-session -d -s redraft 'make node1'
	@tmux split-window -h 'make node2'
	@tmux split-window -v 'make node3'
	@tmux select-layout even-horizontal
	@tmux attach-session -t redraft

# Run all tests (including memory and bolt)
test:
	go test -v ./test/...

# Clean up build artifacts and raft data
clean:
	rm -f $(BINARY)
	rm -rf data/
	rm -rf testdata/

.PHONY: build run node1 node2 node3 cluster test clean
