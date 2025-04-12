#!/bin/bash

echo "Collecting Heap profile..."
curl -s http://localhost:8080/debug/pprof/heap > heap_profile.pprof

echo "-- MEMORY --"
go tool pprof -top heap_profile.pprof | head -20