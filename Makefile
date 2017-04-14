deps:
	glide -q install

test: deps
	go test `glide nv`

.PHONY: deps test