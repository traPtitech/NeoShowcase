package cliutil

import (
	"os"
	"strconv"
)

func GetEnvOrDefault(env, def string) string {
	s, ok := os.LookupEnv(env)
	if ok {
		return s
	}
	return def
}

func GetIntEnvOrDefault(env string, def int) int {
	s, ok := os.LookupEnv(env)
	if ok {
		return parseInt(s)
	}
	return def
}

func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
