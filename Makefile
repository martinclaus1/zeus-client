mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
	go get github.com/golang/mock/gomock
	go generate ./...
build:
	go mod tidy
	go build -v ./...
test: mocks build
	go clean -testcache
	go test -v ./... -count=1
release:
	chmod +x ./build.sh
	./build.sh $(version)
docs:
	go run main.go generate-docs
changelog:
	npm install -g auto-changelog
	auto-changelog