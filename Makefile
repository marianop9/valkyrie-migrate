build:
	go build -o ./bin/valkyrie.exe ./cmd/valkyrie

test:
	go test ./...