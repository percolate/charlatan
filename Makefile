# -*- mode: Makefile-gmake -*-

SHELL := bash

TOP_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

DIAGRAM_DIR := doc/diagram
DIAGRAMS := $(DIAGRAM_DIR)/architecture.png

BUILD_DIR := build
COVERAGE_DIR := $(BUILD_DIR)/coverage
TESTDATA_SOURCES := $(shell find testdata -name "*_def.go")
GENERATED_TESTDATA := $(subst _def,,$(TESTDATA_SOURCES))

all: test

doc/diagram/%.png: doc/%.dot
	@mkdir -p $(DIAGRAM_DIR)
	dot -Tpng $< > $@

doc: $(DIAGRAMS)

$(BUILD_DIR):
	@mkdir -p $@

$(COVERAGE_DIR):
	@mkdir -p $@

clean:
	rm -rf $(BUILD_DIR)

fmt:
	go fmt ./...

vet:
	go vet ./...

build_executable:
	go build

$(GENERATED_TESTDATA): $(TESTDATA_SOURCES) build_executable
	var1=$(subst .go,,$(@F)); var2=`echo $${var1:0:1} | tr  '[a-z]' '[A-Z]'`$${var1:1}; ./charlatan -file=$(subst .go,_def.go,$@) -output=$@ $$var2

test: $(COVERAGE_DIR)
	go test -v -coverprofile=$(TOP_DIR)/$(COVERAGE_DIR)/$(@F)_coverage.out -covermode=atomic ./...

generate_testdata: $(GENERATED_TESTDATA)

.PHONY: clean doc vet fmt test generate_testdata
