GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOTIDY=$(GOCMD) mod tidy 
BINARY_NAME=graphicalboredomkiller
INSTALL_DIR=/usr/local/bin

all: deps build

build: 
	$(GOBUILD) 

deps:
	$(GOTIDY)

run:
	$(GORUN) .

fmt:
	go fmt ./...

install: deps build
	sudo cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
