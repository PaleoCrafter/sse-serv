
BIN       := ./sse-serv
COV_DIR   := coverage
ARTIFACTS := (\.prof|\.test)$$

help:               # Displays this list
	@cat Makefile \
	    | grep "^[a-z0-9_]\+:" \
	    | sed -r "s/:[^#]*?#?(.*)?/\r\t\t-\1/" \
	    | sed "s/^/ • make /"

clean:              # Removes build/test artifacts
	@rm -f ./$(BIN) &> /dev/null
	@rm -rf $(COV_DIR) &> /dev/null
	@find . -type f | grep -E "$(ARTIFACTS)" | xargs -I{} rm {};

build: clean        # Builds binary, sets README.md version
	@go build $(ARGS) && sed -i -r "1 s/\s-\s[^\s]+$$/ - `(./sse-serv -v)`/" README.md

test: clean mocks   # Runs tests with coverage
	@mkdir -p $(COV_DIR) &> /dev/null

	@go test $(shell go list ./... | grep -v mock) -coverprofile=$(COV_DIR)/coverage.out
	@go tool cover -html=$(COV_DIR)/coverage.out -o $(COV_DIR)/index.html && \
		echo "coverage: <file://$(PWD)/$(COV_DIR)/index.html>"

mocks:              # Creates mocks for tests
	@mkdir -p ./mock/amqp &> /dev/null
	@mkdir -p ./mock/logg &> /dev/null

	@mockgen sse-serv/amqp Provider,Consumer > ./mock/amqp/consumer.go
	@mockgen sse-serv/logg Logger            > ./mock/logg/logger.go

run: build          # Builds and executes binary
	@$(BIN)
