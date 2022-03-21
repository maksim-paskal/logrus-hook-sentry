test:
	./scripts/validate-license.sh
	go fmt
	go mod tidy
	go test -race
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v
testScript:
	go run ./test-script
upgrade:
	go get -v -u all
	go mod tidy