ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
else
    SET=export
endif
GOOPT:=-ldflags "-s -w"

all:
	go fmt
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)
