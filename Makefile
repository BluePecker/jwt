APP=jwt
build:
		go build -o ${GOPATH}/bin/${APP} -ldflags '-s -w' ./main.go
linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${GOPATH}/bin/${APP} -ldflags '-s -w' ./main.go
run:
		go run *.go
clean:
		@rm -rf ${GOPATH}/bin/${APP}
