package delta

import (
	"context"
	"testing"
)

func TestDelta(t *testing.T) {
	if _, err := New(); err != nil {
		t.Fatal(err)
	}
}

func TestDB(t *testing.T) {

	set, err := NewBase("")
	if err != nil {
		t.Fatal(err)
	}

	files, err := NewResp("")
	if err != nil {
		t.Fatal(err)
	}

	db, err := NewDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.Init(context.Background(), set, files...); err != nil {
		t.Fatal(err)
	}
}
