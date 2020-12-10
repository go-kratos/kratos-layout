# service-layout

## Install Kratos
```
mkdir $GOPATH/src/github.com/go-kratos/
cd $GOPATH/srcgithub.com/go-kratos/
git clone https://github.com/go-kratos/kratos.git
cd kratos
git checkout v2
go install ./...
```

## Create a project
```
kratos new helloworld
cd helloworld
kratos proto add api/helloworld/helloworld.proto
kratos proto service api/helloworld/helloworld.proto -t internal/service

make proto
make build
make test
```
