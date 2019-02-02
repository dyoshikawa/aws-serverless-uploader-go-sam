deps:
	dep ensure -v
	
build:
	GOOS=linux GOARCH=amd64 go build -o ./index/main ./index && \
	GOOS=linux GOARCH=amd64 go build -o ./store/main ./store && \
	GOOS=linux GOARCH=amd64 go build -o ./destroy/main ./destroy

test:
	go test -v ./index && \
	go test -v ./store
	