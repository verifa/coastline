.PHONY: fe-install
fe-install:
	cd ui && \
	npm install

.PHONY: fe-dev
fe-dev:
	cd ui && \
	npm run dev

.PHONY: fe-build
fe-build:
	cd ui && \
	npm install && \
	npm run build

.PHONY: be-dev
be-dev:
	go run main.go server --dev --cue-dir ./examples/basic

.PHONY: be-gen
be-gen:
	go generate ./...

.PHONY: be-build
be-build: be-gen
	go build -o build/coastline --tags ui,oapi

.PHONY: build
build: fe-build be-build

.PHONY: run
run:
	./build/coastline server --dev
