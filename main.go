package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/freundallein/loadbalancer/bucket"
	"github.com/freundallein/loadbalancer/httpserv"
)

const (
	timeFormat      = "02.01.2006 15:04:05.000"
	serversEnvKey   = "ADDRS"
	portKey         = "PORT"
	staleTimeoutKey = "STALE_TIMEOUT"
)

type logWriter struct {
}

// Write - custom logger formatting
func (writer logWriter) Write(bytes []byte) (int, error) {
	msg := fmt.Sprintf("%s | %s", time.Now().UTC().Format(timeFormat), string(bytes))
	return fmt.Print(msg)
}

func getEnv(key string, fallback string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return fallback, nil
}

func getIntEnv(key string, fallback int) (int, error) {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fallback, err
		}
		return int(i), nil
	}
	return fallback, nil
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	addresses, err := getEnv(serversEnvKey, "")
	if err != nil {
		log.Fatalf("[config] %s", err.Error())
	}
	port, err := getIntEnv(portKey, 8000)
	if err != nil {
		log.Fatalf("[config] %s", err.Error())
	}
	staleTimeout, err := getIntEnv(staleTimeoutKey, 60)
	if err != nil {
		log.Fatalf("[config] %s", err.Error())
	}

	if len(addresses) == 0 {
		log.Fatal("[config] No addresses provided")
	}

	log.Println("[config] starting loadbalancer...")
	buckt, err := bucket.New(bucket.RoundRobin)
	if err != nil {
		log.Fatalf("[config] %s", err.Error())
	}
	items := strings.Split(addresses, ",")
	for _, addr := range items {
		srv, err := bucket.NewServer(addr)
		if err != nil {
			log.Fatal(err)
		}
		buckt.AddServer(srv)
		log.Printf("[config] server %s added\n", addr)
	}
	buckt.RunServices(staleTimeout)
	log.Println("[config] servers bucket started")
	server := httpserv.New(port, buckt)

	log.Printf("[config] httpserv started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
