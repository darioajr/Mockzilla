#!/bin/bash

echo "Collecting CPU profile..."
curl -s http://localhost:8080/debug/pprof/profile?seconds=30 > cpu_profile.pprof

echo "-- CPU --"
go tool pprof -top cpu_profile.pprof | head -20
