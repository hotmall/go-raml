# go-raml
[![Build Status](https://travis-ci.org/Jumpscale/go-raml.svg?branch=master)](https://travis-ci.org/Jumpscale/go-raml)


Table of Contents
=================

* [What is go-raml](#what-is-go-raml)
* [Supported RAML Versions](#raml-versions)
* [Install & Build in Development](#install)
* [Code Generation](#code-generation)
* [Usage](#usage)
* [Generating Server](#generating-server)
  * [Go Server](#go-server)
  * [Python Server](#python-server)
* [Generating Client](#generating-client)
* [Generating Docs](#generating-docs)
* [Using Generated Code](#using-generated-code)
  * [Simple Homepage & API Docs](#simple-home-page-and-api-docs)
  * [Using Go Server](#using-go-server)
  * [Using Python Server](#using-python-server)
  * [Using Go Client Library](#using-go-client-library)
  * [Using Python Client Library](#using-python-client-library)
* [Specification File](#specification-file)
* [Viewing and Editing RAML file](#viewing-and-editing-raml-file)
* [Contribute](#contribute)
* [Roadmap](#roadmap)
* [RAML to Code Translation](#raml-to-code-translation)
* [Tutorial](#tutorial)

## planning

- [kanban for 9.3.0 (includes jumpscale as well)](https://waffle.io/Jumpscale/home?milestone=9.3.0)
- [kanban for 9.3.1 (includes jumpscale as well)](https://waffle.io/Jumpscale/home?milestone=9.3.1)
- [kanban for 9.4.0 (includes jumpscale as well)](https://waffle.io/Jumpscale/home?milestone=9.4.0)

## What is go-raml

When creating and maintaining api's, there are two approaches:
* Design first

    You design the api, all methods, descriptions and types from the api point of view. Afterwards, you generate all the boilerplate code and documentation to bootstrap the code.

* Code first

    When modifying the implementation, the interface definitions need to be kept in sync. Besides the good practice of keeping your specification up to date, it is used by other tools to generate clients, documentation, ...

This tool supports both (or at least, this is on the roadmap).
As a specification format, it uses [RAML 1.0](http://raml.org) .

It currently has these features:

- generate server stub in Go, Python, and Nim from an raml file.
- generate complete client library in Go, Python, and Nim from an RAML file.
- generate [capnp](https://capnproto.org) schema. See [capnp docs](./docs/capnp.md) for details.

## RAML versions
Only RAML version 1.0 RC is supported.

Currently there are still some [limitations](docs/limitations.md) on the RAML 1.0 features that are supported.

## Install

make sure you have at least go 1.8 installed !

`go get -u github.com/hotmall/go-raml`


### Build in development

Install `dep` as package manager, we need it for vendoring tool

`$go get -u github.com/golang/dep/cmd/dep

Install go-bindata, we need it to compile all resource files to .go file

`go get -u github.com/jteeuwen/go-bindata/...`


To rebuild APIDocs files to .go file
```
cd $GOPATH/src/github.com/hotmall/go-raml
sh build_apidocs.sh
```

Build go-raml and all resource files
```
cd $GOPATH/src/github.com/hotmall/go-raml
make install
```

## Code generation

Internally, go templates are used to generate the code, this provides a flexible way to alter the generated code and to add different languages for the client.

## Usage

To use it on the commandline, just execute `go-raml` without any arguments, it will output the help on the stdout.

To invoke the codegeneration using `go generate`, specify the generation in 1 of your go source files:
`//go:generate go-raml ...`

go-raml needs to be on the path for this to work off course.


## Generating Server

`go-raml` able to generate Go & Python server. Nim server is under development.

Generated server will listen on port 5000


### Go Server

To generate the Go code for implementing the server in first design approach, execute

`go-raml server -l go --dir /gopath/src/github.com/mycompany/myapi --ramlfile api.raml --import-path github.com/mycompany/myapi`

`--dir` specify generated code directory. Except you really know what you do, use $GOPATH as root directory. In above example $GOPATH is /gopath.


`--ramlfile` is the path RAML file

`--import-path` is the root import path of the code. Generated code contains sub packages and need it to properly import code from another packages



### Python Server

Executes this command to generates Flask server

`go-raml server -l python --dir ./result_directory --ramlfile api.raml`

Executes this command to generates Sanic server

`go-raml server -l python --dir ./result_directory --ramlfile api.raml --kind sanic`

### Code Generator Options

```
   ----language, -l "go"                        Language to construct a server for
   --kind                                       Kind of server to generate (, sanic) (only for python)
   --dir "."                                    target directory
   --ramlfile "."                               Source raml file
   --package "main"                             package name
   --no-main                                    Do not generate a main.go file
   --no-apidocs                                 Do not generate API Docs in /apidocs/ endpoint
   --import-path "examples.com/ramlcode"        import path of the generated code
   --lib-root-urls								Array of libraries root URLs 
```

## Generating Client

`go-raml client --language go  --dir ./result_directory --ramlfile api.raml`

A go 1.5.x compatible client is generated in result_directory directory.

`go-raml client --language python --dir ./result_directory --ramlfile api.raml` or

`go-raml client --language python --dir ./result_directory --ramlfile api.raml --kind aiohttp`

A python 3.5 compatible client is generated in result_directory directory.

### Code Generator Options

```
OPTIONS:
   --language, -l "go"		Language to construct a client for
   --dir "."			target directory
   --ramlfile "."		Source raml file
   --package "client"		package name
   --import-path 		golang import path of the generated code
   --kind "requests"		Kind of python client to generate (requests,aiohttp)
   --lib-root-urls 		Array of libraries root URLs
   --python-unmarshall-response	set to true for python client to unmarshall the response into python class
```

## Generating Docs
`go-raml docs [--format markdown] --ramlfile api.raml --output api.md`

A single file `markdown` documentation is generated. Note that only markdown format is supported at the moment.

## Using Generated Code

### Simple home page and API Docs

Both servers have simple HTML home page that can be accessed in http://localhost:5000.
it provide simple description of the API server and link to auto generated API Documentation
powered by [api-console](https://github.com/mulesoft/api-console)

### Using Go Server

Generated go server can be used as your API server skeleton. It gives you routing code, implementation skeleton, simple request validation,  and all struct based on raml types and request/response body. You can then extend the generated code to fit your need.

Generated code details:

- an optional main file which contain routing code and other inialization code
- Interfaces files, contain routing code which mapped one to one with RAML resources. It always regenerated
- Interface implementation. It is code skeleton of the resource implemetation, only generated when the file is not present. So you can easily modify it
- Struct for RAML types and request/response body, only generated when the file is not present.
- Validation code
- Helper files

The generated server uses [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux) as HTTP request multiplexer.

install required packages
```
 $go get github.com/gorilla/mux
 $go get gopkg.in/validator.v2
 $go get github.com/justinas/alice
```

Build the code

`go build ./...`

Execute it

`./binary_name`

The server will then run in port 5000, you can go to http://localhost:5000 to see default html page as described above

### Using Python Server

The generated code is utilizing [Flask Blueprint](http://flask.pocoo.org/docs/0.11/blueprints/) to give you modular flask code.
The Sanic version also use Sanic Blueprint.

Generated code details:

- main file named app.py which initialize the app
- one bluperint/module for each root RAML resource
- helper file
- requirements.txt which contains list of packages needed to run the app.
- all files are overwritten when the server is regenerated even if they exist except for the handlers files under handlers/

Install needed packages
```
pip3 install -r requirements.txt
```

You might want to install it inside virtualenv

Run the code
```
python3 app.py
```

The server will then run in port 5000, you can go to http://localhost:5000 to see default html page as described above

### Using Go client library

Generated go client library can be used directly as modular package of your project. 

It has `AuthHeader` field, which if not empty will be used as value of `Authorization` header on each request.

### Using Python Client library

We provide two kind of clients:

- sync client using popular `requests` http library
- async client using `aiohttp` library, this client will give more performance

It has `set_auth_header` method to set `Authorization` header value on each request.
All files are overwritten when the client is regenerated.

## Specification file

Besides generation of a new RAML specification file, updating an existing raml file is also supported. This way the raml filestructure that can be included in the main raml file is honored.

`go-raml spec ...`

## Viewing and Editing RAML File

There are many ways to view and edit RAML file:

- [API Workbench](http://apiworkbench.com/) is an atom package for designing, building, testing, documenting and sharing RESTful HTTP APIs
- Vim : it is enough for simple purpose


## Contribute

When you want to contribute to the development, follow the [contribution guidelines](contributing.md).

## Roadmap
**v0.1**

* Generation of the server using [gorilla mux](http://www.gorillatoolkit.org/pkg/mux)
* Generation of a go client
* Generation of a python 3.5 client
* Generation of a python flask server

**v0.2**

* OAuth 2.0 support
* Input validation according to the RAML type definitions

**v0.3**

* Possibility to generate a default server implementation to MongoDB

**v0.4**

* Generation of a new RAML specification file

**v0.5**

* Update of a RAML specification file


## RAML to Code Translation

Short descriptions about how the generator generates the code from raml file are available in these docs:

- [Go](./docs/go_generator.md)
- [Python](./docs/python_generator.md)
- [Nim](./docs/nim_generator.md)
 
## Tutorial
 
 Tutorial for Go, Python, and Nim is available at [docs/tutorial directory](./docs/tutorial)
 
 Check all the available go commands [here](./docs/tutorial/go/README.md)
 
 Check all the available python commands [here](./docs/tutorial/python/README.md)
 
 Check all the available nim commands [here](./docs/tutorial/nim/README.md)

## CREDITS

- [go-raml/raml](https://github.com/go-raml/raml) for the original parser used by go-raml.
- [razor-1](https://github.com/razor-1) for the python class support.
