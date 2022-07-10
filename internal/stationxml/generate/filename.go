package main

import (
	"path/filepath"
	"regexp"
	"strings"
	//	"unicode"
)

var matchSnakeFirst = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchSnakeRest = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {

	snake := matchSnakeFirst.ReplaceAllString(str, "${1}_${2}")
	snake = matchSnakeRest.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToLower(snake)

	if strings.HasPrefix(snake, "_") {
		snake = "_" + strings.TrimLeft(snake, "_")
	}

	if strings.HasSuffix(snake, "_") {
		snake = strings.TrimRight(snake, "_") + "_"
	}

	for strings.Contains(snake, "__") {
		snake = strings.ReplaceAll(snake, "__", "_")
	}

	return snake
}

// FileName returns a snake-case version of the struct based on it's name, it will attempt to
// replace camel case runs with underscore breaks as well as full uppercase names. It will
// also attempt to avoid adding an extra suffix if it already exists.
func FileName(name, suffix string) string {

	if f := filepath.Base(name); f != "." {
		name = toSnakeCase(f)
	}

	for strings.HasSuffix(name, suffix) {
		name = strings.TrimSuffix(name, suffix)
	}

	return name + suffix
}
