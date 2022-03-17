linq:
	go build -o ./build/linq main.go

linq-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/linq main.go
