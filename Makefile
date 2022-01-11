go ?= go
protoc ?= protoc

define tips
	$(info )
	$(info *************** $(1) ***************)
	$(info )
endef

.PHONY: protocol

protocol:
	$(call tips,Gen Protocol)
	$(protoc) -I ./proto --go_out=./protocol repo.proto sync_file.proto
	$(protoc) -I ./proto --go-grpc_out=./protocol repo.proto sync_file.proto

.PHONY: load-lib

load-lib:
	$(call tips,Load Lib)
	$(go) mod tidy

.PHONY: build

build: protocol load-lib
	$(call tips,Start Build)
	cd ./server && $(go) build

.PHONY: clean

clean:
	rm ./server/server
