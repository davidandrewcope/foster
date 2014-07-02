#
# Ignore existing GOPATH, and force go to look for dependancies locally
#
GOPATH := $(shell pwd)
BIN_NAME = foster
BIN_PATH = .

# Common prefix for installation directories.
# NOTE: This directory must exist when you start the install.
prefix = /usr/local

INSTALL = install
INSTALL_PROGRAM = $(INSTALL)
INSTALL_PREFIX = $(prefix)


all:	compile

compile: 
	go build -a -o $(BIN_PATH)/$(BIN_NAME) ./src/foster.go
	
# Installs to the set path
.PHONY: install
install:	all
	@echo "Installing to $(DESTDIR)$(INSTALL_PREFIX)/bin"
	@$(INSTALL_PROGRAM) $(BIN_PATH)/$(BIN_NAME) $(DESTDIR)$(INSTALL_PREFIX)/bin

# Uninstalls the program
.PHONY: uninstall
uninstall:
	@echo "Removing $(DESTDIR)$(INSTALL_PREFIX)/bin/$(BIN_NAME)"
	@$(RM) $(DESTDIR)$(INSTALL_PREFIX)/bin/$(BIN_NAME)

.PHONY : clean
clean:
	-@$(RM) $(BIN_NAME)
