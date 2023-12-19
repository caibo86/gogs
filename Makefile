.PHONY: .FORCE

GO?=go
PWD=$(shell pwd)
GOPATH:=$(PWD)
GOBIN:=$(GOPATH)/bin
PATH:=$(PATH):$(GOBIN)


BUILD_VERSION?=undefined
BUILD_REVISION?=$(shell ./src/scripts/gen_revision.sh)
BUILD_TIME=$(shell TZ='Asia/Shanghai' date '+%FT%T')

GO_TAGS=$(GO_BUILD_TAGS)

LD_FLAGS="\
-X $(VER_PKG).version=$(BUILD_VERSION) \
-X $(VER_PKG).revision=$(BUILD_REVISION) \
-X $(VER_PKG).buildTime=$(BUILD_TIME) "

GO_FLAGS=-ldflags=$(LD_FLAGS) -tags=$(GO_TAGS)

gsc:
	cd src && GOBIN=$(GOBIN) $(GO) install $(GO_FLAGS) gogs/apps/gsc/gs2go