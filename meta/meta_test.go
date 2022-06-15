package meta

import (
	"testing"
	"time"
)

func TestSpan_Overlaps(t *testing.T) {

	today := Span{
		Start: time.Date(2022, 6, 12, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2022, 6, 13, 0, 0, 0, 0, time.UTC),
	}

	if !today.Overlaps(today) {
		t.Error("exact spans should match themselves")
	}

	yesterday := Span{
		Start: time.Date(2022, 6, 11, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2022, 6, 12, 0, 0, 0, 0, time.UTC),
	}

	tomorrow := Span{
		Start: time.Date(2022, 6, 13, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2022, 6, 14, 0, 0, 0, 0, time.UTC),
	}

	if !today.Overlaps(tomorrow) {
		t.Error("exact spans should match the end times")
	}

	if yesterday.Overlaps(tomorrow) {
		t.Error("disjoint spans should not match")
	}

	upcoming := Span{
		Start: time.Date(2022, 6, 12, 0, 0, 0, 0, time.UTC),
		End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	future := Span{
		Start: time.Date(2022, 6, 13, 0, 0, 0, 0, time.UTC),
		End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	if !upcoming.Overlaps(future) {
		t.Error("open ended spans should match")
	}

}

func TestSpan_Extents(t *testing.T) {

	today := Span{
		Start: time.Date(2022, 6, 12, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2022, 6, 13, 0, 0, 0, 0, time.UTC),
	}

	yesterday := Span{
		Start: time.Date(2022, 6, 11, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2022, 6, 12, 0, 0, 0, 0, time.UTC),
	}

	tomorrow := Span{
		Start: time.Date(2022, 6, 13, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2022, 6, 14, 0, 0, 0, 0, time.UTC),
	}

	if _, ok := today.Extent(yesterday, tomorrow); ok {
		t.Error("disjoint spans should not have an extent")
	}

	past := Span{
		Start: yesterday.Start,
		End:   tomorrow.End,
	}

	next := Span{
		Start: today.Start,
		End:   tomorrow.End,
	}

	if x, ok := today.Extent(past, next); !ok || x != today {
		t.Error("crossing spans should have an extent")
	}

	upcoming := Span{
		Start: time.Date(2022, 6, 12, 0, 0, 0, 0, time.UTC),
		End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	future := Span{
		Start: time.Date(2022, 6, 13, 0, 0, 0, 0, time.UTC),
		End:   time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	if x, ok := upcoming.Extent(future); !ok || x != future {
		t.Error("open ended spans should have a matching extent")
	}
}
