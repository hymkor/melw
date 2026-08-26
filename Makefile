ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
    NUL=nul
    DEL=del
else
    SET=export
    NUL=/dev/null
    DEL=rm
endif

NAME:=$(subst go-,,$(notdir $(CURDIR)))
VERSION:=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT:=-ldflags "-s -w -X main.version=$(VERSION)"
EXE:=$(shell go env GOEXE)

all: $(wildcard *.go) github.css
	go fmt
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)

_dist:
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)
	zip -9 $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(NAME)$(EXE)

dist:
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _dist

bump:
	$(GO) run github.com/hymkor/latest-notes@latest -suffix "-goinstall" -gosrc main CHANGELOG*.md > version.go

github.css :
	curl https://raw.githubusercontent.com/sindresorhus/github-markdown-css/gh-pages/github-markdown-light.css > github.css

release:
	$(GO) run github.com/hymkor/latest-notes@latest | gh release create -d --notes-file - -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

manifest:
	$(GO) run github.com/hymkor/make-scoop-manifest@latest -all *-windows-*.zip > $(NAME).json

.PHONY: all test dist _dist clean manifest release
