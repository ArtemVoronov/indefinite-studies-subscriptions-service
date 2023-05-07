#!/bin/sh

down() {
    docker-compose down
} 

build() {
    docker-compose build api
}

start() {
    docker-compose up -d
}

tail() {
    docker-compose logs -f
}

case "$1" in
  start)
    down
    build
    start
    tail
    ;;
  stop)
    down
    ;;
  tail)
    tail
    ;;
  purge)
    down
    ;;
  *)
    echo "Usage: $0 {start|stop|purge|tail}"
esac