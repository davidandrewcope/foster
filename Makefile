#
# Ignore existing GOPATH, and force go to look for dependancies locally
#
GOPATH := $(shell pwd)

all:	foster

foster: 
	go build -a -o foster ./src/foster.go

.PHONY : clean
clean:
	rm foster
