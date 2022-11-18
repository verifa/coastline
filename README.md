# Coastline

> The self-service portal to automate the mundane

[![Go Report Card](https://goreportcard.com/badge/github.com/verifa/coastline)](https://goreportcard.com/report/github.com/verifa/coastline)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## What is Coastline?

Coastline is a web service that enables users to make Requests for predefined Request Templates, that trigger workflows to automate mundane tasks (like giving users access, creating new resources, etc.).

### How it works?

![how-it-works](./docs/how-it-works.excalidraw.png)

The Platform Team defines Request Templates and associated Workflows using CUE, and Coastline provides a web portal for users to login and make Requests (based on the Request Templates).

There is an approval process for Requests, and once approved, the server triggers a workflow.

### Who is it for?

Coastline is aimed at Platform Teams who maintain infrastructure and tooling for software teams, to provide a self-service portal for automating mundane tasks.

### Goal

The goal of Coastline is to automate simple tasks to provide better developer experience and allow platform teams to focus on things that matter.

### Example request

Below is a simple example request and workflow for requesting Cat Facts from: <https://catfact.ninja/>

```cue
package demo

import (
 "encoding/json"
 
 "github.com/verifa/coastline/tasks/http"
)

request: #CatFact: {
 kind:        "CatFact"
 description: "Cat fact max length \(spec.maxLength)"
 serviceSelector: {
  matchLabels: {
   tool: "cat-facts"
  }
 }
 spec: {
  // Max length of cat fact
  maxLength: int | *100
 }
}

workflow: CatFact: {
 input: request.#CatFact

 step: api: http.Get & {
  url: "https://catfact.ninja/fact"
  request: {
   params: {
    max_length: "\(input.spec.maxLength)"
   }
  }
 }

 output: {
  fact: json.Unmarshal(step.api.response.body).fact
 }
}
```

## Run dev server

Coastline's server has a `--dev` mode which is intended for exploring Coastline without setting anything up.

To give you a quick feel for Coastline, it comes packaged with some demo data to request useful things like:

1. Cat facts
2. Pokemon facts

To run a demo environment of Coastline with these capabilities you will need the request templates and workflows available in the [demo](./examples/demo) folder.

```bash
git clone https://github.com/verifa/coastline.git

# Build from source and run coastline
# TODO: add docs to download and run from binary: https://github.com/verifa/coastline/releases
# TODO: add docs on running from Docker: https://hub.docker.com/r/verifa/coastline
make build run

# Go to http://localhost:3000
# Login (without credentials in dev mode)
# Explore and enjoy, and don't forget to tell us what you learnt :)
```

## Terminology

**Request Template** - A request template is written in CUE and defines the specification (i.e. inputs/parameters) for a Request which users can make

**Request** - A request is made by users according to a specific Request Template

**Trigger** - A trigger is automatically created when a Request is approved

**Workflow** - A workflow is written in CUE and defines what should happen when a trigger for a Request is made

**Task** - Workflows execute tasks that actually do something (like making HTTP requests)

## License

This code is released under the [Apache-2.0 License](./LICENSE).
