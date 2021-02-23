#!/bin/sh

cd ./graphql
go get -u github.com/99designs/gqlgen
go run github.com/99designs/gqlgen
go mod tidy