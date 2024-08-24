package config

import (
	"flag"
	"log"
	"time"
)

type Config struct {
	Port      string
	JwtSecret string
	JwtTTL    time.Duration
}

func ParseFlags() *Config {
	port := flag.String("port", "8080", "Port to listen on")
	jwtSecret := flag.String("jwt-secret", "", "JWT Secret Key")
	jwtTTL := flag.Duration("jwt-ttl", time.Hour, "JWT TTL (default 1h)")

	flag.Parse()

	if *jwtSecret == "" {
		log.Fatal("JWT Secret Key must be provided")
	}

	return &Config{
		Port:      *port,
		JwtSecret: *jwtSecret,
		JwtTTL:    *jwtTTL,
	}
}
