package metadb

import (
	"reflect"
	"testing"

	"github.com/GeoNet/delta/meta"
	"github.com/hashicorp/go-memdb"
)

func TestMetaDB_Schema(t *testing.T) {

	db, err := memdb.NewMemDB(NewSchema())
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		t string
		i string
		a []interface{}
		r interface{}
	}{

		{"asset", "id", []interface{}{"TEST", "TEST", "TEST"}, meta.Asset{
			Equipment: meta.Equipment{
				Make:   "TEST",
				Model:  "TEST",
				Serial: "TEST",
			},
		}},
		{"constituent", "id", []interface{}{"TEST", 1}, meta.Constituent{
			Gauge:  "TEST",
			Number: 1,
		}},
		{"gauge", "id", []interface{}{"TEST"}, meta.Gauge{
			Reference: meta.Reference{
				Code:    "TEST",
				Network: "TEST",
			},
		}},
		{"mark", "id", []interface{}{"TEST"}, meta.Mark{
			Reference: meta.Reference{
				Code:    "TEST",
				Network: "TEST",
			},
		}},
		{"monument", "id", []interface{}{"TEST"}, meta.Monument{
			Mark: "TEST",
		}},
		{"mount", "id", []interface{}{"TEST"}, meta.Mount{
			Reference: meta.Reference{
				Code:    "TEST",
				Network: "TEST",
			},
		}},
		{"network", "id", []interface{}{"TEST"}, meta.Network{
			Code: "TEST",
		}},
		{"site", "id", []interface{}{"TEST", "TEST"}, meta.Site{
			Station:  "TEST",
			Location: "TEST",
		}},
		{"station", "id", []interface{}{"TEST"}, meta.Mark{
			Reference: meta.Reference{
				Code:    "TEST",
				Network: "TEST",
			},
		}},
	}

	txn := db.Txn(true)
	for _, test := range tests {
		if err := txn.Insert(test.t, test.r); err != nil {
			t.Fatal(err)
		}
	}
	txn.Commit()

	txn = db.Txn(false)
	defer txn.Abort()

	for _, test := range tests {
		v, err := txn.First(test.t, test.i, test.a...)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(v, test.r) {
			t.Errorf("db invalid check: %s [%s] got %v, expected %v", test.t, test.i, v, test.r)
		}
	}
}

func TestMetaDB_Multiple(t *testing.T) {

	db, err := memdb.NewMemDB(NewSchema())
	if err != nil {
		t.Fatal(err)
	}

	var entries = []meta.Asset{
		meta.Asset{
			Equipment: meta.Equipment{
				Make:   "A",
				Model:  "B",
				Serial: "C",
			},
		},
		meta.Asset{
			Equipment: meta.Equipment{
				Make:   "A",
				Model:  "B",
				Serial: "D",
			},
		},
		meta.Asset{
			Equipment: meta.Equipment{
				Make:   "A",
				Model:  "B",
				Serial: "E",
			},
		},
		meta.Asset{
			Equipment: meta.Equipment{
				Make:   "A",
				Model:  "BB",
				Serial: "CC",
			},
		},
		meta.Asset{
			Equipment: meta.Equipment{
				Make:   "A",
				Model:  "BB",
				Serial: "DD",
			},
		},
		meta.Asset{
			Equipment: meta.Equipment{
				Make:   "A",
				Model:  "A",
				Serial: "BB",
			},
		},
	}

	txn := db.Txn(true)
	for _, a := range entries {
		if err := txn.Insert("asset", a); err != nil {
			t.Fatal(err)
		}
	}
	txn.Commit()

	var tests = []struct {
		a []interface{}
		n int
	}{
		{[]interface{}{"A", "B"}, 3},
		{[]interface{}{"A", "BB"}, 2},
	}

	txn = db.Txn(false)
	defer txn.Abort()

	for _, test := range tests {
		var n int

		r, err := txn.Get("asset", "model", test.a...)
		if err != nil {
			t.Fatal(err)
		}

		for v := r.Next(); v != nil; v = r.Next() {
			n++
		}

		if n != test.n {
			t.Errorf("db invalid multiple check: %s [%s] got %d results, expected %d", "asset", "model", n, test.n)
		}
	}
}

func TestMetaDB_Load(t *testing.T) {
	db, err := NewMetaDB("../..")
	if err != nil {
		t.Fatal(err)
	}
	txn := db.Txn(false)
	defer txn.Abort()

	var tests = []struct {
		t string
		i string
		a []interface{}
	}{
		{"asset", "id", []interface{}{"Milne", "MILNE", "16"}},
		{"gauge", "id", []interface{}{"WLGT"}},
		{"constituent", "id", []interface{}{"WLGT", 1}},
		{"mark", "id", []interface{}{"AVLN"}},
		{"monument", "id", []interface{}{"AVLN"}},
		{"mount", "id", []interface{}{"TCEA"}},
		{"network", "id", []interface{}{"NZ"}},
		{"site", "id", []interface{}{"WEL", "20"}},
		{"station", "id", []interface{}{"WEL"}},
	}

	for _, test := range tests {
		v, err := txn.First(test.t, test.i, test.a...)
		if err != nil {
			t.Fatal(err)
		}
		if v == nil {
			t.Errorf("db invalid check: %s [%s] unable to find entry: %v", test.t, test.i, test.a)
		}
	}
}
