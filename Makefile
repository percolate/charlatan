# -*- mode: Makefile-gmake -*-

SHELL := bash

TOP_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

DIAGRAM_DIR := doc/diagram
DIAGRAMS := $(DIAGRAM_DIR)/architecture.png

BUILD_DIR := build
COVERAGE_DIR := $(BUILD_DIR)/coverage
TESTDATA_SOURCES := $(shell find testdata -name "*_def.go")
IGNORED_TESTDATA := testdata/_/__def.go testdata/emptier/emptier_def.go
GENERATED_TESTDATA := $(subst _def,,$(filter-out $(IGNORED_TESTDATA),$(TESTDATA_SOURCES)))

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

charlatan:
	go build -mod=vendor ./...

# Get the capitalized interface name from the filename and pass it to charlatan
%.go: %_def.go
	rm -f $@
	iface=$(*F); ./charlatan -dir=testdata/$(*F) -output=$@ $${iface^}

test: $(COVERAGE_DIR)
	go test -v -coverprofile=$(TOP_DIR)/$(COVERAGE_DIR)/$(@F)_coverage.out -covermode=atomic ./...

testdata: charlatan $(GENERATED_TESTDATA)

.PHONY: clean doc vet fmt charlatan test testdata
