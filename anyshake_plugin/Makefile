BINDIR = .
GO_TARGET = $(BINDIR)/anyshake_plugin

GO = go
GOENV = CGO_ENABLED=0
GOFLAGS = -v -trimpath
GOLDFLAGS = -s -w

.PHONY: all go clean

all: $(GO_TARGET)

$(GO_TARGET):
	@mkdir -p $(BINDIR)
	$(GOENV) $(GO) build -ldflags="$(GOLDFLAGS)" $(GOFLAGS) -o $(GO_TARGET)

go: all

clean:
	rm -f $(GO_TARGET)
