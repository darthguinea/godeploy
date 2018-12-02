export APP=$(shell basename $(CURDIR))
export GOPATH=${PWD}
export GOBIN=${PWD}
export INSTALL_PATH=${HOME}/Programming/Scripts

all:
	go get
	go build

install:
	cp ${APP} ${INSTALL_PATH}

clean:
	go clean

clean-all:
	go clean
	rm -rvf src/github.com
