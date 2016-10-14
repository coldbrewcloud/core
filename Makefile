VERSION := $(shell cat VERSION)

deps:
	glide -q install -s -u

test: deps
	go test `glide nv`

.PHONY: deps test