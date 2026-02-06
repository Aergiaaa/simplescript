.PHONY: build run clean test install

BIN_NAME=idiot
GO=go

# check os
ifeq ($(OS),Windows_NT)
	BIN_EXT=.exe
	PATH_INSTALL=$(USERPROFILE)/bin
	MKDIR=if not exist $(subst /,\,$(PATH_INSTALL)) mkdir $(subst /,\,$(PATH))
	CP=copy
	RM=del /Q
else
	BIN_EXT=
	PATH_INSTALL=/usr/local/bin
	MKDIR=mkdir -p $(PATH_INSTALL)
	CP=cp
	RM=rm -f
endif

BIN=$(BIN_NAME)$(BIN_EXT)

build:
	@$(GO) build -o $(BIN)

run:
	@$(GO) run main.go

test:
	@$(GO) test ./...

clean:
	@$(GO) clean
	$(RM) $(BIN_NAME)

install: build
	$(MKDIR)
ifeq ($(OS),Windows_NT)
	$(CP) $(BIN) $(PATH_INSTALL)/$(BIN)
	@echo idiotic interpreter has been Installed
	@echo add $(PATH_INSTALL) to PATH if not already present
else
	sudo $(CP) $(BIN) $(PATH_INSTALL)/$(BIN)
	@echo "idiotic interpreter has been Installed"
endif

uninstall:
ifeq ($(OS),Windows_NT)
	$(RM) $(PATH_INSTALL)/$(BIN)
	@echo idiotic interpreter has been removed
else
	sudo $(RM) $(PATH_INSTALL)/$(BIN)
	@echo "idiotic interpreter has been removed"
endif
