#!/bin/bash
# Runs tests and generates code coverage HTML file.
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
