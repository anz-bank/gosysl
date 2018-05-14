Sysl plugin for Golang REST api generation
==========================================
[![Build Status](https://travis-ci.com/anz-bank/gosysl.svg?token=WkxpvzDNrfMxp7HeKSyA&branch=master)](https://travis-ci.com/anz-bank/gosysl)
[![Codecov](https://codecov.io/gh/anz-bank/gosysl/branch/master/graph/badge.svg?token=lRZ30tCTGK)](https://codecov.io/gh/anz-bank/gosysl)
[![Go Report Card](https://goreportcard.com/badge/github.com/anz-bank/gosysl)](https://goreportcard.com/report/github.com/anz-bank/gosysl)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)

`sysl-go-rest` generates RESTful API in the Go programming language for given Sysl specification.

Usage
-----
```bash
go install github.com/anz-bank/gosysl/cmd/sysl-go-rest
sysl pb example.sysl -o example.pb
sysl-go-rest example.pb
```

Compiling the protobuf file
---------------------------
[Protoc](https://github.com/google/protobuf/releases) and [Golang-Protobuf-plugin](https://github.com/golang/protobuf)
Use protoc compiler with golang Plugin to generate golang source code `sysl.pb.go`:

	protoc --go_out=. pb/sysl.proto
