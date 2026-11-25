build:
	go build -o main main.go torrent.go

clean:
	rm -f main

test:
	go test -v ./...