package main

import (
	"fmt"
	"hash/fnv"
	"time"
)

func hash(input []byte) string {
	hash := fnv.New64a()
	hash.Write(input)
	return fmt.Sprintf("%x", hash.Sum64())
}

func hashString(input string) string {
	return hash([]byte(input))
}

func hashTime(input time.Time) string {
	return hashString(input.Format(time.RFC3339Nano))
}

func hashTag(input string, at time.Time) string {
	return hashString(input + ":" + at.Format(time.RFC3339Nano))
}
