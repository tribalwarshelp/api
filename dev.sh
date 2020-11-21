#!/bin/sh
export MODE=development
export GIN_MODE=debug
go run -race main.go