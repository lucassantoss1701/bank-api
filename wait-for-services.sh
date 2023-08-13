#!/bin/bash
wait_for_service() {
  local host=$1
  local port=$2
  local timeout=$3

  until nc -z $host $port >/dev/null 2>&1; do
    echo "Waiting for availability $host:$port..."
    sleep 1

    if (( timeout-- == 0 )); then
      echo "Timeout while waiting for availability of $host:$port"
      exit 1
    fi
  done
}

wait_for_service db 3306 30

echo "All services are ready. Starting the application..."
