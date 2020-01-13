package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"./bucket"
	"./httpserv"
)

const timeFormat = "02.01.2006 15:04:05.000"

type logWriter struct {
}

// Write - custom logger formatting
func (writer logWriter) Write(bytes []byte) (int, error) {
	msg := fmt.Sprintf("%s | %s", time.Now().UTC().Format(timeFormat), string(bytes))
	return fmt.Print(msg)
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	var addresses string
	var port int
	var staleTimeout int

	flag.StringVar(&addresses, "servers", "", "set server addresses, separated by `,`")
	flag.IntVar(&port, "port", 8000, "load balancer's port")
	flag.IntVar(&staleTimeout, "stale-timeout", 60, "set minutes when unreachable host becomes stale")
	flag.Parse()

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
