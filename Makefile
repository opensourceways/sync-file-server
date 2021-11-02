go ?= go
bazel ?= bazel-4.2.1
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
	$(bazel) run //:gazelle -- update-repos -from_file=go.mod

.PHONY: gen-bzl

gen-bzl:
	$(call tips,Generate bazel File)
	$(bazel) run //:gazelle

.PHONY: build

build: protocol load-lib gen-bzl
	$(call tips,Start Build)
	$(bazel) build //grpc/server:go_default_library
	$(bazel) build //grpc/client:go_default_library
	$(bazel) build //repo-file-cache:go_default_library
	$(bazel) build //server

.PHONY: clean

clean:
	$(bazel) clean
