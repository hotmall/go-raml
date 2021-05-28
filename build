#!/bin/bash

go generate

# build
go install -v $(go list ./... | grep -v '/vendor/')

