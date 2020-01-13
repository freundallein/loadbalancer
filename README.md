# loadbalancer
[![Build Status](https://travis-ci.org/freundallein/loadbalancer.svg?branch=master)](https://travis-ci.org/freundallein/loadbalancer)  

Round-robin http load balancer

Proxy incoming request to provided servers bucket with round-robin.  
Every 5 sec check server's availability.  
Every STALE_TIMEOUT minutes delete unreachable servers from bucket.


## Configuration
Application supports configuration via environment variables:
```
PORT=8000 (default 8000)
STALE_TIMEOUT=60 (default 60 - minutes)
ADDRS=http://service-1:9000,http://service-2:9001 (default empty)
```
## Installation
### With docker  
```
$> docker pull freundallein/go-lb
```
### With source
```
$> git clone git@github.com:freundallein/loadbalancer.git
$> cd loadbalancer
$> make build
```

## Usage

```
version: "3.5"

networks:
  network:
    name: example-network
    driver: bridge

services:
  loadbalancer:
    image: freundallein/go-lb:latest
    container_name: loadbalancer
    restart: always
    environment: 
      - PORT=8000
      - ADDRS=http://service-one:8000,http://service-two:8000,http://service-three:8000
      - STALE_TIMEOUT=1
    networks: 
      - network
    ports:
      - 8000:8000

  service-one:
    image: crccheck/hello-world
    container_name: hello-world-one
    restart: always
    networks: 
      - network

  service-two:
    image: crccheck/hello-world
    container_name: hello-world-two
    restart: always
    networks: 
      - network

  service-three:
    image: crccheck/hello-world
    container_name: hello-world-three
    restart: always
    networks: 
      - network

```
## Metrics
Default prometheus metrics are available on `/metrics`  
Custom metric - `lb_bucket_size 1` - represents the total number of servers in bucket

Good luck.