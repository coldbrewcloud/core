deps:
	glide -q install

vet:
	go vet `glide nv`

test: deps vet
	go test -race -cover -p 1 `glide nv`

fmt:
	go fmt `glide nv`

build: test
	GOOS=linux GOARCH=amd64 go build -tags production -o bin/linux/amd64/launchablebot .

deploy: build
	coldbrew deploy

.PHONY: test deps vet fmt