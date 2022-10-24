package server

//go:generate mkdir -p oapi/
//go:generate go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0
//go:generate oapi-codegen --config=config.yaml spec.yaml
