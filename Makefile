VERSION:=0.1.0

TARGETNAME:=rescueSummary
EXENAME:=$(TARGETNAME).exe
PACKANAME:=$(TARGETNAME)_$(VERSION)

GOOS=windows
GOARCH=amd64

SHELL=/bin/bash

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(EXENAME) main.go

.PHONY: pack
pack: build
	- rm -r target
	mkdir -p target/$(PACKANAME)
	cp $(EXENAME) target/$(PACKANAME)/
	cp README.md target/$(PACKANAME)/

.PHONY: clean
clean:
	- rm $(EXENAME)
	- rm -rf target
	- rm error.log
	- rm summary.html
