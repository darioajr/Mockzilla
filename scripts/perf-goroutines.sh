#!/bin/bash

echo "Collecting Goroutine Analysis..."
curl -s http://localhost:8080/debug/pprof/goroutine > goroutine.pprof

echo "-- Goroutine --"
go tool pprof -top goroutine.pprof | head -20
