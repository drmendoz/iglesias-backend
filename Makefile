GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=practical
    

all: build
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
	sudo nohup ./${BINARY_NAME} &

clean:
	sudo kill ${BINARY_NAME}
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
