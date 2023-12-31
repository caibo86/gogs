.PHONY: .FORCE

GO?=go
PWD=$(shell pwd)
GOPATH:=$(PWD)
GOBIN:=$(GOPATH)/bin
PATH:=$(PATH):$(GOBIN)


BUILD_VERSION?=undefined
BUILD_REVISION?=$(shell ./src/scripts/gen_revision.sh)
BUILD_TIME=$(shell TZ='Asia/Shanghai' date '+%FT%T')

BOILERPLATE=$(PWD)/src/cmd/gengo/boilerplate/boilerplate.go.txt

GO_TAGS=$(GO_BUILD_TAGS)

LD_FLAGS="\
-X $(VER_PKG).version=$(BUILD_VERSION) \
-X $(VER_PKG).revision=$(BUILD_REVISION) \
-X $(VER_PKG).buildTime=$(BUILD_TIME) "

GO_FLAGS=-ldflags=$(LD_FLAGS) -tags=$(GO_TAGS)

cbc:
	cd src && GOBIN=$(GOBIN) $(GO) install $(GO_FLAGS) gogs/apps/cbc/cb2go

t:
	@rm -f src/base/cblang/*.cb.go
	@cd src && GOBIN=$(GOBIN) $(GO) install $(GO_FLAGS) gogs/apps/cbc/cb2go
	GOPATH=$(GOPATH) cb2go --module gogs gs cb

getdeepcopy:
	cd src/cmd/gengo/examples/deepcopy-gen && go build . && cp deepcopy-gen $(GOPATH)/bin

deepcopy:
	cd src && GOPATH=/home/cb/go/src/gogs deepcopy-gen -i gogs/gs -v=5 --trim-path-prefix /home/cb/go/src/gogs/src/gogs/ --logtostderr -h $(BOILERPLATE)

dev:
	cd ./src && GOPATH=$(GOPATH) $(GO) install -v github.com/golang/protobuf/protoc-gen-go@v1.5.2
	cd ./src && GOPATH=$(GOPATH) $(GO) install -v github.com/gogo/protobuf/protoc-gen-gofast@v1.3.1

g:
	@protoc --proto_path=src/pb --gofast_out=./src/pb test.proto

f:
	clang-format -i src/gs/test.gs

mod:
	cd ./src && go mod tidy && go mod vendor

gate:
	cd src && GOBIN=$(GOBIN) $(GO) install $(GO_FLAGS) gogs/apps/gate_main

game:
	cd src && GOBIN=$(GOBIN) $(GO) install $(GO_FLAGS) gogs/apps/game_main

login:
	cd src && GOBIN=$(GOBIN) $(GO) install $(GO_FLAGS) gogs/apps/login_main

simulator:
	cd src && GOBIN=$(GOBIN) $(GO) install $(GO_FLAGS) gogs/apps/simulator_main

etcd:
	sh ./src/scripts/init_config_in_etcd.sh