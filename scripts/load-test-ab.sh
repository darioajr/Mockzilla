#!/bin/bash

# Configuration
HOST="http://localhost:8080"
ENDPOINT="/message"
DURATION=30          # Seconds
CONCURRENCY=1000
TOTAL_REQUESTS=1000000

echo "-- INITIAL VALIDATION --"
curl -X POST -H "Content-Type: application/json" -d '{"message":"test"}' $HOST$ENDPOINT || exit 1

echo "-- SETTING UP ENVIRONMENT --"
echo '{"message":"test"}' > message.json  # Create payload file

echo "-- ADJUSTING SYSTEM LIMITS --"
ulimit -n 100000 2>/dev/null  # Increase file descriptor limit

echo "-- EXECUTING LOAD TEST --"
ab -k -n $TOTAL_REQUESTS -c $CONCURRENCY -p message.json -T 'application/json' $HOST$ENDPOINT | tee load_test.txt

echo "-- KEY METRICS --"
grep -E "Time taken|Requests per second|Transfer rate|Failed requests" load_test.txt

echo "-- ERROR ANALYSIS --"
grep "Non-2xx responses" load_test.txt

echo "-- LATENCY DISTRIBUTION --"
grep "Total:" load_test.txt -A 10

# Cleanup
rm -f message.json

echo "-- TEST COMPLETED --"