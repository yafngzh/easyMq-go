#!/bin/bash
cd ./cmd/server && go run main.go
cd ../producer && go run main.go
cd ../consumer && go run main.go
