build:
	go build -o ./bin/migrate.exe

run: build
	./bin/migrate.exe

test:
	go test ./...