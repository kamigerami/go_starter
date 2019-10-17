package env

import (
	"log"
	"os"
	"strconv"
)

func GetStringOrDie(key string) string {
	v, found := os.LookupEnv(key)
	if !found || len(v) == 0 {
		log.Fatalf("%s must be set!", key)
	}
	return v
}

func GetBoolOrDie(key string) bool {
	v, found := os.LookupEnv(key)
	if !found {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatalf("%s must be bool!", key)
	}
	return b
}

func GetUintOrDie(key string) uint64 {
	v, found := os.LookupEnv(key)
	if !found {
		log.Fatalf("%s must be set!", key)
	}
	i, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		log.Fatalf("%s must be int8!", key)
	}
	return i
}
