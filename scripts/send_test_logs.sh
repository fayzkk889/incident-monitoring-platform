#!/bin/bash

# Script to send test logs to the Go API
# Usage: ./send_test_logs.sh [count]

COUNT=${1:-10}
API_URL="http://localhost:8080/api/logs"

echo "Sending $COUNT test log entries..."

for i in $(seq 1 $COUNT); do
  LEVEL="info"
  if [ $((i % 5)) -eq 0 ]; then
    LEVEL="error"
  elif [ $((i % 3)) -eq 0 ]; then
    LEVEL="warn"
  fi

  SERVICE="service-$((RANDOM % 3 + 1))"
  
  curl -X POST "$API_URL" \
    -H "Content-Type: application/json" \
    -d "{
      \"logs\": [{
        \"service\": \"$SERVICE\",
        \"level\": \"$LEVEL\",
        \"message\": \"Test log message #$i - $LEVEL level\",
        \"metadata\": {\"test\": true, \"iteration\": $i}
      }]
    }" > /dev/null 2>&1

  echo "Sent log $i/$COUNT"
  sleep 0.1
done

echo "Done! Check the dashboard at http://localhost:5173"
