cli:
	go build -mod vendor -o bin/append-dataurl cmd/append-dataurl/main.go
	go build -mod vendor -o bin/wunderkammer-db cmd/wunderkammer-db/main.go
