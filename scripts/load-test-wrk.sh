#!/bin/bash

# Configuration
HOST="https://localhost:8443"
ENDPOINT="/message"
DURATION=30          # Seconds
THREADS=4            # CPU cores
CONNECTIONS=1000

echo "-- INITIAL VALIDATION --"
curl -k --http2 -X POST -H "Content-Type: application/json" -d '{"message":"test"}' $HOST$ENDPOINT || exit 1

echo "-- PREPARING TEST ENVIRONMENT --"
cat > post.lua <<EOF
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.body = '{"transactionId": "123e4567-e89b-12d3-a456-426614174000","timestamp": "2025-04-12T04:14:32Z","payload": "QUTDZ9emPNlWhVZVvXa9XNbthZBnfgImlNuUZtQ/g9vGXHMj/lB5ThJL9q+91acSEjS6hp+n/BnoQYP/X8F8Tv4U/vLh7YMMjX/s5ov2nT5Jw9PjgsYDwbgBpQ0cxZjQfEE9g74p1HL4xk/JsbwMYWvju8J7Rfvy70sIDebh0TN0cBsyiyrLCXa4Ozbmv1LJlulkXD4hPHWkeevz5o1Kqawts7tZI+a6IzLJrPqaCc4HOzXEDTv6WOts9Eq6M1+aEZNU8bIY8OC3lGNpvverCwOcZBLFvKLxbReZFCZWhVuQidsTHMq7i+Ob1me+EXLo0hqcFvgKh97uKP/vuCSEoOI3hqBWoxD6NLqMd0Lvzgcuj8vIqOo7lCB1pO8Ss0RoQoVoLkkKH687srl3z6qXu6DOHpBdcAnDWTTp2JzLrtrXr/wAcmzbqnKPiCY8CMYuMlquu+dLq27wvSmW1fosN1LNCfqYFavQedtMiyvlptCoXWmPeikCPKjB7+2UaPJuIzw3o3jWtRsTFJsofaBCS5Eii5Xck0Llu9K6nQ0ktEF8FkUK8r6YE1I9J2wnJ65gWai4tTaO7ETFPPuO8N/qx943RW9DLDL0JlH91Qm3G3SCr0Iq7tqTsdsAZ0w4f432Gw/jKBaUbM9OxdsFQUganX7PeOBkpNGsJtRyHS7UoEI980CgMdAR1nXJmsa67MRLwnDWDeio49blYLvd0+DSmfT86MFGSMGPhMNls0i21T8KQ/m2wffvp7OSiNmsQJJV9TAEUgLGA3UyWWlG+1ISfYtUR6TZZ2R734QYs/N3eRqQv+dnYVQuGgF5XnpI+FOfKb9YUBGDzHGEQfzWkTrBB+Msi3lqJKFhRq3EoRuGa4WR/o7Gx9YMfHDyfPGo/cj8Aff89mymPY+OqtMNS7a3G7tazgrvehCLniejwN222CFLe2gFeLm/SQnFbGE6PmOg0UIQTuZ+eqDokyR+zsQ/DgBZKOUFj3JNS4WyMODooaZhqrelZsztKvupuYW38Y1AmHbC0YvXAyov5HO4vWQS7Q3sTKbzhoPxaqN6SaXOPrz7WKvXrebCGl7xZaBCmPD5+s1iraXOgcjP6Bn6Q5n0VheybXC2JS04T+6HIVRr1VrJKPp/VoK4qa7MPcqBUGtHYYSdNhb1eGKooGIWliCTHm8nUChcrOJQOVuLE+iENccqgg5d9xxzdu/U7MdV3coLlyE5tjhqH5sZb43EzCN4oANkcUCooM7ZBfzfubs9sYdbRJ2mWeYOqi6nk9X2gnGgqiuBCS7dypp7CrFZd36uE3mTA2ynQbwKdDktayzs6pE4JApQ3cNThoopt94lVRf1f/hnsyS8hqSiL0PezQa6lBuPrWNMYunrik2xbl4a8yTa08LwgUETczF/YLM7X28Nd4rYjaj97oV/iQpHe9nCmQf4RHUEL+lM2NU33qDPUzf7K+/hP3oI8rDcMCaRsfwqh5pV6Rr4pz8GfGgsqaNOePwm5anBeBVLGs5YwzoIZSYxw14wIvmEAqnxQXexEs2Bkg5/BDXreFlD7NujOxs8mklfVmcye+GZSO0foU9oQ9Kusqle4uIle+Qyf7CuoCwprv4ONFL7YkERfvJFKqrt8sJlQuw2M7btiblhvqygYC5mqcLcEWi9k4lssEKg36NwYxH2uuWOu+U7+aWgrClzFEOUGJwFmai0WrXCxOHYokuXdtrQZl9wB4M6zwmsveiV3XddhhEzMttx4onctDci7FY7hYC2obRUNYhkYUKMFOlROxB2jcDgQmQ/c2Ygt","metadata": {"service": "service-gateway","version": "1.2.3"}}'
EOF

echo "-- SYSTEM OPTIMIZATION --"
ulimit -n 100000 2>/dev/null  # Increase file descriptors

echo "-- EXECUTING LOAD TEST --"
wrk -t$THREADS -c$CONNECTIONS -d${DURATION}s -s post.lua --latency $HOST$ENDPOINT | tee load_test.txt

echo "-- PERFORMANCE METRICS --"
grep -E "Requests/sec|Latency|Socket errors" load_test.txt

echo "-- DETAILED LATENCY DISTRIBUTION --"
grep -A 5 "Latency Distribution" load_test.txt

echo "-- ERROR ANALYSIS --"
grep -E "errors|Non-2xx" load_test.txt

echo "-- CLEANUP --"
rm -f post.lua

echo "-- LOAD TEST COMPLETE --"