#
# Ignore existing GOPATH, and force go to look for dependancies locally
#
GOPATH := $(shell pwd)

all:	foster

foster: 
	go build  -o foster ./src/foster.go

clean:
	rm foster
