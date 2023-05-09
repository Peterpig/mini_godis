all: linux	darwin

darwin:
	go build -o target/godis-darwin ./


linux:
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o target/godis-linux ./