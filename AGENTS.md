# Repository Guidelines

## Project Structure & Module Organization

This is a Kratos service template. Keep API source and generated files under
`api/<domain>/<version>/`. Runtime code follows Kratos layers: service packages
handle transports, biz packages own usecases and repo interfaces, data packages
own persistence, and server packages wire HTTP/gRPC setup. Tests live beside the
code they cover.

## Build, Test, and Development Commands

- `make init`: install protoc, Kratos, OpenAPI, Wire generators.
- `make api`: regenerate API `.pb.go`, HTTP/gRPC bindings, and `openapi.yaml`.
- `make config`: regenerate internal config protobuf code.
- `make all`: run generation, Wire, and `go mod tidy`.
- `make build`: build all packages into `bin/`.
- `go run ./cmd/<app> -conf ./configs`: run locally; replace `<app>`.
- `go test ./...`: run all tests.

## Coding Style & Naming Conventions

Use `gofmt`. Keep package names short and lowercase; use MixedCaps for exported
identifiers. Write complete-sentence comments for exported types and functions.
Do not edit generated files manually; use `make api`, `make config`, or
`go generate ./...`.

## Kratos CRUD Template Reference

Use the sample CRUD implementation as a template. Replace the sample resource
name and keep the same boundaries.

- Proto: define `Create<Resource>`, `Get<Resource>`, `List<Resources>`,
  `Update<Resource>`, and `Delete<Resource>` under the resource API package.
  Include `google.api.http`, `google.api.field_behavior`, `FieldMask`, and list
  fields `page_size`, `page_token`, `filter`, and `order_by`.
- Service: embed generated `Unimplemented...Server`, convert protobuf messages
  to biz models, parse AIP fields with `filtering`, `ordering`, and
  `pagination`, and apply updates with `fieldmask.Update`.
- Biz: define the domain model, usecase, repo interface, typed errors, and
  `ListOption` helpers: `ListFilter`, `ListOrderBy`, `ListOffset`, `ListLimit`.
- Data: keep template repos simple. In-memory repos should store copies and
  paginate deterministically. Real database repos can apply `Filter` and
  `OrderBy` later.
- Wiring: add constructors to provider sets, register HTTP/gRPC services, run
  `make api`, then regenerate the application injector.

## Testing Guidelines

Use Go's `testing` package. Name tests `Test<TypeOrBehavior>` in `*_test.go`.
Cover CRUD behavior, validation errors, pagination, field masks, and streaming
examples when changing the sample. Run `go test ./...` before submitting.

## Commit & Pull Request Guidelines

Recent commits use subjects such as `feat:`, `fix:`, `refactor:`, and
`chore(deps):`. Keep subjects imperative and scoped. PRs should describe the
change, note API/generated-code updates, link issues, and list validation
commands. For proto changes, include regenerated files and `openapi.yaml`.

## Security & Configuration Tips

Do not commit real credentials in `configs/config.yaml`. Keep `third_party/`
protobufs stable unless updating generator dependencies.
