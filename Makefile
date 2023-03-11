mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
	go get github.com/golang/mock/gomock
	go generate ./...
test: mocks
	go clean -testcache
	go test -v ./... -count=1
build: test
	chmod +x ./build.sh
	./build.sh