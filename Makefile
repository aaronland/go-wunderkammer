cli:
	go build -mod vendor -o bin/append cmd/append/main.go
	go build -mod vendor -o bin/emit cmd/emit/main.go
	go build -mod vendor -o bin/wunderkammer-db cmd/wunderkammer-db/main.go
