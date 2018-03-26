MAKEFILEDIR = $(shell pwd)

GO      = go
TARGET  = photobooth
GOLINT  = $(GOPATH)/bin/gometalinter

MKDIR_P = mkdir -p

GLIDE_VERSION := $(shell glide --version 2>/dev/null)
DEP_VERSION := $(shell dep version 2>/dev/null)

GPHOTO := $(shell command -v gphoto2 2> /dev/null)

include config

all: $(TARGET) run

$(TARGET): vendor clean
	$(GO) build -ldflags="-s -w" -o $@ ./cmd/main.go

vendor:
ifdef DEP_VERSION
	dep ensure
else ifdef GLIDE_VERSION
	glide install
else
	go get ./cmd/main.go
endif

clean:
	rm -f $(TARGET)

test: vendor lint

run: directories
ifndef GPHOTO
	$(error "gphoto2 is required.")
endif
	./$(TARGET) --port=$(APPLICATION_PORT) --limit=$(IMAGE_LIMIT) --photopath=$(PHOTO_DIR) --thumbpath=$(THUMB_DIR)

$(GOPATH)/bin/goconst:
	$(GO) get github.com/jgautheron/goconst/cmd/goconst

directories:
	$(MKDIR_P) $(PHOTO_DIR)
	cp ./static/default/default.png $(PHOTO_DIR)
	$(MKDIR_P) $(THUMB_DIR)
	cp ./static/default/thumb/default.png $(THUMB_DIR)

lint: \
	$(GOLINT) \
	$(GOPATH)/bin/goconst \
	$(GOPATH)/bin/ineffassign \
	$(GOPATH)/bin/varcheck \
	$(GOPATH)/bin/structcheck \
	$(GOPATH)/bin/aligncheck \
	$(GOPATH)/bin/gocyclo \
	$(GOPATH)/bin/interfacer \
	$(GOPATH)/bin/gosimple \
	$(GOPATH)/bin/deadcode \
	$(GOPATH)/bin/unconvert \
	$(GOPATH)/bin/staticcheck \
	$(GOPATH)/bin/gas
		$(GOLINT) --deadline 30s ./cmd/

$(GOPATH)/bin/ineffassign:
	$(GO) get github.com/gordonklaus/ineffassign

$(GOLINT):
	$(GO) get -u github.com/alecthomas/gometalinter

$(GOPATH)/bin/aligncheck:
	$(GO) get github.com/opennota/check/cmd/aligncheck

$(GOPATH)/bin/structcheck:
	$(GO) get github.com/opennota/check/cmd/structcheck

$(GOPATH)/bin/varcheck:
	$(GO) get github.com/opennota/check/cmd/varcheck

$(GOPATH)/bin/gocyclo:
	$(GO) get github.com/fzipp/gocyclo

$(GOPATH)/bin/interfacer:
	$(GO) get mvdan.cc/interfacer

$(GOPATH)/bin/gosimple:
	$(GO) get honnef.co/go/tools/cmd/gosimple

$(GOPATH)/bin/deadcode:
	$(GO) get github.com/tsenart/deadcode

$(GOPATH)/bin/unconvert:
	$(GO) get github.com/mdempsky/unconvert

$(GOPATH)/bin/staticcheck:
	$(GO) get honnef.co/go/tools/cmd/staticcheck

$(GOPATH)/bin/gas:
	$(GO) get github.com/GoASTScanner/gas

.PHONY:test lint vendor run