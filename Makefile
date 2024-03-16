build:
	go build -o ./bin/migrate.exe ./cmd/valkyrie-migrate

run: build
	./bin/migrate.exe

test:
	go test ./...