fusp: $(wildcard */*/*.go)
	go build -o bin/fusp ./cmd/fusp/main.go