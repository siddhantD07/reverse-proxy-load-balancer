GOPATH := ${PWD}:${GOPATH}
export GOPATH

default: build

build:
	# hint: "go get github.com/go-zookeeper/zk"
	go build -tags netgo -v -o ./bin/grproxy ./src/grproxy
