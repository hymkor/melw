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
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _dist
	$(SET) "GOOS=darwin"  && $(SET) "GOARCH=amd64" && $(MAKE) _dist
	$(SET) "GOOS=darwin"  && $(SET) "GOARCH=arm64" && $(MAKE) _dist
	$(SET) "GOOS=freebsd" && $(SET) "GOARCH=amd64" && $(MAKE) _dist
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _dist

bump:
	$(GO) run github.com/hymkor/latest-notes@latest -suffix "-goinstall" -gosrc main CHANGELOG*.md > version.go

release:
	$(GO) run github.com/hymkor/latest-notes@latest | gh release create -d --notes-file - -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

manifest:
	$(GO) run github.com/hymkor/make-scoop-manifest@latest -all *-windows-*.zip > $(NAME).json

# MINIPAGE=$(GO) run github.com/hymkor/minipage@latest 
MINIPAGE=minipage

docs:
	$(MINIPAGE) -title "Melw - Markdown Editor Like WikiEngine" -outline-in-sidebar -readme-to-index README.md > docs/index.html
	$(MINIPAGE) -title "Melw - Markdown Editor Like WikiEngine" -outline-in-sidebar -readme-to-index README_ja.md > docs/index_ja.html
	$(MINIPAGE) -title "Changelog" -outline-in-sidebar -readme-to-index CHANGELOG.md > docs/CHANGELOG.html
	$(MINIPAGE) -title "Changelog" -outline-in-sidebar -readme-to-index CHANGELOG_ja.md > docs/CHANGELOG_ja.html

.PHONY: all test dist _dist clean manifest release docs
