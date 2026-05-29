# Kratos Project Template

A project template for creating new Kratos services with HTTP and gRPC
transports, protobuf-first APIs, Wire dependency injection, OpenAPI generation,
and a small CRUD example.

Use this repository as a starting point for a new service. The included sample
resource is only reference code for API shape, layering, code generation, and
testing. Replace it with your own domain model when creating a real project.

## Create a New Project

1. Copy or generate a repository from this template.
2. Update the Go module path:

```bash
go mod edit -module github.com/your-org/your-service
```

3. Replace existing import paths that reference this template module.
4. Rename the command, service metadata, and sample API package to match your
   service.
5. Replace the sample CRUD resource with your own resource.
6. Regenerate code and verify the project:

```bash
make all
go test ./...
```

## What Is Included

- Kratos HTTP and gRPC server setup.
- Protobuf API definitions and generated Go code.
- OpenAPI generation.
- Wire-based dependency injection.
- Layered `service`, `biz`, and `data` packages.
- A lightweight in-memory repository for the sample resource.
- Unit tests for the service layer.
- Server-streaming and bidirectional-streaming examples.

## Project Layout

```text
api/                  Protobuf APIs and generated bindings
cmd/                  Application entrypoints
configs/              Local configuration
internal/server/      HTTP and gRPC server construction
internal/service/     Transport-facing service methods
internal/biz/         Usecases, entities, errors, repository interfaces
internal/data/        Repository implementations
third_party/          Protobuf dependencies
openapi.yaml          Generated OpenAPI document
```

## API Template Practices

The sample CRUD API demonstrates common conventions for Kratos projects:

- Resource-oriented methods: create, get, list, update, delete.
- HTTP annotations with `google.api.http`.
- Required fields with `google.api.field_behavior`.
- List requests with `page_size`, `page_token`, `filter`, and `order_by`.
- Pagination with `go.einride.tech/aip/pagination`.
- Partial updates with `google.protobuf.FieldMask` and `fieldmask.Update`.
- Streaming RPC definitions for one-way and bidirectional streams.

The in-memory data layer intentionally stays simple. It demonstrates flow across
layers, but does not implement a full query engine. Real repositories can apply
parsed filters and ordering in SQL, Ent, or another storage layer.

## Development Commands

Install generators:

```bash
make init
```

Regenerate API bindings and OpenAPI:

```bash
make api
```

Regenerate config protobufs:

```bash
make config
```

Run all generation steps, Wire, and module cleanup:

```bash
make all
```

Build:

```bash
make build
```

Test:

```bash
go test ./...
```

## Run Locally

```bash
go run ./cmd/server -conf ./configs
```

Default local ports are configured in `configs/config.yaml`:

- HTTP: `0.0.0.0:8000`
- gRPC: `0.0.0.0:9000`

## Docker

```bash
docker build -t <your-image-name> .
docker run --rm -p 8000:8000 -p 9000:9000 \
  -v </path/to/your/configs>:/data/conf \
  <your-image-name>
```
