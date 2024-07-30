
tidy:
	go mod tidy -v

## test: 🚦 Execute all tests
test:
	go run gotest.tools/gotestsum@latest -f testname -- ./... -race -count=1 -shuffle=on
