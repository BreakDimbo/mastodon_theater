export GOPATH := $(CURDIR)
export LD_LIBRARY_PATH=/Users/break/Documents/Geek/cmx_bot/lib/libmsc.so

theater:
	@echo "Building Theater ..."
	go build -o bin/theater theater

deps:
	@echo "Install Installing dependencies"
	@go get -u github.com/golang/dep/cmd/dep
	cd src/theater; ${GOPATH}/bin/dep init; ${GOPATH}/bin/dep ensure -v