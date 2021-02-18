# Kratos Layout

## Install Kratos
```
go get github.com/go-kratos/kratos/cmd/kratos/v2
go get github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
go get github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2

# from source
cd cmd/kratos && go install
cd cmd/protoc-gen-go-http && go install
cd cmd/protoc-gen-go-errors && go install
```
## Create a service
```
# create a template project
kratos new helloworld

cd helloworld
# Add a proto template
kratos proto add api/helloworld/helloworld.proto
# Generate the source code of service by proto file
kratos proto service api/helloworld/helloworld.proto -t internal/service

make proto
make build
make test
```
## Automated Initialization (wire)
```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```
