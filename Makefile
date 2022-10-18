fe-install:
	cd ui && \
	npm install && \
	npm run prepare && \
	npx openapi-typescript ../server/spec.yaml --output src/lib/oapi/spec.ts

fe-dev:
	cd ui && \
	npm run dev

be-dev:
	go run main.go server --dev
