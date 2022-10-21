.PHONY: fe-install
fe-install:
	cd ui && \
	npm install && \
	npx openapi-typescript ../server/spec.yaml --output src/lib/oapi/spec.ts

.PHONY: fe-dev
fe-dev:
	cd ui && \
	npm run dev

.PHONY: fe-build
fe-build:
	cd ui && \
	npm run build

.PHONY: be-dev
be-dev:
	go run main.go server --dev

.PHONY: be-gen
be-gen:
	go generate ./...
	# Generate OpenAPI specification
	npx redoc-cli build -o server/oapi/index.html server/spec.yaml

.PHONY: be-build
be-build: be-gen
	go build -o build/coastline --tags ui,oapi

.PHONY: build
build: fe-build be-build

.PHONY: run
run:
	./build/coastline server --dev
