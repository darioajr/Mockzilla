#!/bin/bash

echo "Starting load test 4 threads and 100 connections..."
wrk -t4 -c100 -d30s -s <(echo 'wrk.method="POST"; wrk.headers["Content-Type"]="application/json"; wrk.body="{\"message\":\"test\"}"') --latency http://localhost:8080/message > load_test.txt

echo "Collecting CPU profile..."
curl -s http://localhost:8080/debug/pprof/profile?seconds=30 > cpu_profile.pprof

echo "Collecting Memory profile..."
curl -s http://localhost:8080/debug/pprof/heap > heap_profile.pprof

echo "-- LOAD TEST --"
cat load_test.txt
echo -e "\n---------------------------------------------------\n"

echo "-- CPU --"
go tool pprof -top cpu_profile.pprof | head -20
echo -e "\n---------------------------------------------------\n"

echo "-- MEMORY --"
go tool pprof -top heap_profile.pprof | head -20
echo -e "\n---------------------------------------------------\n"