Sysl plugin for Golang REST api generation
==========================================
`sysl-go-rest` generates RESTful API in the Go programming language for given Sysl specification.

Usage
-----
```bash
go install github.com/anz-bank/sysl-go-rest/cmd/sysl-go-rest
sysl pb example.sysl -o example.pb
sysl-go-rest example.pb
```

Compiling the protobuf file
---------------------------
[Protoc](https://github.com/google/protobuf/releases) and [Golang-Protobuf-plugin](https://github.com/golang/protobuf)
Use protoc compiler with golang Plugin to generate golang source code `sysl.pb.go`:

	protoc --go_out=. pb/sysl.proto
