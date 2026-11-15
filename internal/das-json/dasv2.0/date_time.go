package dasjson

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

import (
	"time"
)

// only include down to seconds
const DateTimeFormat = "2006-01-02T15:04:05Z"

type DateTime struct {
	time.Time `xml:",chardata"`
}

func Now() DateTime {
	return DateTime{Time: time.Now()}
}

func NewDateTime(t time.Time) DateTime {
	return DateTime{
		Time: t,
	}
}

func ParseDateTime(s string) (DateTime, error) {
	x, err := time.Parse(DateTimeFormat, s)
	return DateTime{x}, err
}

func MustParseDateTime(s string) DateTime {
	x, err := time.Parse(DateTimeFormat, s)
	if err != nil {
		panic(err)
	}
	return DateTime{x}
}

func MustParseDateTimePtr(s string) *DateTime {
	x := MustParseDateTime(s)
	return &x
}
