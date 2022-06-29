.PHONY: all build compile clean bin/setup-linux-amd64-$(VERSION)-$(BUILDNUMER)

BUILDTIME ?= $(shell date +%Y-%m-%d_%I:%M:%S)
GITCOMMIT ?= $(shell git rev-parse -q HEAD)
ifeq ($(CI_PIPELINE_ID),)
	BUILDNUMER := private
else
	BUILDNUMER := $(CI_PIPELINE_ID)
endif
VERSION ?= $(shell git describe --tags --always --dirty)

LDFLAGS = -extldflags \
		  -static \
		  -X "main.Version=$(VERSION)" \
		  -X "main.BuildTime=$(BUILDTIME)" \
		  -X "main.GitCommit=$(GITCOMMIT)" \
		  -X "main.BuildNumber=$(BUILDNUMER)"

all: build

clean:
	rm -rf bin setup-*.zip

build:
	go build -o bin/setup -ldflags "$(LDFLAGS)"

bin/setup-linux-amd64-$(VERSION)-$(BUILDNUMER):
	GOOS=linux GOARCH=amd64 go build -o bin/setup-linux-amd64-$(VERSION)-$(BUILDNUMER) -ldflags "$(LDFLAGS)"

package: bin/setup-linux-amd64-$(VERSION)-$(BUILDNUMER)
	echo $(GITCOMMIT) > commit.txt
	echo $(VERSION) > version.txt
	zip -r setup-$(VERSION)-$(BUILDNUMER).zip bin commit.txt version.txt

compile: bin/setup-linux-amd64-$(VERSION)-$(BUILDNUMER)

