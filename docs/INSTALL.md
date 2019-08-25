# operative-framework
[![Go Report Card](https://goreportcard.com/badge/github.com/graniet/operative-framework)](https://goreportcard.com/report/github.com/graniet/operative-framework) [![GoDoc](https://godoc.org/github.com/graniet/operative-framework?status.svg)](http://godoc.org/github.com/graniet/operative-framework) [![GitHub release](https://img.shields.io/github/release/graniet/operative-framework.svg)](https://github.com/graniet/operative-framework/releases/latest) [![LICENSE](https://img.shields.io/github/license/graniet/operative-framework.svg)](https://github.com/graniet/operative-framework/blob/master/LICENSE)

## Installing

### Manually

#### Download sources
```
go get -d github.com/graniet/operative-framework
cd $GOPATH/src/github.com/graniet/operative-framework
```

#### Get dependencies
```
go get github.com/Masterminds/glide
glide install --strip-vendor
```
#### Build binary
```
go build
./operative-framework
```

### Starting the `operative-framework` Shell

Once installed, run the optional `operative-framework` autocompleter with interactive help:

    $ operative-framework

Running the optional `operative-framework` shell will provide you with autocompletion, interactive help, fish-style suggestions, etc

