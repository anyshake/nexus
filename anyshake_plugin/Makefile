BINDIR = .
SRC = ../../libs/plugin/plugin.c
OBJ = plugin.o

CLIB = libplugin.a
GO_TARGET = $(BINDIR)/anyshake_plugin

CC = gcc
GOENV = CGO_ENABLED=1
CFLAGS = -Wall -fPIC -I../
LDFLAGS =

GO = go
GOENV = 
GOFLAGS = -v -trimpath
GOLDFLAGS = -s -w

.PHONY: all lib go clean

all: $(GO_TARGET)

lib: $(CLIB)

$(CLIB): $(OBJ)
	ar rcs $@ $^

$(OBJ): $(SRC)
	$(CC) $(CFLAGS) -c $< -o $@

go: lib
	@mkdir -p $(BINDIR)
	$(GOENV) $(GO) build -ldflags="$(GOLDFLAGS)" $(GOFLAGS) -o $(GO_TARGET)

$(GO_TARGET): go

clean:
	rm -f $(CLIB) $(OBJ) $(GO_TARGET)
