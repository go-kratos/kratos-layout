# Kratos Layout

## Install Kratos
```
go get github.com/go-kratos/kratos/cmd/kratos/v2

```
## Create a service
```
# create a template project
kratos new helloworld

cd helloworld
# Add a proto template
kratos proto add api/helloworld/helloworld.proto
# Generate the source code of service by proto file
kratos proto server api/helloworld/helloworld.proto -t internal/service

go generate ./...
cd cmd/helloworld
go build
./helloworld
```
## Automated Initialization (wire)
```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```
